package ast

import "time"

type Request struct {
	Meta         *Meta        `json:"meta"`
	Orchestrator Orchestrator `json:"orchestrator"`
	Catalog      *Catalog     `json:"catalog,omitempty"`
}

type Meta struct {
	RequestID string    `json:"request_id"`
	Version   uint64    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Orchestrator struct {
	Lifecycle Lifecycle         `json:"lifecycle"`
	Entities  map[string]Entity   `json:"entities"`
	Resources map[string]Resource `json:"resources"`
	Flows     map[string]Flow     `json:"flows"`
	Policies  []Policy          `json:"policies,omitempty"`
}

type Lifecycle struct {
	States      []string     `json:"states"`
	Initial     string       `json:"initial"`
	Transitions []Transition `json:"transitions"`
}

type Transition struct {
	From    string       `json:"from"`
	To      string       `json:"to"`
	Guard   *Expr        `json:"guard,omitempty"`
	Effects []ActionCall `json:"effects,omitempty"`
}

type ActionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args,omitempty"`
}

type Entity struct {
	ID    string            `json:"id"`
	Typ   string            `json:"typ"`
	Attrs map[string]AttrVal `json:"attrs"`
}

type AttrVal struct {
	Value      interface{} `json:"value"`
	Provenance *string     `json:"provenance,omitempty"`
	NeededBy   []string    `json:"needed_by,omitempty"`
}

type Resource struct {
	ID       string                 `json:"id"`
	Typ      string                 `json:"typ"`
	Requires []RequireItem          `json:"requires,omitempty"`
	Config   map[string]interface{} `json:"config,omitempty"`
}

type RequireItem struct {
	Kind string `json:"kind"` // "entity" or "attr"
	ID   string `json:"id"`   // id or refpath
}

type Flow struct {
	ID    string  `json:"id"`
	Doc   *string `json:"doc,omitempty"`
	Steps []Step  `json:"steps"`
}

type Step struct {
	Kind string `json:"kind"` // "task"|"gate"|"fork"|"join"
	Task *Task  `json:"task,omitempty"`
	Gate *Gate  `json:"gate,omitempty"`
	Fork *Fork  `json:"fork,omitempty"`
	Join *Join  `json:"join,omitempty"`
}

type Task struct {
	ID       string                 `json:"id"`
	On       string                 `json:"on"`
	Op       string                 `json:"op"`
	Args     map[string]interface{} `json:"args,omitempty"`
	Needs    []string               `json:"needs,omitempty"`
	Produces []string               `json:"produces,omitempty"`
	Labels   []string               `json:"labels,omitempty"`
}

type Gate struct {
	ID        string `json:"id"`
	Condition string `json:"condition"`
}

type Fork struct {
	ID       string   `json:"id"`
	Branches []string `json:"branches"`
}

type Join struct {
	ID    string   `json:"id"`
	After []string `json:"after"`
}

type Policy struct {
	Name string                 `json:"name"`
	KV   map[string]interface{} `json:"kv,omitempty"`
}

type Catalog struct {
	Attributes map[string]AttrDef  `json:"attributes"`
	Actions    map[string]ActionDef `json:"actions"`
}

type AttrDef struct {
	Typ    string    `json:"typ"`
	Enum   *[]string `json:"enum,omitempty"`
	Format *string   `json:"format,omitempty"`
	PII    *bool     `json:"pii,omitempty"`
}

type ActionDef struct {
	Params   map[string]ParamDef `json:"params"`
	Needs    []string            `json:"needs"`
	Produces []string            `json:"produces"`
}

type ParamDef struct {
	Typ      string    `json:"typ"`
	Required bool      `json:"required"`
	Enum     *[]string `json:"enum,omitempty"`
}

type Expr struct {
	Kind string `json:"kind"` // minimal placeholder
	Path string `json:"path,omitempty"`
}
