# Session Notes: DSL-Go Project Development

**Date**: 2024-10-29  
**Project**: dsl-go - S-expression based onboarding/KYC orchestration DSL  
**Session Focus**: Architecture discussion, mock data system, generator implementation

---

## Project Overview

### Domain: Onboarding & KYC for Financial Services
- **Target clients**: Investment Managers, Asset Owners, Pension Funds, SICAVs, Management Companies
- **Products**: Custody, Fund Accounting, Fund Administration, Regulatory Reporting
- **Use case**: Long-running, multi-step, parallel workflow orchestration

### Core Concept: Agent-Generated DSL
The DSL is NOT hand-written. Instead:
1. **Agent (LLM)** generates DSL from client data + RAG context
2. **Human reviews** the generated DSL (explainable AI via code)
3. **Compile** DSL to execution plan
4. **Execute** long-running workflow (hours to days)
5. **State = Updated DSL** (self-describing, versionable, auditable)

This is different from typical DSLs - the DSL instance IS the state machine state.

---

## Architectural Decisions

### Language Choice: Go (Not Rust)

**Rationale**:
- **I/O-bound workload**: Orchestrating external services (KYC APIs, AML screening, document verification)
- **Fast development**: Materially faster iteration for prototyping
- **Good enough performance**: Adding 1ms overhead when API calls take 500ms+ is irrelevant
- **Mature ecosystem**: Excellent for APIs, gRPC, concurrent I/O

**When Rust would be needed**:
- Custom VM/interpreter for DSL execution
- WASM compilation target
- Untrusted code sandbox
- Microsecond-level latency requirements
- No GC allowed (embedded, WASM)

**Quote from user**: "Unless GO or its libraries cannot really handle coding something - and compilers or wasm modules is a case in point then it wins"

### Recommended Hybrid Approach (Future)
```
┌─────────────────────────────────────┐
│  Go: Orchestration Layer            │
│  - Parse DSL                         │
│  - Coordinate task execution         │
│  - State management                  │
│  - API calls, I/O, networking        │
└─────────────────────────────────────┘
              │ (If needed later)
              ↓
┌─────────────────────────────────────┐
│  Rust: Hot Path Components           │
│  - Expression evaluator              │
│  - Policy engine                     │
│  - WASM runtime for extensions       │
└─────────────────────────────────────┘
```

### Scale Expectations
- **Volume**: Thousands of onboarding requests total (not per day)
- **Concurrent**: Maybe 100 cases/day at peak
- **Latency**: Minutes to days per case (not microseconds)
- **Bottleneck**: LLM generation (2-10 sec), external APIs (30 sec - 5 min)
- **Go's role**: <10ms overhead, not the bottleneck

---

## System Architecture

### Complete Workflow Pipeline

```
1. GENERATE
   Agent + RAG → DSL Instance (populated with customer data)
   - Agent queries: compliance rules, operations, examples
   - Generates: S-expr with entities, products, flows
   - Output: Deterministic, reviewable DSL

2. REVIEW
   Human reviews generated DSL in UI
   - Inspect entities, resources, flows
   - Verify compliance requirements
   - Check task sequencing

3. REFINE (optional)
   - Human edits DSL directly, OR
   - Human provides feedback → Agent regenerates

4. COMPILE
   DSL (declarative) → Execution Plan (imperative)
   - Resolve dependencies
   - Build task DAG
   - Generate gRPC service calls

5. EXECUTE
   Long-running workflow execution
   - Parallel task execution (goroutines)
   - State updates (DSL versioning)
   - Gate evaluations (manual approval)
   - Human interventions
```

### Key Architectural Insight: DSL IS State

The DSL instance is not just a template - it's the **source of truth** for execution state:
- Initial DSL: entities declared, tasks pending
- During execution: DSL updated with task results
- Each update: new version, new hash
- Terminal state: DSL reflects "onboarded" or "kyc-cleared"

**Example evolution**:
```
v1: (task :id verify-identity :status pending)
v2: (task :id verify-identity :status completed :result "passed" :completed-at "...")
v3: (task :id verify-identity :status completed) + (entity customer (attrs (kyc-status "verified")))
```

### Multi-tenancy Model
- **Shared grammar**: All tenants use same S-expr structure
- **Tenant-specific RAG**: Different compliance rules, operations, providers
- **Tenant-specific validation**: What operations are available
- **Firewalls**: Data isolation at storage + execution level

### gRPC Service Design

**One gRPC service per DSL type** (Onboarding, KYC, etc.):
```protobuf
service OnboardingOrchestrator {
  // Task execution
  rpc VerifyIdentity(VerifyIdentityRequest) returns (TaskResult);
  rpc CheckAML(AMLCheckRequest) returns (TaskResult);
  
  // State management
  rpc GetState(GetStateRequest) returns (OnboardingState);
  rpc UpdateEntity(UpdateEntityRequest) returns (UpdateResult);
  
  // Flow control
  rpc ExecuteFlow(FlowRequest) returns (stream FlowEvent);
  rpc EvaluateGate(GateRequest) returns (GateDecision);
}
```

