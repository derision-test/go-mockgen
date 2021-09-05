package generation

type Options struct {
	ImportPaths       []string
	PkgName           string
	Interfaces        []string
	Exclude           []string
	OutputFilename    string
	OutputDir         string
	OutputImportPath  string
	Prefix            string
	Force             bool
	DisableFormatting bool
	GoImportsBinary   string
}
