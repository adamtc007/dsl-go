# Mock Data Quick Start Guide

**TL;DR**: Use JSON files instead of a database for fast development iteration.

## What is this?

The `data-mocks/` directory contains pre-built JSON files representing:
- **Entities**: Legal entities (investment managers, pension funds, SICAVs, etc.)
- **Products**: Services being offered (custody, fund accounting, etc.)
- **Scenarios**: Complete onboarding workflows with multiple entities and products

## Why use mock data?

✅ **No database setup** - Edit JSON files directly  
✅ **Fast iteration** - See changes immediately  
✅ **Version controlled** - Track changes in git  
✅ **Reproducible** - Same data every time  
✅ **Offline work** - No network dependencies  

## Quick Examples

### Example 1: Use an existing scenario

```go
package main

import (
    "fmt"
    "github.com/example/dsl-go/internal/mocks"
    "github.com/example/dsl-go/internal/generator"
)

func main() {
    loader := mocks.NewDefaultLoader()
    
    // Load pre-built scenario
    scenario, _ := loader.LoadScenario("institutional-onboarding-001.json")
    
    // Generate DSL
    gen := generator.New()
    response, _ := gen.Generate(scenario)
    
    fmt.Println(response.DSL)
}
```

### Example 2: Mix and match entities/products

```go
loader := mocks.NewDefaultLoader()

// Pick specific entities and products
scenario, _ := loader.BuildCustomScenario(
    "my-custom-onboard-001",
    []string{
        "investment-manager-001.json",
        "sicav-001.json",
    },
    []string{
        "custody-safekeeping-eur.json",
    },
)

gen := generator.New()
response, _ := gen.Generate(scenario)
```

### Example 3: Load all entities of a specific role

```go
loader := mocks.NewDefaultLoader()

// Get all investment managers
investmentMgrs, _ := loader.LoadEntitiesByRole(generator.RoleInvestmentManager)

// Get all SICAVs
sicavs, _ := loader.LoadEntitiesByRole(generator.RoleSicav)
```

## Available Mock Data

Run this to see what's available:

```go
loader := mocks.NewDefaultLoader()

entities, _ := loader.ListEntities()
products, _ := loader.ListProducts()
scenarios, _ := loader.ListScenarios()

fmt.Println("Entities:", entities)
fmt.Println("Products:", products)
fmt.Println("Scenarios:", scenarios)
```

Current files:
- **Entities**: `investment-manager-001.json`, `asset-owner-pension-001.json`, `sicav-001.json`
- **Products**: `custody-safekeeping-eur.json`
- **Scenarios**: `institutional-onboarding-001.json`

## Creating Your Own Mock Data

### Add a new entity:

Create `data-mocks/entities/my-entity.json`:

```json
{
  "id": "le:my-fund-001",
  "name": "My Fund Management Ltd",
  "role": "investment-manager",
  "entity_type": "LegalEntity",
  "lei": "1234567890ABCDEF",
  "country": "LU",
  "attributes": {
    "aum": "1000000000",
    "regulated": true,
    "risk_tier": "medium"
  }
}
```

Then load it:

```go
entity, _ := loader.LoadEntity("my-entity.json")
```

### Add a new product:

Create `data-mocks/products/my-product.json`:

```json
{
  "id": "prod:my-service",
  "product_type": "FundAccounting",
  "config": {
    "frequency": "daily",
    "nav_calculation": true
  }
}
```

### Save programmatically:

```go
newEntity := &generator.ClientEntity{
    ID:   "le:example-001",
    Name: "Example Corp",
    Role: generator.RoleInvestmentManager,
    // ... more fields
}

loader.SaveEntity(newEntity, "example-001.json")
```

## Complete Workflow

```go
// 1. Load mock data
loader := mocks.NewDefaultLoader()
scenario, _ := loader.LoadScenario("institutional-onboarding-001.json")

// 2. Generate DSL
gen := generator.New()
response, _ := gen.Generate(scenario)

// 3. Save to storage
mgr := manager.New(manager.Config{
    DataDir:     "./data",
    RegistryDir: "./registry",
})
version, hash, _ := mgr.CreateRequest(response.RequestID, response.DSL)

fmt.Printf("Created: %s v%d (hash: %s)\n", response.RequestID, version, hash)

// 4. View it
v, dsl, _ := mgr.GetCurrentText(response.RequestID)
fmt.Println(dsl)
```

## Tips

- **Start with scenarios**: Use `institutional-onboarding-001.json` as a template
- **Copy and modify**: Duplicate existing files and change the data
- **Keep IDs unique**: Use prefixes like `le:` for entities, `prod:` for products
- **Realistic data**: Include all attributes for proper testing
- **Commit to git**: Share mock data with your team

## Next Steps

1. Look at existing files in `data-mocks/entities/` to understand the structure
2. Create your own entity JSON file
3. Build a custom scenario with your entities
4. Generate DSL and see the S-expression output
5. Validate, compile, and test execution

## See Also

- `data-mocks/README.md` - Detailed documentation
- `examples/use-mocks.go` - Comprehensive examples
- `internal/mocks/loader.go` - API documentation