Each task in DSL maps to gRPC call. Compiler generates service definitions from DSL.

---

## What We Built This Session

### 1. Fixed Build Errors
- **Problem**: `lexer.MustSimple` type mismatch in participle v2
- **Solution**: Changed from `[]lexer.Rule` to `[]lexer.SimpleRule` with named fields
- **Result**: Project builds cleanly

### 2. Generator Package (`internal/generator/`)

**Purpose**: Populate DSL templates with client entities and products

**Key types**:
- `ClientEntity`: Legal entities with roles (investment-manager, asset-owner, sicav, etc.)
- `ProductSpec`: Products being onboarded (custody, fund-accounting, etc.)
- `ResourceSpec`: Additional resources (KYC services, document storage, etc.)
- `GenerateRequest`: Complete specification for DSL generation
- `GenerateResponse`: Generated DSL + metadata

**API**:
```go
gen := generator.New()

// From scratch
response, err := gen.Generate(request)

// From template
response, err := gen.GenerateFromTemplate(templateDSL, request)
```

**Features**:
- Generates complete DSL structure (meta, orchestrator, lifecycle, entities, resources, flows)
- Creates verification tasks for each entity based on role
- Generates AML screening tasks
- Adds compliance review gates
- Creates product setup tasks
- Handles multiple entities and products in single onboarding

### 3. Mock Data System (`data-mocks/`)

**Problem**: Database slows down iteration cycles
**Solution**: JSON files as data source

**Directory Structure**:
```
data-mocks/
├── entities/           # Legal entity definitions
├── products/          # Product specifications
├── scenarios/         # Complete onboarding scenarios
├── README.md          # Full documentation
└── QUICKSTART.md      # 5-minute guide
```

**Mock Data Loader** (`internal/mocks/loader.go`):
```go
loader := mocks.NewDefaultLoader()

// Load individual items
entity, _ := loader.LoadEntity("investment-manager-001.json")
product, _ := loader.LoadProduct("custody-safekeeping-eur.json")
scenario, _ := loader.LoadScenario("institutional-onboarding-001.json")

// Load all of a type
entities, _ := loader.LoadAllEntities()
investmentMgrs, _ := loader.LoadEntitiesByRole(RoleInvestmentManager)

// Build custom combinations
scenario, _ := loader.BuildCustomScenario(
    "my-onboard-001",
    []string{"investment-manager-001.json", "sicav-001.json"},
    []string{"custody-safekeeping-eur.json"},
)

// Save new mocks
loader.SaveEntity(newEntity, "my-entity.json")
```

**Benefits**:
- ✅ No database setup required
- ✅ Edit JSON files directly in editor
- ✅ Version controlled (git)
- ✅ Fast iteration (instant feedback)
- ✅ Reproducible test data
- ✅ Easy to create edge cases
- ✅ Team can share scenarios

### 4. Realistic Mock Data Files

Created detailed, production-realistic mock data:

**`investment-manager-001.json`** (89 lines):
- Full company details, LEI, registration
- Regulatory info (CSSF, license, passporting)
- Business metrics (AUM, client count, strategies)
- Ownership structure with UBOs
- Risk profile (KYC level, PEP exposure, sanctions)
- Document list

**`asset-owner-pension-001.json`** (137 lines):
- US pension fund (ERISA, DOL regulated)
- Participant counts, funded ratio
- Board of trustees, governance
- Investment policy, allocation targets
- Tax status (FATCA, CRS)
- Operational details

**`sicav-001.json`** (285 lines):
- UCITS-compliant umbrella SICAV
- 4 sub-funds with ISINs
- Multiple share classes
- Service providers (depositary, administrator, auditor)
- Board of directors, conducting officers
- Fees, risk profile, SFDR classification
- Operational details (NAV calculation, dealing)

**`custody-safekeeping-eur.json`** (128 lines):
- Account types, settlement methods
- Markets and asset classes supported
- Corporate actions processing
- Pricing and valuation
- Reporting deliverables
- Fees, SLAs
- Regulatory compliance (CSDR, MiFID II, EMIR)

**`institutional-onboarding-001.json`** (265 lines):
- Complete scenario with 4 entities
- 6 products
- 4 resources (KYC verification, AML screening, document vault, LEI verification)
- Metadata (relationship manager, revenue estimates, strategic importance)
- Workflow preferences

---

## Key Technical Insights

### Go vs Rust Trade-offs

**Developer Experience**:
- Go: "Materially faster for me" (user quote)
- Rust: "Compiler outputs and clippy are a massive help. If it compiles, it runs." (user quote)

**When to use each**:
- Go: Default choice for services, APIs, orchestration (this project)
- Rust: When Go literally can't do it (WASM, compilers, no-GC requirements)

**User's philosophy**: "Unless GO or its libraries cannot really handle coding something...then it wins"

