package command

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/derision-test/go-mockgen/internal/genlib/extraction"
	"github.com/derision-test/go-mockgen/internal/genlib/types"
)

func Extract(
	typeGetter types.TypeGetter,
	importPaths []string,
	targetNames []string,
) ([]*types.Interface, error) {
	extractor, err := extraction.NewExtractor()
	if err != nil {
		return nil, err
	}

	pkgs, err := extractor.Extract(importPaths)
	if err != nil {
		return nil, err
	}

	ifaces := []*types.Interface{}

	for _, name := range pkgs.GetNames() {
		if !shouldInclude(name, targetNames) {
			continue
		}

		iface, err := typeGetter(pkgs, name)
		if err != nil {
			return nil, err
		}

		if iface == nil {
			continue
		}

		for _, method := range iface.Methods {
			if !unicode.IsUpper([]rune(method.Name)[0]) {
				return nil, fmt.Errorf(
					"type '%s' has unexported an method '%s'",
					name,
					method.Name,
				)
			}
		}

		ifaces = append(ifaces, iface)
	}

	return ifaces, nil
}

func shouldInclude(name string, targetNames []string) bool {
	for _, v := range targetNames {
		if strings.ToLower(v) == strings.ToLower(name) {
			return true
		}
	}

	return len(targetNames) == 0
}
