# Mock Data for DSL Generation

This directory contains JSON mock data for rapid development and testing without requiring a database connection.

## Directory Structure

```
data-mocks/
├── entities/           # Legal entities with various roles
├── products/          # Product specifications
├── scenarios/         # Complete onboarding scenarios
└── README.md
```

## Entity Roles

Entities can have the following roles:

- **investment-manager** - Investment management firms
- **asset-owner** - Pension funds, endowments, family offices
- **management-company** - Fund management companies
- **sicav** - SICAV (Investment Company with Variable Capital)
- **custodian** - Custodian banks
- **administrator** - Fund administrators

## Available Mock Data

### Entities

- `investment-manager-001.json` - Luxembourg-based investment manager
- `asset-owner-pension-001.json` - US pension fund
- `sicav-001.json` - Luxembourg SICAV with multiple sub-funds

### Products

- `custody-safekeeping-eur.json` - EUR custody and safekeeping service

### Scenarios

- `institutional-onboarding-001.json` - Complex multi-entity institutional onboarding

## Usage

### Using the CLI

#### List available mock data:
```bash
# List all entities
./bin/dsl-go list-mocks entities

# List all products
./bin/dsl-go list-mocks products

# List all scenarios
./bin/dsl-go list-mocks scenarios
```

#### Generate from a scenario:
```bash
./bin/dsl-go generate-from-scenario institutional-onboarding-001.json
```

#### Build custom combination:
```bash
./bin/dsl-go generate-from-mocks \
  --request-id onboard-custom-001 \
  --entities investment-manager-001.json,sicav-001.json \
  --products custody-safekeeping-eur.json
```

### Using the Go API

```go
import "github.com/example/dsl-go/internal/mocks"

// Load a single entity
loader := mocks.NewDefaultLoader()
entity, err := loader.LoadEntity("investment-manager-001.json")

// Load a complete scenario
scenario, err := loader.LoadScenario("institutional-onboarding-001.json")

// Build a custom scenario
customScenario, err := loader.BuildCustomScenario(
    "onboard-req-001",
    []string{"investment-manager-001.json", "sicav-001.json"},
    []string{"custody-safekeeping-eur.json"},
)

// Generate DSL from scenario
gen := generator.New()
response, err := gen.Generate(scenario)
```

## Creating New Mock Data

### Entity Template

```json
{
  "id": "le:entity-id",
  "name": "Legal Name",
  "role": "investment-manager|asset-owner|management-company|sicav|custodian|administrator",
  "entity_type": "LegalEntity",
  "lei": "LEI_CODE_IF_AVAILABLE",
  "country": "ISO_COUNTRY_CODE",
  "attributes": {
    "registration_number": "...",
    "regulated": true,
    "regulator": "...",
    "risk_tier": "low|medium|high",
    "kyc_level": "standard|enhanced",
    ...
  }
}
```

### Product Template

```json
{
  "id": "prod:product-id",
  "product_type": "CustodySafekeeping|FundAccounting|FundAdministration|...",
  "currency": "EUR|USD|...",
  "config": {
    "account_type": "...",
    ...
  }
}
```

### Scenario Template

```json
{
  "request_id": "unique-request-id",
  "tenant_id": "tenant-identifier",
  "entities": [ /* array of entity objects */ ],
  "products": [ /* array of product objects */ ],
  "resources": [ /* optional additional resources */ ],
  "metadata": { /* optional metadata */ }
}
```

## Benefits

### Fast Iteration
- No database setup required
- Edit JSON files directly
- See changes immediately
- Version control friendly

### Testing
- Consistent test data
- Easy to create edge cases
- Reproducible scenarios
- No test data pollution

### Development
- Work offline
- No dependencies on external systems
- Prototype new entity types
- Experiment with product combinations

## Best Practices

1. **Naming Convention**: Use descriptive filenames like `investment-manager-001.json`
2. **IDs**: Use consistent ID prefixes (`le:` for entities, `prod:` for products)
3. **Completeness**: Include all relevant attributes for realistic testing
4. **Documentation**: Add comments in attributes for complex fields
5. **Validation**: Test your mock data through the generator before committing

## Example Workflow

1. **Create entity mocks** for each role in your onboarding
2. **Create product mocks** for services being offered
3. **Combine into scenario** for complete onboarding flow
4. **Generate DSL** to see the populated S-expression
5. **Validate** the generated DSL
6. **Compile** to execution plan
7. **Refine** mock data based on results

## Extending Mock Data

To add new entity types or product types:

1. Create JSON file in appropriate directory
2. Follow the template structure above
3. Test with generator: `./bin/dsl-go generate-from-json your-file.json`
4. Adjust attributes as needed
5. Commit to repository

## Integration with Agent

These mocks serve as:
- **Examples** for LLM context (RAG)
- **Templates** for agent-generated DSL
- **Test cases** for validation
- **Reference data** for entity/product attributes

The agent can use these mocks to learn what realistic onboarding data looks like and generate similar structures for new clients.