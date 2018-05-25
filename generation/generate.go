package generation

import (
	"fmt"
	"path"
	"strings"

	"github.com/efritz/go-mockgen/paths"
	"github.com/efritz/go-mockgen/specs"
)

func Generate(specs specs.Specs, pkgName, prefix, dirname, filename string, force bool) error {
	importPath, err := inferImportPath(dirname)
	if err != nil {
		return err
	}

	for _, spec := range specs {
		if spec.ImportPath == importPath {
			spec.ImportPath = ""
		}
	}

	if dirname != "" && filename == "" {
		return generateDirectory(specs, pkgName, prefix, dirname, force)
	}

	return generateFile(specs, pkgName, prefix, path.Join(dirname, filename), force)
}

func inferImportPath(path string) (string, error) {
	gopath := paths.Gopath()
	if strings.HasPrefix(path, gopath) {
		// gopath + /src/
		return path[len(gopath)+5:], nil
	}

	return "", fmt.Errorf("destination is outside $GOPATH")
}

func generateFile(specs specs.Specs, pkgName, prefix, filename string, force bool) error {
	content, err := generateContent(specs, pkgName, prefix)
	if err != nil {
		return err
	}

	if filename != "" {
		exists, err := paths.Exists(filename)
		if err != nil {
			return err
		}

		if exists && !force {
			return fmt.Errorf("filename %s already exists", filename)
		}

		return writeFile(filename, content)
	}

	fmt.Printf("%s\n", content)
	return nil
}

func generateDirectory(allSpecs specs.Specs, pkgName, prefix, dirname string, force bool) error {
	if !force {
		allPaths := []string{}
		for name := range allSpecs {
			allPaths = append(allPaths, getFilename(dirname, name, prefix))
		}

		conflict, err := paths.AnyExists(allPaths)
		if err != nil {
			return err
		}

		if conflict != "" {
			return fmt.Errorf("filename %s already exists", conflict)
		}
	}

	for name, spec := range allSpecs {
		content, err := generateContent(specs.Specs{name: spec}, pkgName, prefix)
		if err != nil {
			return err
		}

		if err := writeFile(getFilename(dirname, name, prefix), content); err != nil {
			return err
		}
	}

	return nil
}
