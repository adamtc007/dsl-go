package parse

import (
	"fmt"
	"strconv"
	"time"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/example/dsl-go/internal/ast"
)

/*
We parse S-expr into a generic tree (Sexpr) using Participle, then map to AST.

Grammar (tokens only; structure handled programmatically):

  "("  -> LParen
  ")"  -> RParen
  "->" -> Arrow (reserved for future)
  String: " ... "
  ColonIdent: :meta :orchestrator :lifecycle etc.
  Ident:      onboarding-request states initial draft ...

*/

var sexprLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Whitespace", Pattern: `\s+`},
	{Name: "Comment", Pattern: `\;[^\n]*`},
	{Name: "LParen", Pattern: `\(`},
	{Name: "RParen", Pattern: `\)`},
	{Name: "Arrow", Pattern: `->`},
	{Name: "String", Pattern: `"(?:\\.|[^"])*"`},
	{Name: "ColonIdent", Pattern: `:[A-Za-z][A-Za-z0-9_-]*`},
	{Name: "Ident", Pattern: `[A-Za-z][A-Za-z0-9_-]*`},
	{Name: "Number", Pattern: `[0-9]+(?:\.[0-9]+)?`}, // Add number support
})

type Sexpr struct {
	Pos lexer.Position
	// Either a list or an atom:
	List *List `  @@`
	Atom *Atom `| @@`
}

type List struct {
	Pos      lexer.Position
	Elements []*Sexpr `"(" @@* ")"`
}

type Atom struct {
	Pos    lexer.Position
	String *string `  @String`
	Number *string `| @Number` // Capture as string, parse later
	Sym    *string `| @Ident | @ColonIdent`
}

func buildParser() *participle.Parser[Sexpr] {
	return participle.MustBuild[Sexpr](
		participle.Lexer(sexprLexer),
		participle.Unquote("String"),
	)
}

// Public parser API
type Parser interface {
	Parse(text string) (*ast.Request, error)
}

type PartParser struct {
	p *participle.Parser[Sexpr]
}

func New() *PartParser {
	return &PartParser{p: buildParser()}
}

func (pp *PartParser) Parse(text string) (*ast.Request, error) {
	root, err := pp.p.ParseString("", text)
	if err != nil {
		return nil, err
	}
	return sexprToRequest(root)
}

/* ---------------- mapping Sexpr -> AST ---------------- */

func sexprToRequest(root *Sexpr) (*ast.Request, error) {
	if root == nil || root.List == nil || len(root.List.Elements) == 0 {
		return nil, fmt.Errorf("top level must be a list")
	}
	first := root.List.Elements[0]
	if !(first.Atom != nil && first.Atom.Sym != nil && *first.Atom.Sym == "onboarding-request") {
		return nil, fmt.Errorf("expected (onboarding-request ...)")
	}

	req := &ast.Request{
		Orchestrator: ast.Orchestrator{
			Lifecycle: ast.Lifecycle{},
			Entities:  map[string]ast.Entity{},
			Resources: map[string]ast.Resource{},
			Flows:     map[string]ast.Flow{},
		},
	}

	for _, sec := range root.List.Elements[1:] {
		if sec.List == nil || len(sec.List.Elements) < 2 {
			continue
		}
		head := sec.List.Elements[0]
		body := sec.List.Elements[1]
		if head.Atom == nil || head.Atom.Sym == nil {
			continue
		}
		switch *head.Atom.Sym {
		case ":meta":
			m, err := parseMeta(body)
			if err != nil {
				return nil, err
			}
			req.Meta = m
		case ":orchestrator":
			if body.List == nil {
				continue
			}
			for _, pair := range body.List.Elements {
				if pair.List == nil || len(pair.List.Elements) < 2 {
					continue
				}
				h := pair.List.Elements[0]
				val := pair.List.Elements[1]
				if h.Atom == nil || h.Atom.Sym == nil {
					continue
				}
				switch *h.Atom.Sym {
				case ":lifecycle":
					lc, _ := parseLifecycle(val)
					if lc != nil {
						req.Orchestrator.Lifecycle = *lc
					}
				case ":entities":
					ent, _ := parseEntities(val)
					req.Orchestrator.Entities = ent
				case ":resources":
					res, _ := parseResources(val)
					req.Orchestrator.Resources = res
				case ":flows":
					fl, _ := parseFlows(val)
					req.Orchestrator.Flows = fl
				default:
					// :policies parsed later
				}
			}
		default:
			// ignore :catalog etc. for now
		}
	}

	return req, nil
}

