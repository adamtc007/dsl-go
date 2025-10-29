package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/example/dsl-go/internal/ebnf"
	"github.com/example/dsl-go/internal/generator"
	"github.com/example/dsl-go/internal/manager"
	"github.com/example/dsl-go/internal/mocks"
	"github.com/example/dsl-go/internal/parse"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	dataDir := "./data"
	regDir := "./registry"

	cmd := os.Args[1]
	args := os.Args[2:]

	mgr := manager.New(manager.Config{
		DataDir:     dataDir,
		RegistryDir: regDir,
	})

	switch cmd {
	case "create":
		fs := flag.NewFlagSet("create", flag.ExitOnError)
		fs.Parse(args)
		if fs.NArg() != 2 {
			fmt.Println("usage: dsl-go create <request_id>
