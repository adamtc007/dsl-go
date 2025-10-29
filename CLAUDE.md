# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go implementation of a Domain Specific Language (DSL) for platform-level orchestration. The project parses S-expressions that describe onboarding requests with entities, resources, and workflows.

## Architecture

The codebase follows a modular architecture with these key components:

- **Parser (`internal/parse/`)**: Uses Participle library to parse S-expressions into a generic tree structure, then maps to AST
- **AST (`internal/ast/`)**: Defines the core data structures for requests, orchestrators, entities, resources, and flows
- **Manager (`internal/manager/`)**: Provides high-level operations for creating, validating, and compiling requests
- **Storage (`internal/storage/`)**: File-based storage system for persisting requests with versioning
- **Print (`internal/print/`)**: Converts AST back to S-expression format
- **CLI (`cmd/dsl-go/`)**: Command-line interface with multiple operations

## Common Commands

### Building
```bash
go build -o dsl-go cmd/dsl-go/main.go
```

### Running the CLI
The built binary supports these commands:
- `./dsl-go create <request_id> <template.sexpr>` - Create a new request from S-expression file
- `./dsl-go show <request_id>` - Display current version of a request
- `./dsl-go validate <file.sexpr>` - Validate S-expression syntax
- `./dsl-go compile <file.sexpr>` - Compile to execution plan (stub implementation)
- `./dsl-go plan-delta <from.sexpr> <to.sexpr>` - Compare two versions (stub implementation)
- `./dsl-go ebnf` - Display grammar specification
- `./dsl-go parse-summary <file.sexpr>` - Show parsed structure summary
- `./dsl-go ast-json <file.sexpr>` - Output AST as JSON

### Testing
```bash
go test ./...
```
Note: Currently no test files exist in the codebase.

## Development Notes

- The project uses the Participle parser generator library for S-expression parsing
- Example files are available in `examples/` directory
- Grammar specification is documented in `docs/ebnf_v0_1.txt`
- The manager provides versioning and hashing for request storage in `./data/` directory
- Registry functionality is stubbed but directory expected at `./registry/`

## Data Flow

1. S-expression input → Participle parser → Generic Sexpr tree
2. Sexpr tree → AST mapping → Typed Request structure
3. Request → Manager operations (validate/compile/store)
4. Storage uses SHA256 hashing for content verification