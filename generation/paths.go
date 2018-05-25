package generation

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

func writeFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func getFilename(dirname, interfaceName, prefix string) string {
	filename := fmt.Sprintf("%s_mock.go", interfaceName)
	if prefix != "" {
		filename = fmt.Sprintf("%s_%s", prefix, filename)
	}

	return path.Join(dirname, strings.Replace(strings.ToLower(filename), "-", "_", -1))
}