func parseMeta(body *Sexpr) (*ast.Meta, error) {
	m := &ast.Meta{}
	if body.List == nil {
		return m, nil
	}
	for _, kv := range body.List.Elements {
		if kv.List == nil || len(kv.List.Elements) < 2 {
			continue
		}
		k := atomText(kv.List.Elements[0])
		v := kv.List.Elements[1]
		switch k {
		case "request-id":
			m.RequestID = atomText(v)
		case "version":
			m.Version, _ = strconv.ParseUint(atomText(v), 10, 64)
		case "created-at":
			if t, err := time.Parse(time.RFC3339, atomText(v)); err == nil {
				m.CreatedAt = t
			}
		case "updated-at":
			if t, err := time.Parse(time.RFC3339, atomText(v)); err == nil {
				m.UpdatedAt = t
			}
		}
	}
	// defaults
	if m.CreatedAt.IsZero() {
		now := time.Now().UTC()
		m.CreatedAt, m.UpdatedAt = now, now
	} else if m.UpdatedAt.IsZero() {
		m.UpdatedAt = m.CreatedAt
	}
	return m, nil
}

func parseLifecycle(body *Sexpr) (*ast.Lifecycle, error) {
	lc := &ast.Lifecycle{}
	if body.List == nil {
		return lc, nil
	}
	for _, el := range body.List.Elements {
		if el.List == nil || len(el.List.Elements) == 0 {
			continue
		}
		key := atomText(el.List.Elements[0])
		switch key {
		case "states":
			for _, s := range el.List.Elements[1:] {
				lc.States = append(lc.States, atomText(s))
			}
		case "initial":
			if len(el.List.Elements) > 1 {
				lc.Initial = atomText(el.List.Elements[1])
			}
		case "transitions":
			// TODO: parse transitions (guards/effects)
		}
	}
	return lc, nil
}

func parseEntities(body *Sexpr) (map[string]ast.Entity, error) {
	m := make(map[string]ast.Entity)
	if body.List == nil {
		return m, nil
	}
	for _, el := range body.List.Elements {
		if el.List == nil || atomText(el.List.Elements[0]) != "entity" {
			continue
		}
		kmap := parseKeywordMap(el.List.Elements[1:])
		id := atomText(kmap[":id"])
		ent := ast.Entity{
			ID:  id,
			Typ: atomText(kmap[":type"]),
		}
		// Find (attrs ...)
		for _, subEl := range el.List.Elements[1:] {
			if subEl.List != nil && atomText(subEl.List.Elements[0]) == "attrs" {
				ent.Attrs, _ = parseAttrs(subEl.List.Elements[1:])
				break
			}
		}
		m[id] = ent
	}
	return m, nil
}

func parseAttrs(list []*Sexpr) (map[string]ast.AttrVal, error) {
	m := make(map[string]ast.AttrVal)
	for _, el := range list {
		if el.List == nil || len(el.List.Elements) < 2 {
			continue
		}
		key := atomText(el.List.Elements[0])
		val := atomValue(el.List.Elements[1])
		attr := ast.AttrVal{Value: val}
		// Check for metadata like :provenance
		kmap := parseKeywordMap(el.List.Elements[2:])
		if p, ok := kmap[":provenance"]; ok {
			s := atomText(p)
			attr.Provenance = &s
		}
		m[key] = attr
	}
	return m, nil
}

func parseResources(body *Sexpr) (map[string]ast.Resource, error) {
	m := make(map[string]ast.Resource)
	if body.List == nil {
		return m, nil
	}
	for _, el := range body.List.Elements {
		if el.List == nil || atomText(el.List.Elements[0]) != "resource" {
			continue
		}
		kmap := parseKeywordMap(el.List.Elements[1:])
		id := atomText(kmap[":id"])
		res := ast.Resource{
			ID:  id,
			Typ: atomText(kmap[":type"]),
		}
		// Find (requires ...) and (config ...)
		for _, subEl := range el.List.Elements[1:] {
			if subEl.List == nil || len(subEl.List.Elements) == 0 {
				continue
			}
			key := atomText(subEl.List.Elements[0])
			switch key {
			case "requires":
				res.Requires, _ = parseRequires(subEl.List.Elements[1:])
			case "config":
				res.Config, _ = parseConfig(subEl.List.Elements[1:])
			}
		}
		m[id] = res
	}
	return m, nil
}

func parseRequires(list []*Sexpr) ([]ast.RequireItem, error) {
	var items []ast.RequireItem
	for _, el := range list {
		if el.List == nil || len(el.List.Elements) != 2 {
			continue
		}
		items = append(items, ast.RequireItem{
			Kind: atomText(el.List.Elements[0]),
			ID:   atomText(el.List.Elements[1]),
		})
	}
	return items, nil
}

