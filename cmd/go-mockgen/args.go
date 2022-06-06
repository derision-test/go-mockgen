package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/derision-test/go-mockgen/internal/mockgen/consts"
	"github.com/derision-test/go-mockgen/internal/mockgen/generation"
	"github.com/derision-test/go-mockgen/internal/mockgen/paths"
	"gopkg.in/yaml.v3"
)

func parseAndValidateOptions() ([]*generation.Options, error) {
	allOptions, err := parseOptions()
	if err != nil {
		return nil, err
	}

	validators := []func(opts *generation.Options) (bool, error){
		validateOutputPaths,
		validateOptions,
	}

	for _, opts := range allOptions {
		for _, f := range validators {
			if fatal, err := f(opts); err != nil {
				if !fatal {
					kingpin.Fatalf("%s, try --help", err.Error())
				}

				return nil, err
			}
		}
	}

	return allOptions, nil
}

func parseOptions() ([]*generation.Options, error) {
	if len(os.Args) == 1 {
		return parseManifest()
	}

	opts, err := parseFlags()
	if err != nil {
		return nil, err
	}

	return []*generation.Options{opts}, nil
}

func parseFlags() (*generation.Options, error) {
	opts := &generation.Options{
		ImportPaths: []string{},
		Interfaces:  []string{},
	}

	app := kingpin.New(consts.Name, consts.Description).Version(consts.Version)
	app.UsageWriter(os.Stdout)

	app.Arg("path", "The import paths used to search for eligible interfaces").Required().StringsVar(&opts.ImportPaths)
	app.Flag("package", "The name of the generated package. It will be inferred from the output options by default.").Short('p').StringVar(&opts.PkgName)
	app.Flag("interfaces", "A list of target interfaces to generate defined in the given the import paths.").Short('i').StringsVar(&opts.Interfaces)
	app.Flag("exclude", "A list of interfaces to exclude from generation. Mocks for all other exported interfaces defined in the given import paths are generated.").Short('e').StringsVar(&opts.Exclude)
	app.Flag("dirname", "The target output directory. Each mock will be written to a unique file.").Short('d').StringVar(&opts.OutputDir)
	app.Flag("filename", "The target output file. All mocks are written to this file.").Short('o').StringVar(&opts.OutputFilename)
	app.Flag("import-path", "The import path of the generated package. It will be inferred from the target directory by default.").StringVar(&opts.PkgName)
	app.Flag("prefix", "A prefix used in the name of each mock struct. Should be TitleCase by convention.").StringVar(&opts.Prefix)
	app.Flag("constructor-prefix", "A prefix used in the name of each mock constructor function (after the initial `New`/`NewStrict` prefixes). Should be TitleCase by convention.").StringVar(&opts.ConstructorPrefix)
	app.Flag("force", "Do not abort if a write to disk would overwrite an existing file.").Short('f').BoolVar(&opts.Force)
	app.Flag("disable-formatting", "Do not run goimports over the rendered files.").BoolVar(&opts.DisableFormatting)
	app.Flag("goimports", "Path to the goimports binary.").Default("goimports").StringVar(&opts.GoImportsBinary)
	app.Flag("for-test", "Append _test suffix to generated package names and file names.").Default("false").BoolVar(&opts.ForTest)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return opts, nil
}

