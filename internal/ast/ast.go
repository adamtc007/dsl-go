package ast

import (
	"time"

	"github.com/alecthomas/participle/v2/lexer"
)

type Request struct {
	Pos lexer.Position

	Meta         *Meta         `parser:"'(' 'onboarding-request' @@"`
	Orchestrator *Orchestrator `parser:"@@"`
	Catalog      *Catalog      `parser:"@@? ')'"`
}

type Meta struct {
	Pos lexer.Position

	RequestID string    `parser:"'(' ':meta' '(' 'request-id' @String ')'"`
	Version   uint64    `parser:"'(' 'version' @Int ')'"`
	CreatedAt time.Time `parser:"('(' 'created-at' @String ')')?"`
	UpdatedAt time.Time `parser:"('(' 'updated-at' @String ')')? ')'"`
}

type Orchestrator struct {
	Pos lexer.Position

	Lifecycle *Lifecycle  `parser:"'(' ':orchestrator' @@"`
	Entities  []*Entity   `parser:"@@*"`
	Resources []*Resource `parser:"@@*"`
	Flows     []*Flow     `parser:"@@*"`
	Policies  []*Policy   `parser:"@@* ')'"`
}

type Lifecycle struct {
	Pos lexer.Position

	States      []string      `parser:"'(' ':lifecycle' '(' 'states' @Ident* ')'"`
	Initial     string        `parser:"'(' 'initial' @Ident ')'"`
	Transitions []*Transition `parser:"'(' 'transitions' @@* ')' ')'"`
}

type Transition struct {
	Pos lexer.Position

	From    string        `parser:"'(' '->' @Ident"`
	To      string        `parser:"@Ident"`
	Guard   *Expr         `parser:"@@?"`
	Effects []*ActionCall `parser:"'(' 'do' @@* ')'? ')'"`
}

type ActionCall struct {
	Pos lexer.Position

	Name string    `parser:"'(' @Ident"`
	Args []*KVPair `parser:"@@* ')'"`
}

type Entity struct {
	Pos lexer.Position

	ID    string     `parser:"'(' ':entities' '(' 'entity' ':id' @String"`
	Typ   string     `parser:"':type' @Ident"`
	Attrs []*AttrVal `parser:"'(' 'attrs' @@* ')' ')' ')'"`
}

type AttrVal struct {
	Pos lexer.Position

	Key        string   `parser:"'(' @Ident"`
	Value      *Value   `parser:"@@"`
	Provenance *string  `parser:"(':provenance' @String)?"`
	NeededBy   []string `parser:"(':needed-by' '(' @Ident* ')')? ')'"`
}

type Resource struct {
	Pos lexer.Position

	ID       string         `parser:"'(' ':resources' '(' 'resource' ':id' @String"`
	Typ      string         `parser:"':type' @Ident"`
	Requires []*RequireItem `parser:"('(' 'requires' @@* ')')?"`
	Config   []*KVPair      `parser:"('(' 'config' @@* ')')? ')' ')'"`
}

type RequireItem struct {
	Pos lexer.Position

	Kind string `parser:"'(' @'entity'"`
	ID   string `parser:"@String ')'"`
}

type Flow struct {
	Pos lexer.Position

	ID    string  `parser:"'(' ':flows' '(' 'flow' ':id' @String"`
	Doc   *string `parser:"(@String)?"`
	Steps []*Step `parser:"'(' 'steps' @@* ')' ')' ')'"`
}

type Step struct {
	Pos lexer.Position

	Task *Task `parser:"( @@"`
	Gate *Gate `parser:"| @@"`
	Fork *Fork `parser:"| @@"`
	Join *Join `parser:"| @@ )"`
}

type Task struct {
	Pos lexer.Position

	ID       string    `parser:"'task' ':id' @String"`
	On       string    `parser:"':on' @String"`
	Op       string    `parser:"':op' @Ident"`
	Args     []*KVPair `parser:"'(' 'args' @@* ')'"`
	Needs    []string  `parser:"('(' 'needs' @String* ')')?"`
	Produces []string  `parser:"('(' 'produces' @String* ')')?"`
	Labels   []string  `parser:"('(' 'labels' @Ident* ')')?"`
}

type Gate struct {
	Pos lexer.Position

	ID        string `parser:"'gate' ':id' @String"`
	Condition string `parser:"'(' 'when' @String ')'"`
}

type Fork struct {
	Pos lexer.Position

	ID       string   `parser:"'fork' ':id' @String"`
	Branches []string `parser:"'(' 'branches' @String* ')'"`
}

type Join struct {
	Pos lexer.Position

	ID    string   `parser:"'join' ':id' @String"`
	After []string `parser:"'(' 'after' @String* ')'"`
}

type Policy struct {
	Pos lexer.Position

	Name string    `parser:"'(' ':policies' '(' 'policy' @Ident"`
	KV   []*KVPair `parser:"@@* ')' ')'"`
}

type Catalog struct {
	Pos lexer.Position

	Attributes []*AttrDef   `parser:"'(' ':catalog' '(' ':attributes' @@* ')'"`
	Actions    []*ActionDef `parser:"'(' ':actions' @@* ')' ')'"`
}

type AttrDef struct {
	Pos lexer.Position

	Name   string   `parser:"'(' @Ident"`
	Typ    string   `parser:"':' 'type' @Ident"`
	Enum   []string `parser:"(':' 'enum' '(' @Ident* ')')?"`
	Format *string  `parser:"(':' 'format' @Ident)?"`
	PII    *bool    `parser:"(':' 'pii' @('true' | 'false'))? ')'"`
}

type ActionDef struct {
	Pos lexer.Position

	Name     string      `parser:"'(' @Ident"`
	Params   []*ParamDef `parser:"'(' 'params' @@* ')'"`
	Needs    []string    `parser:"'(' 'needs' @String* ')'"`
	Produces []string    `parser:"'(' 'produces' @String* ')' ')'"`
}

type ParamDef struct {
	Pos lexer.Position

	Name     string   `parser:"'(' @Ident"`
	Typ      string   `parser:"':' 'type' @Ident"`
	Required *bool    `parser:"(':' 'required' @('true' | 'false'))?"`
	Enum     []string `parser:"(':' 'enum' '(' @Ident* ')')? ')'"`
}

type Expr struct {
	Pos lexer.Position

	Kind string `parser:"@Ident"`
	Path string `parser:"@String?"`
}

type KVPair struct {
	Pos lexer.Position

	Key   string `parser:"@Ident"`
	Value *Value `parser:"@@"`
}

type Value struct {
	Pos lexer.Position

	String *string  `parser:"@String"`
	Int    *int64   `parser:"| @Int"`
	Float  *float64 `parser:"| @Float"`
	Bool   *bool    `parser:"| @('true' | 'false')"`
	Symbol *string  `parser:"| @Ident"`
}
