package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/example/dsl-go/internal/ebnf"
	"github.com/example/dsl-go/internal/generator"
	"github.com/example/dsl-go/internal/manager"
	"github.com/example/dsl-go/internal/mocks"
	"github.com/example/dsl-go/internal/parse"
)

func Run() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	dataDir := "./data"
	regDir := "./registry"

	mgr, err := manager.New(manager.Config{
		DataDir:     dataDir,
		RegistryDir: regDir,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating manager: %v\n", err)
		os.Exit(1)
	}

	cmds := map[string]func(){
		"create": func() {
			fs := flag.NewFlagSet("create", flag.ExitOnError)
			fs.Usage = func() {
				fmt.Println("usage: dsl-go create <request_id> <template_file>")
				fs.PrintDefaults()
			}
			if err := fs.Parse(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
				os.Exit(1)
			}
			if fs.NArg() != 2 {
				fs.Usage()
				return
			}
			reqID, templateFile := fs.Arg(0), fs.Arg(1)
			template, err := os.ReadFile(templateFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading template: %v\n", err)
				os.Exit(1)
			}
			version, hash, err := mgr.CreateRequest(reqID, string(template))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating request: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("created request %s, version %d, hash %s\n", reqID, version, hash)
		},
		"get": func() {
			fs := flag.NewFlagSet("get", flag.ExitOnError)
			fs.Usage = func() {
				fmt.Println("usage: dsl-go get <request_id>")
				fs.PrintDefaults()
			}
			if err := fs.Parse(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
				os.Exit(1)
			}
			if fs.NArg() != 1 {
				fs.Usage()
				return
			}
			reqID := fs.Arg(0)
			_, text, err := mgr.GetCurrentText(reqID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting request: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(text)
		},
		"validate": func() {
			fs := flag.NewFlagSet("validate", flag.ExitOnError)
			fs.Usage = func() {
				fmt.Println("usage: dsl-go validate <file>")
				fs.PrintDefaults()
			}
			if err := fs.Parse(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
				os.Exit(1)
			}
			if fs.NArg() != 1 {
				fs.Usage()
				return
			}
			file := fs.Arg(0)
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
				os.Exit(1)
			}
			issues, err := mgr.ValidateText(string(content))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error validating: %v\n", err)
				os.Exit(1)
			}
			if len(issues) > 0 {
				fmt.Println("Validation issues:")
				for _, issue := range issues {
					fmt.Printf("- %s\n", issue)
				}
				os.Exit(1)
			}
			fmt.Println("Validation successful")
		},
		"plan": func() {
			fs := flag.NewFlagSet("plan", flag.ExitOnError)
			fs.Usage = func() {
				fmt.Println("usage: dsl-go plan <file>")
				fs.PrintDefaults()
			}
			if err := fs.Parse(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
				os.Exit(1)
			}
			if fs.NArg() != 1 {
				fs.Usage()
				return
			}
			file := fs.Arg(0)
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
				os.Exit(1)
			}
			plan, err := mgr.CompilePlan(string(content))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error compiling plan: %v\n", err)
				os.Exit(1)
			}
			jsonPlan, _ := json.MarshalIndent(plan, "", "  ")
			fmt.Println(string(jsonPlan))
		},
		"gen": func() {
			fs := flag.NewFlagSet("gen", flag.ExitOnError)
			templateFile := fs.String("template", "", "Template file to use")
			fs.Usage = func() {
				fmt.Println("usage: dsl-go gen -template=<template_file> <scenario_file>")
				fs.PrintDefaults()
			}
			if err := fs.Parse(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
				os.Exit(1)
			}
			if fs.NArg() != 1 || *templateFile == "" {
				fs.Usage()
				return
			}
			scenarioFile := fs.Arg(0)

			loader := mocks.NewDefaultLoader()
			req, err := loader.LoadScenario(scenarioFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error loading scenario: %v\n", err)
				os.Exit(1)
			}

			req.DataDictionary = mgr.GetDataDictionary()

			gen, err := generator.New()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating generator: %v\n", err)
				os.Exit(1)
			}
			resp, err := gen.GenerateFromTemplateFile(*templateFile, req)

			if err != nil {
				fmt.Fprintf(os.Stderr, "error generating dsl: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(resp.DSL)
		},
		"dictionary": func() {
			fs := flag.NewFlagSet("dictionary", flag.ExitOnError)
			fs.Usage = func() {
				fmt.Println("usage: dsl-go dictionary <attribute_id>")
				fs.PrintDefaults()
			}
			if err := fs.Parse(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
				os.Exit(1)
			}
			if fs.NArg() != 1 {
				fs.Usage()
				return
			}
			attrID := fs.Arg(0)
			attr, ok := mgr.GetAttribute(attrID)
			if !ok {
				fmt.Fprintf(os.Stderr, "error: attribute %q not found\n", attrID)
				os.Exit(1)
			}
			fmt.Printf("AttributeID: %s\n", attr.AttributeID)
			fmt.Printf("Description: %s\n", attr.Description)
			fmt.Printf("VectorID:    %s\n", attr.VectorID)
		},
		"ebnf": func() {
			fmt.Println(ebnf.Text)
		},
		"ast-json": func() {
			fs := flag.NewFlagSet("ast-json", flag.ExitOnError)
			fs.Usage = func() {
				fmt.Println("usage: dsl-go ast-json <file>")
				fs.PrintDefaults()
			}
			if err := fs.Parse(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
				os.Exit(1)
			}
			if fs.NArg() != 1 {
				fs.Usage()
				return
			}
			file := fs.Arg(0)
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
				os.Exit(1)
			}
			parser, err := parse.New()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating parser: %v\n", err)
				os.Exit(1)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating parser: %v\n", err)
				os.Exit(1)
			}
			ast, err := parser.Parse(string(content))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error parsing file: %v\n", err)
				os.Exit(1)
			}
			jsonAST, _ := json.MarshalIndent(ast, "", "  ")
			fmt.Println(string(jsonAST))
		},
	}

	cmd, ok := cmds[os.Args[1]]
	if !ok {
		usage()
		return
	}
	cmd()
}

func usage() {
	fmt.Println("usage: dsl-go <command> [<args>]")
	fmt.Println("Commands:")
	fmt.Println("  create      Create a new onboarding request from a template")
	fmt.Println("  get         Get the latest version of an onboarding request")
	fmt.Println("  validate    Validate a DSL file")
	fmt.Println("  plan        Compile a DSL file into a plan")
	fmt.Println("  gen         Generate a DSL file from a scenario")
	fmt.Println("  ebnf        Print the EBNF grammar")
	fmt.Println("  ast-json    Print the AST of a DSL file as JSON")
	fmt.Println("  dictionary  Get information about a data dictionary attribute")
}