func parseConfig(list []*Sexpr) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for _, el := range list {
		if el.List == nil || len(el.List.Elements) != 2 {
			continue
		}
		m[atomText(el.List.Elements[0])] = atomValue(el.List.Elements[1])
	}
	return m, nil
}

func parseFlows(body *Sexpr) (map[string]ast.Flow, error) {
	m := make(map[string]ast.Flow)
	if body.List == nil {
		return m, nil
	}
	for _, el := range body.List.Elements {
		if el.List == nil || atomText(el.List.Elements[0]) != "flow" {
			continue
		}
		kmap := parseKeywordMap(el.List.Elements[1:])
		id := atomText(kmap[":id"])
		flow := ast.Flow{ID: id}
		// Find (steps ...)
		for _, subEl := range el.List.Elements[1:] {
			if subEl.List != nil && atomText(subEl.List.Elements[0]) == "steps" {
				flow.Steps, _ = parseSteps(subEl.List.Elements[1:])
				break
			}
		}
		m[id] = flow
	}
	return m, nil
}

func parseSteps(list []*Sexpr) ([]ast.Step, error) {
	var steps []ast.Step
	for _, el := range list {
		if el.List == nil || len(el.List.Elements) == 0 {
			continue
		}
		kind := atomText(el.List.Elements[0])
		switch kind {
		case "task":
			t, _ := parseTask(el)
			steps = append(steps, ast.Step{Kind: "task", Task: t})
		case "gate":
			g, _ := parseGate(el)
			steps = append(steps, ast.Step{Kind: "gate", Gate: g})
			// TODO: fork, join
		}
	}
	return steps, nil
}

func parseTask(node *Sexpr) (*ast.Task, error) {
	kmap := parseKeywordMap(node.List.Elements[1:])
	task := &ast.Task{
		ID: atomText(kmap[":id"]),
		On: atomText(kmap[":on"]),
		Op: atomText(kmap[":op"]),
	}
	for _, subEl := range node.List.Elements[1:] {
		if subEl.List != nil && atomText(subEl.List.Elements[0]) == "args" {
			task.Args, _ = parseConfig(subEl.List.Elements[1:])
			break
		}
	}
	return task, nil
}

func parseGate(node *Sexpr) (*ast.Gate, error) {
	kmap := parseKeywordMap(node.List.Elements[1:])
	gate := &ast.Gate{ID: atomText(kmap[":id"])}
	for _, subEl := range node.List.Elements[1:] {
		if subEl.List != nil && atomText(subEl.List.Elements[0]) == "when" {
			gate.Condition = atomText(subEl.List.Elements[1])
			break
		}
	}
	return gate, nil
}

/* ---------------- helpers ---------------- */

// atomText returns the string content of an atom, or ""
func atomText(n *Sexpr) string {
	if n == nil || n.Atom == nil {
		return ""
	}
	if n.Atom.String != nil {
		return *n.Atom.String
	}
	if n.Atom.Sym != nil {
		return *n.Atom.Sym
	}
	if n.Atom.Number != nil {
		return *n.Atom.Number
	}
	return ""
}

// atomValue converts an atom to interface{} (string, uint64, bool)
func atomValue(n *Sexpr) interface{} {
	if n == nil || n.Atom == nil {
		return nil
	}
	if n.Atom.String != nil {
		return *n.Atom.String
	}
	if n.Atom.Number != nil {
		// Try uint, then float, then string
		if v, err := strconv.ParseUint(*n.Atom.Number, 10, 64); err == nil {
			return v
		}
		if v, err := strconv.ParseFloat(*n.Atom.Number, 64); err == nil {
			return v
		}
		return *n.Atom.Number
	}
	if n.Atom.Sym != nil {
		s := *n.Atom.Sym
		if s == "true" {
			return true
		}
		if s == "false" {
			return false
		}
		return s // It's an identifier like 'draft'
	}
	return nil
}

// parseKeywordMap turns a list like (:id "foo" :type "bar" (...))
// into a map {":id": Sexpr<"foo">, ":type": Sexpr<"bar">}
func parseKeywordMap(list []*Sexpr) map[string]*Sexpr {
	m := make(map[string]*Sexpr)
	for i := 0; i < len(list)-1; i++ {
		keyNode := list[i]
		if keyNode.Atom == nil || keyNode.Atom.Sym == nil {
			continue
		}
		key := *keyNode.Atom.Sym
		if key[0] == ':' {
			m[key] = list[i+1]
			i++ // skip value
		}
	}
	return m
}
