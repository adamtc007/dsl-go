package parse

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/example/dsl-go/internal/ast"
)

var sexprLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Whitespace", Pattern: `\s+`},
	{Name: "Comment", Pattern: `\;[^\n]*`},
	{Name: "LParen", Pattern: `\(`},
	{Name: "RParen", Pattern: `\)`},
	{Name: "Arrow", Pattern: `->`},
	{Name: "String", Pattern: `"(?:\\.|[^\"])*"`},
	{Name: "ColonIdent", Pattern: `:[A-Za-z][A-Za-z0-9_-]*`},
	{Name: "Ident", Pattern: `[A-Za-z][A-Za-z0-9_-]*`},
	{Name: "Number", Pattern: `[0-9]+(?:\.[0-9]+)?`}, // Add number support
})

// Parser interface
type Parser interface {
	Parse(text string) (*ast.Request, error)
}

// ParticipleParser is a parser that uses participle
type ParticipleParser struct {
	parser *participle.Parser[ast.Request]
}

// New creates a new participle parser
func New() (Parser, error) {
	parser, err := participle.Build[ast.Request](
		participle.Lexer(sexprLexer),
		participle.Unquote("String"),
		participle.Elide("Whitespace", "Comment"),
	)
	if err != nil {
		return nil, err
	}
	return &ParticipleParser{parser: parser}, nil
}

// Parse parses the given text into an AST
func (p *ParticipleParser) Parse(text string) (*ast.Request, error) {
	return p.parser.ParseString("", text)
}