func parseManifest() ([]*generation.Options, error) {
	contents, err := os.ReadFile("mockgen.yaml")
	if err != nil {
		return nil, err
	}

	var payload struct {
		// Global options
		Exclude           []string `yaml:"exclude"`
		Prefix            string   `yaml:"prefix"`
		ConstructorPrefix string   `yaml:"constructor-prefix"`
		Force             bool     `yaml:"force"`
		DisableFormatting bool     `yaml:"disable-formatting"`
		Goimports         string   `yaml:"goimports"`
		ForTest           bool     `yaml:"for-test"`

		Mocks []struct {
			Path              string   `yaml:"path"`
			Paths             []string `yaml:"paths"`
			Package           string   `yaml:"package"`
			Interfaces        []string `yaml:"interfaces"`
			Exclude           []string `yaml:"exclude"`
			Dirname           string   `yaml:"dirname"`
			Filename          string   `yaml:"filename"`
			ImportPath        string   `yaml:"import-path"`
			Prefix            string   `yaml:"prefix"`
			ConstructorPrefix string   `yaml:"constructor-prefix"`
			Force             bool     `yaml:"force"`
			DisableFormatting bool     `yaml:"disable-formatting"`
			Goimports         string   `yaml:"goimports"`
			ForTest           bool     `yaml:"for-test"`
		} `yaml:"mocks"`
	}
	if err := yaml.Unmarshal(contents, &payload); err != nil {
		return nil, err
	}

	allOptions := make([]*generation.Options, 0, len(payload.Mocks))
	for _, opts := range payload.Mocks {
		// Mix
		opts.Exclude = append(opts.Exclude, payload.Exclude...)

		// Set if not overwritten in this entry
		if opts.Prefix == "" {
			opts.Prefix = payload.Prefix
		}
		if opts.ConstructorPrefix == "" {
			opts.ConstructorPrefix = payload.ConstructorPrefix
		}
		if opts.Goimports == "" {
			opts.Goimports = payload.Goimports
		}

		// Overwrite
		if payload.Force {
			opts.Force = true
		}
		if payload.DisableFormatting {
			opts.DisableFormatting = true
		}
		if payload.ForTest {
			opts.ForTest = true
		}

		// Validation
		paths := opts.Paths
		if opts.Path != "" {
			paths = append(paths, opts.Path)
		}
		if opts.Goimports == "" {
			opts.Goimports = "goimports"
		}

		allOptions = append(allOptions, &generation.Options{
			ImportPaths:       paths,
			PkgName:           opts.Package,
			Interfaces:        opts.Interfaces,
			Exclude:           opts.Exclude,
			OutputDir:         opts.Dirname,
			OutputFilename:    opts.Filename,
			OutputImportPath:  opts.ImportPath,
			Prefix:            opts.Prefix,
			ConstructorPrefix: opts.ConstructorPrefix,
			Force:             opts.Force,
			DisableFormatting: opts.DisableFormatting,
			GoImportsBinary:   opts.Goimports,
			ForTest:           opts.ForTest,
		})
	}

	return allOptions, nil
}

func validateOutputPaths(opts *generation.Options) (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return true, fmt.Errorf("failed to get current directory")
	}

	if opts.OutputFilename == "" && opts.OutputDir == "" {
		opts.OutputDir = wd
	}

	if opts.OutputFilename != "" && opts.OutputDir != "" {
		return false, fmt.Errorf("dirname and filename are mutually exclusive")
	}

	if opts.OutputFilename != "" {
		opts.OutputDir = path.Dir(opts.OutputFilename)
		opts.OutputFilename = path.Base(opts.OutputFilename)
	}

	if err := paths.EnsureDirExists(opts.OutputDir); err != nil {
		return true, fmt.Errorf(
			"failed to make output directory %s: %s",
			opts.OutputDir,
			err.Error(),
		)
	}

	if opts.OutputDir, err = cleanPath(opts.OutputDir); err != nil {
		return true, err
	}

	return false, nil
}

var goIdentifierPattern = regexp.MustCompile("^[A-Za-z]([A-Za-z0-9_]*[A-Za-z])?$")

func validateOptions(opts *generation.Options) (bool, error) {
	if len(opts.Interfaces) != 0 && len(opts.Exclude) != 0 {
		return false, fmt.Errorf("interface lists and exclude lists are mutually exclusive")
	}

	if opts.OutputImportPath == "" {
		path, ok := paths.InferImportPath(opts.OutputDir)
		if !ok {
			return false, fmt.Errorf("could not infer output import path")
		}

		opts.OutputImportPath = path
	}

	if opts.PkgName == "" {
		opts.PkgName = opts.OutputImportPath[strings.LastIndex(opts.OutputImportPath, string(os.PathSeparator))+1:]
	}

	if !goIdentifierPattern.Match([]byte(opts.PkgName)) {
		return false, fmt.Errorf("package name `%s` is illegal", opts.PkgName)
	}

	if opts.Prefix != "" && !goIdentifierPattern.Match([]byte(opts.Prefix)) {
		return false, fmt.Errorf("prefix `%s` is illegal", opts.Prefix)
	}

	if opts.ConstructorPrefix != "" && !goIdentifierPattern.Match([]byte(opts.ConstructorPrefix)) {
		return false, fmt.Errorf("constructor-`prefix `%s` is illegal", opts.ConstructorPrefix)
	}

	return false, nil
}

func cleanPath(path string) (cleaned string, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return "", err
	}

	if path, err = filepath.EvalSymlinks(path); err != nil {
		return "", err
	}

	return path, nil
}