### Performance Reality Check

For this use case:
- LLM generation: **2-10 seconds** (unavoidable)
- External API calls: **30 sec - 5 min** (KYC providers, AML screening)
- Go parsing/validation: **<10ms** (negligible)
- Database I/O: **10-100ms** (can optimize with caching)

**Conclusion**: Go adding 1ms overhead when total flow takes hours/days is irrelevant.

### Agent Integration Points

Mock data serves multiple purposes:
1. **RAG examples**: Train agent on realistic entity/product structures
2. **Templates**: Agent learns patterns from successful onboardings
3. **Validation**: Generated DSL compared to known-good examples
4. **Testing**: Reproducible test cases for DSL generation
5. **Reference data**: Attribute schemas, required fields, valid values

---

## Next Steps / Future Work

### Immediate (In Progress)
1. Fix corrupted `cmd/dsl-go/main.go` file
2. Add CLI commands for mock data:
   - `dsl-go list-mocks entities|products|scenarios`
   - `dsl-go generate-from-scenario <file>`
   - `dsl-go generate-from-mocks --entities X --products Y`

### Short Term
1. Complete compilation pipeline (DSL → Execution Plan)
2. Build executor service (gRPC-based)
3. State persistence layer (versioned DSL storage)
4. Basic UI for DSL review and visualization

### Medium Term
1. Agent integration (LLM-based generation)
2. RAG system with mock data as corpus
3. Grammar-constrained generation
4. Human review workflows with feedback loop

### Long Term (Optimization Phase)
1. Profile actual performance bottlenecks
2. Consider Rust for hot paths IF needed (likely not)
3. Distributed execution with Temporal or similar
4. Multi-region deployment
5. Advanced features (securities lending, complex workflows)

---

## Important Context for Future Sessions

### Project Goals
- **NOT building a traditional DSL**: Building an agent-driven workflow generation system
- **DSL is the artifact**: Deterministic, reviewable, auditable output from AI
- **Explainable AI via code**: Compliance officers review S-expressions, not model weights
- **State = DSL**: The DSL instance evolves as the workflow executes

### Design Constraints
- **Multi-tenant**: Shared infrastructure, isolated data
- **Long-running**: Workflows span hours to days
- **Human-in-loop**: Manual review gates, approval workflows
- **Audit trail**: Every DSL version is hashed and stored
- **Resumability**: Can restart execution from last saved state

### User Preferences
- **Pragmatic over perfect**: Ship features fast, optimize if needed
- **Go unless impossible**: Only reach for Rust when Go can't do it
- **Mock data over DB**: Fast iteration is more valuable than production-like setup
- **Agent-generated**: Templates are filled by AI, not humans

### Testing Philosophy
- Use mock data for rapid iteration
- Don't worry about performance until it's actually slow
- Validate generated DSL against known-good examples
- Test edge cases with custom JSON scenarios

---

## File Locations Reference

### Core Implementation
- `internal/generator/types.go` - Data structures for entities, products, resources
- `internal/generator/generator.go` - DSL generation logic
- `internal/mocks/loader.go` - Mock data loader API
- `internal/ast/ast.go` - AST structures for parsed DSL
- `internal/parse/sexpr.go` - S-expression parser (participle-based)
- `internal/print/print.go` - AST → S-expression serializer
- `internal/manager/manager.go` - Request lifecycle management

### Mock Data
- `data-mocks/entities/*.json` - Entity definitions
- `data-mocks/products/*.json` - Product specifications
- `data-mocks/scenarios/*.json` - Complete scenarios
- `data-mocks/README.md` - Full documentation
- `data-mocks/QUICKSTART.md` - Quick start guide

### Examples
- `examples/full.sexpr` - Complete DSL example
- `examples/minimal.sexpr` - Minimal DSL example
- `examples/use-mocks.go` - Mock data usage examples
- `examples/onboarding-input.json` - JSON input format example

### CLI
- `cmd/dsl-go/main.go` - Main CLI (has corruption, needs fix)

---

## Agent Prompt Context

When using this in future sessions, the agent should understand:

1. **This is NOT a traditional compiler project** - It's an AI-generated workflow system
2. **The DSL is generated, not written** - Templates are filled by LLM with RAG
3. **Performance is not the concern** - I/O-bound workload, orchestration not computation
4. **Mock data is the development strategy** - Fast iteration over production-like setup
5. **Go is the right choice** - Only Rust if absolutely necessary
6. **Scale is moderate** - Thousands of requests total, not millions per second

The user values **pragmatic engineering**: ship fast, optimize when there's evidence it's needed, use boring technology that works.

---

## Questions to Ask User in Future Sessions

- Has the agent integration strategy evolved?
- What RAG system are you using for the LLM?
- Have you hit any Go performance limitations?
- How is the human review workflow being implemented?
- What's the multi-tenancy isolation strategy (DB-level? Service-level?)
- Are you using Temporal or building custom orchestration?