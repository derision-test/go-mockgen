package main

import (
	"log"

	"github.com/efritz/go-mockgen/internal/genlib/command"
	"github.com/efritz/go-mockgen/internal/genlib/types"
	"github.com/efritz/go-mockgen/internal/mockgen"
)

const (
	name        = "go-mockgen"
	description = "go-mockgen generates mock implementations from interface definitions."
	packageName = "github.com/efritz/go-mockgen"
	version     = "0.1.0"
)

var Main = main

func init() {
	log.SetFlags(0)
	log.SetPrefix("go-mockgen: ")
}

func main() {
	if err := command.Run(name, description, version, types.GetInterface, mockgen.Generate); err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}
