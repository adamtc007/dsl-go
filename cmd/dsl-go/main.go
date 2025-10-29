package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/example/dsl-go/internal/ebnf"
	"github.com/example/dsl-go/internal/manager"
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
		if fs.NArg() != 2 { fmt.Println("usage: dsl-go create <request_id> <template.sexpr>"); os.Exit(2) }
		id := fs.Arg(0); path := fs.Arg(1)
		b, err := os.ReadFile(path); must(err)
		v, hash, err := mgr.CreateRequest(id, string(b)); must(err)
		fmt.Printf("created %s v%d hash=%s\nstored at: %s/%s\n", id, v, hash, dataDir, id)

	case "show":
		fs := flag.NewFlagSet("show", flag.ExitOnError)
		fs.Parse(args)
		if fs.NArg() != 1 { fmt.Println("usage: dsl-go show <request_id>"); os.Exit(2) }
		id := fs.Arg(0)
		v, txt, err := mgr.GetCurrentText(id); must(err)
		fmt.Printf("-- %s v%d --\n%s\n", id, v, txt)

	case "validate":
		fs := flag.NewFlagSet("validate", flag.ExitOnError)
		fs.Parse(args)
		if fs.NArg() != 1 { fmt.Println("usage: dsl-go validate <file.sexpr>"); os.Exit(2) }
		path := fs.Arg(0)
		b, err := os.ReadFile(path); must(err)
		issues, err := mgr.ValidateText(string(b)); must(err)
		if len(issues)==0 { fmt.Println("OK") } else { for _, i := range issues { fmt.Println(i) } }

	case "compile":
		fs := flag.NewFlagSet("compile", flag.ExitOnError)
		fs.Parse(args)
		if fs.NArg() != 1 { fmt.Println("usage: dsl-go compile <file.sexpr>"); os.Exit(2) }
		path := fs.Arg(0)
		b, err := os.ReadFile(path); must(err)
		plan, err := mgr.CompilePlan(string(b)); must(err)
		j, _ := json.MarshalIndent(plan, "", "  ")
		fmt.Println(string(j))

	case "plan-delta":
		fs := flag.NewFlagSet("plan-delta", flag.ExitOnError)
		fs.Parse(args)
		if fs.NArg() != 2 { fmt.Println("usage: dsl-go plan-delta <from.sexpr> <to.sexpr>"); os.Exit(2) }
		a := fs.Arg(0); b := fs.Arg(1)
		ab, err := os.ReadFile(a); must(err)
		bb, err := os.ReadFile(b); must(err)
		delta, err := mgr.PlanDelta(string(ab), string(bb)); must(err)
		j, _ := json.MarshalIndent(delta, "", "  ")
		fmt.Println(string(j))

	case "ebnf":
		fmt.Println(ebnf.Text)

	case "parse-summary":
		fs := flag.NewFlagSet("parse-summary", flag.ExitOnError)
		fs.Parse(args)
		if fs.NArg() != 1 { fmt.Println("usage: dsl-go parse-summary <file.sexpr>"); os.Exit(2) }
		path := fs.Arg(0)
		b, err := os.ReadFile(path); must(err)
		p := parse.New()
		req, err := p.Parse(string(b)); must(err)
		if req.Meta != nil {
			fmt.Printf("request-id: %s\n", req.Meta.RequestID)
		}
		fmt.Printf("entities: %d\n", len(req.Orchestrator.Entities))
		fmt.Printf("resources: %d\n", len(req.Orchestrator.Resources))
		fmt.Printf("flows: %d\n", len(req.Orchestrator.Flows))
		if f, ok := req.Orchestrator.Flows["main"]; ok {
			fmt.Printf("  - main steps: %d\n", len(f.Steps))
		}

	case "ast-json":
		fs := flag.NewFlagSet("ast-json", flag.ExitOnError)
		fs.Parse(args)
		if fs.NArg() != 1 { fmt.Println("usage: dsl-go ast-json <file.sexpr>"); os.Exit(2) }
		path := fs.Arg(0)
		b, err := os.ReadFile(path); must(err)
		p := parse.New()
		req, err := p.Parse(string(b)); must(err)
		j, _ := json.MarshalIndent(req, "", "  ")
		fmt.Println(string(j))

	default:
		usage()
	}
}

func must(err error) {
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(os.Stderr, "not found")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
func usage() {
	fmt.Println(`dsl-go (Participle)
Commands:
  create <request_id> <template.sexpr>
  show <request_id>
  validate <file.sexpr>
  compile <file.sexpr>
  plan-delta <from.sexpr> <to.sexpr>
  ebnf
  parse-summary <file.sexpr>
  ast-json <file.sexpr>
`)
}
