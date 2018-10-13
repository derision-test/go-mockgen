package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/efritz/go-genlib/command"
	"github.com/efritz/go-genlib/generator"
	"github.com/efritz/go-genlib/types"
)

const (
	Name        = "go-mockgen"
	Description = "go-mockgen generates mock implementations from interface definitions."
	Version     = "0.1.0"
)

func main() {
	if err := command.Run(Name, Description, Version, typeGetter, generate); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
}

func typeGetter(pkgs *types.Packages, name string) (*types.Interface, error) {
	return pkgs.GetInterface(name)
}

func generate(ifaces []*types.Interface, opts *command.Options) error {
	return generator.Generate(
		"github.com/efritz/go-mockgen",
		ifaces,
		opts,
		filenameGenerator,
		ifaceGenerator,
	)
}

func filenameGenerator(ifaceName string) string {
	return fmt.Sprintf("%s_mock.go", ifaceName)
}

func ifaceGenerator(file *jen.File, iface *types.Interface, prefix string) {
	newInterfaceGenerator(file, title(iface.Name), prefix, iface).generate()
}

func title(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}
