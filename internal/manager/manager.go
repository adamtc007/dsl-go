package manager

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/example/dsl-go/internal/ast"
	"github.com/example/dsl-go/internal/parse"
	"github.com/example/dsl-go/internal/print"
	"github.com/example/dsl-go/internal/storage"
)

type Config struct {
	RegistryDir string
	DataDir     string
}

type Manager struct {
	store          *storage.FileStore
	parser         parse.Parser
	cfg            Config
	dataDictionary *DataDictionary
}

func New(cfg Config) (*Manager, error) {
	parser, err := parse.New()
	if err != nil {
		return nil, err
	}
	m := &Manager{
		store:  storage.NewFileStore(cfg.DataDir),
		parser: parser,
		cfg:    cfg,
	}
	if err := m.LoadDataDictionary(); err != nil {
		// For now, we'll just log the error. In a real application, you might want to handle this more gracefully.
		fmt.Printf("warning: could not load data dictionary: %v\n", err)
	}
	return m, nil
}

func (m *Manager) LoadDataDictionary() error {
	path := filepath.Join(m.cfg.RegistryDir, "data-dictionary.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read data dictionary: %w", err)
	}

	var dict DataDictionary
	if err := json.Unmarshal(data, &dict); err != nil {
		return fmt.Errorf("failed to parse data dictionary: %w", err)
	}

	m.dataDictionary = &dict

	return nil
}

func (m *Manager) GetDataDictionary() *DataDictionary {
	return m.dataDictionary
}

func (m *Manager) GetAttribute(id string) (Attribute, bool) {
	for _, attr := range m.dataDictionary.Attributes {
		if attr.AttributeID == id {
			return attr, true
		}
	}
	return Attribute{}, false
}

func (m *Manager) CreateRequest(id string, template string) (version uint64, canonicalHash string, err error) {
	req, err := m.parser.Parse(template) // strict
	if err != nil {
		return 0, "", err
	}

	now := time.Now().UTC()
	if req.Meta == nil {
		req.Meta = &ast.Meta{}
	}
	req.Meta.RequestID = id
	req.Meta.Version = 1
	if req.Meta.CreatedAt.IsZero() {
		req.Meta.CreatedAt = now
	}
	req.Meta.UpdatedAt = now

	txt := print.ToSexpr(req)
	if err := m.store.Put(id, 1, txt); err != nil {
		return 0, "", fmt.Errorf("failed to store request: %w", err)
	}
	return 1, hash(txt), nil
}

func (m *Manager) GetCurrentText(id string) (version uint64, text string, err error) {
	return m.store.GetLatest(id)
}

func (m *Manager) ValidateText(text string) (issues []string, err error) {
	_, err = m.parser.Parse(text)
	if err != nil {
		return []string{err.Error()}, nil
	}
	return nil, nil
}

// Compile/Plan/Delta are stubs (parity with Rust baseline)
type Plan struct {
	Steps    []PlanStep `json:"steps"`
	PlanHash string     `json:"plan_hash"`
}
type PlanStep struct {
	ID     string      `json:"id"`
	Action string      `json:"action"`
	Inputs [][2]string `json:"inputs"`
	After  []string    `json:"after"`
}

func (m *Manager) CompilePlan(text string) (*Plan, error) {
	_, err := m.parser.Parse(text)
	if err != nil {
		return nil, err
	}
	return &Plan{Steps: []PlanStep{}, PlanHash: "todo"}, nil
}

type PlanDelta struct {
	Added   []PlanStep    `json:"added"`
	Removed []PlanStep    `json:"removed"`
	Changed [][2]PlanStep `json:"changed"`
}

func (m *Manager) PlanDelta(fromText, toText string) (*PlanDelta, error) {
	_, err := m.parser.Parse(fromText)
	if err != nil {
		return nil, err
	}
	_, err = m.parser.Parse(toText)
	if err != nil {
		return nil, err
	}
	return &PlanDelta{Added: nil, Removed: nil, Changed: nil}, nil
}

func hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return "sha256:" + hex.EncodeToString(h[:])
}

var ErrNotFound = errors.New("not found")

// expose AST type to CLI (for ast-json)
type Request = ast.Request
