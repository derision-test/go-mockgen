package generation

type Options struct {
	ImportPaths       []string
	PkgName           string
	Interfaces        []string
	OutputFilename    string
	OutputDir         string
	OutputImportPath  string
	Prefix            string
	Force             bool
	DisableFormatting bool
	GoImportsBinary   string
}
