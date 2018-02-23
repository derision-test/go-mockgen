package main

import "strings"

func decomposePackage(name string) (string, string) {
	index := strings.LastIndexByte(name, '.')
	if index == -1 {
		return "", name
	}

	return name[:index], name[index+1:]
}

func stringInSlice(needle string, haystack []string) bool {
	for _, elem := range haystack {
		if needle == elem {
			return true
		}
	}

	return false
}
