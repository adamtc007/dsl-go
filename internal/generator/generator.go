package generator

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/example/dsl-go/internal/ast"
	"github.com/example/dsl-go/internal/parse"
	"github.com/example/dsl-go/internal/print"
)

// Generator generates populated DSL instances from templates and client data
type Generator struct {
	parser parse.Parser
}

// New creates a new Generator instance
func New() (*Generator, error) {
	parser, err := parse.New()
	if err != nil {
		return nil, err
	}
	return &Generator{
		parser: parser,
	}, nil
}

// Generate creates a populated DSL instance from the request
func (g *Generator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
	if err := g.validate(req); err != nil {
		return nil, err
	}

	// Create base request structure
	dslRequest := g.createBaseRequest(req)

	// Add client entities
	g.addEntities(dslRequest, req.Entities)

	// Add products as resources
	g.addResources(dslRequest, req.Products, req.Resources)

	// Generate onboarding flows
	g.generateFlows(dslRequest)

	// Convert to S-expression format
	dslText := print.ToSexpr(dslRequest)

	// Prepare response
	response := &GenerateResponse{
		RequestID:      req.RequestID,
		DSL:            dslText,
		Version:        1,
		Hash:           "",
		GeneratedAt:    time.Now().UTC(),
		EntitiesAdded:  len(req.Entities),
		ResourcesAdded: len(req.Products) + len(req.Resources),
		FlowsGenerated: 1, // main flow
	}

	return response, nil
}

// GenerateFromTemplate generates a DSL instance from an existing template
func (g *Generator) GenerateFromTemplate(templateDSL string, req *GenerateRequest) (*GenerateResponse, error) {
	if err := g.validate(req); err != nil {
		return nil, err
	}

	// Parse the template
	dslRequest, err := g.parser.Parse(templateDSL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Update metadata
	if dslRequest.Meta == nil {
		dslRequest.Meta = &ast.Meta{}
	}
	dslRequest.Meta.RequestID = req.RequestID
	dslRequest.Meta.Version = 1
	now := time.Now().UTC()
	dslRequest.Meta.CreatedAt = now
	dslRequest.Meta.UpdatedAt = now

	// Add client entities
	g.addEntities(dslRequest, req.Entities)

	// Add products and resources
	g.addResources(dslRequest, req.Products, req.Resources)

	// Convert to S-expression format
	dslText := print.ToSexpr(dslRequest)

	response := &GenerateResponse{
		RequestID:      req.RequestID,
		DSL:            dslText,
		Version:        1,
		GeneratedAt:    time.Now().UTC(),
		EntitiesAdded:  len(req.Entities),
		ResourcesAdded: len(req.Products) + len(req.Resources),
		FlowsGenerated: len(dslRequest.Orchestrator.Flows),
	}

	return response, nil
}

func (g *Generator) GenerateFromTemplateFile(templatePath string, req *GenerateRequest) (*GenerateResponse, error) {
	if err := g.validate(req); err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template file: %w", err)
	}

	req.Now = time.Now()

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, req); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	dslText := buf.String()

	response := &GenerateResponse{
		RequestID:      req.RequestID,
		DSL:            dslText,
		Version:        1,
		GeneratedAt:    req.Now.UTC(),
		EntitiesAdded:  len(req.Entities),
		ResourcesAdded: len(req.Products) + len(req.Resources),
		FlowsGenerated: 1, // This is now controlled by the template
	}

	return response, nil
}

// validate checks that the GenerateRequest has required fields
func (g *Generator) validate(req *GenerateRequest) error {
	if req.RequestID == "" {
		return &ValidationError{Field: "RequestID", Message: "required"}
	}
	if len(req.Entities) == 0 {
		return &ValidationError{Field: "Entities", Message: "at least one entity required"}
	}
	return nil
}

// createBaseRequest creates a minimal DSL request structure
func (g *Generator) createBaseRequest(req *GenerateRequest) *ast.Request {
	now := time.Now().UTC()

	return &ast.Request{
		Meta: &ast.Meta{
			RequestID: req.RequestID,
			Version:   1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Orchestrator: &ast.Orchestrator{
			Lifecycle: &ast.Lifecycle{
				States:      []string{"draft", "validated", "in-progress", "kyc-complete", "onboarded", "failed"},
				Initial:     "draft",
				Transitions: []*ast.Transition{},
			},
			Entities:  []*ast.Entity{},
			Resources: []*ast.Resource{},
			Flows:     []*ast.Flow{},
		},
	}
}

// addEntities adds client entities to the DSL
func (g *Generator) addEntities(dslReq *ast.Request, entities []ClientEntity) {
	for _, clientEntity := range entities {
		attrs := []*ast.AttrVal{}

		// Add standard attributes
		attrs = append(attrs, &ast.AttrVal{
			Key:        "name",
			Value:      &ast.Value{String: &clientEntity.Name},
			Provenance: stringPtr("client-provided"),
		})
		attrs = append(attrs, &ast.AttrVal{
			Key:        "role",
			Value:      &ast.Value{Symbol: stringPtr(string(clientEntity.Role))},
			Provenance: stringPtr("system-assigned"),
		})

		if clientEntity.Country != "" {
			attrs = append(attrs, &ast.AttrVal{
				Key:        "country",
				Value:      &ast.Value{String: &clientEntity.Country},
				Provenance: stringPtr("client-provided"),
			})
		}

		if clientEntity.LEI != "" {
			attrs = append(attrs, &ast.AttrVal{
				Key:        "lei",
				Value:      &ast.Value{String: &clientEntity.LEI},
				Provenance: stringPtr("client-provided"),
			})
		}

		// Add any additional attributes
		for key, value := range clientEntity.Attributes {
			strVal, ok := value.(string)
			if !ok {
				continue
			}
			attrs = append(attrs, &ast.AttrVal{
				Key:        key,
				Value:      &ast.Value{String: &strVal},
				Provenance: stringPtr("client-provided"),
			})
		}

		entity := &ast.Entity{
			ID:    clientEntity.ID,
			Typ:   clientEntity.EntityType,
			Attrs: attrs,
		}

		dslReq.Orchestrator.Entities = append(dslReq.Orchestrator.Entities, entity)
	}
}

// addResources adds products and resources to the DSL
func (g *Generator) addResources(dslReq *ast.Request, products []ProductSpec, resources []ResourceSpec) {
	// Add products as resources
	for _, product := range products {
		requires := []*ast.RequireItem{}
		// Products typically require at least one entity
		if len(dslReq.Orchestrator.Entities) > 0 {
			for _, entity := range dslReq.Orchestrator.Entities {
				requires = append(requires, &ast.RequireItem{
					Kind: "entity",
					ID:   entity.ID,
				})
				break // Just require the first entity for now
			}
		}

		config := []*ast.KVPair{}
		if product.Currency != "" {
			config = append(config, &ast.KVPair{
				Key:   "currency",
				Value: &ast.Value{String: &product.Currency},
			})
		}

		resource := &ast.Resource{
			ID:       product.ID,
			Typ:      product.ProductType,
			Requires: requires,
			Config:   config,
		}

		dslReq.Orchestrator.Resources = append(dslReq.Orchestrator.Resources, resource)
	}

	// Add explicit resources
	for _, resSpec := range resources {
		requires := []*ast.RequireItem{}
		for _, reqID := range resSpec.Requires {
			requires = append(requires, &ast.RequireItem{
				Kind: "entity",
				ID:   reqID,
			})
		}

		config := []*ast.KVPair{}
		for k, v := range resSpec.Config {
			strVal, ok := v.(string)
			if !ok {
				continue
			}
			config = append(config, &ast.KVPair{
				Key:   k,
				Value: &ast.Value{String: &strVal},
			})
		}

		resource := &ast.Resource{
			ID:       resSpec.ID,
			Typ:      resSpec.Type,
			Requires: requires,
			Config:   config,
		}

		dslReq.Orchestrator.Resources = append(dslReq.Orchestrator.Resources, resource)
	}
}

// generateFlows generates onboarding flows based on entities and products
func (g *Generator) generateFlows(dslReq *ast.Request) {
	steps := []*ast.Step{}

	// Step 1: Verify each entity
	for _, entity := range dslReq.Orchestrator.Entities {
		taskID := fmt.Sprintf("verify-%s", sanitizeID(entity.ID))

		// Determine verification type based on role
		var role string
		for _, attr := range entity.Attrs {
			if attr.Key == "role" {
				role = *attr.Value.Symbol
				break
			}
		}
		verificationLevel := "standard"
		if role == string(RoleSicav) || role == string(RoleManagementCompany) {
			verificationLevel = "enhanced"
		}

		step := &ast.Step{
			Task: &ast.Task{
				ID: taskID,
				On: "kyc-service",
				Op: "verify-entity",
				Args: []*ast.KVPair{
					{Key: "entity-id", Value: &ast.Value{String: &entity.ID}},
					{Key: "verification-level", Value: &ast.Value{String: &verificationLevel}},
				},
			},
		}
		steps = append(steps, step)
	}

	// Step 2: AML screening for all entities
	for _, entity := range dslReq.Orchestrator.Entities {
		taskID := fmt.Sprintf("aml-check-%s", sanitizeID(entity.ID))
		step := &ast.Step{
			Task: &ast.Task{
				ID: taskID,
				On: "aml-service",
				Op: "screen-entity",
				Args: []*ast.KVPair{
					{Key: "entity-id", Value: &ast.Value{String: &entity.ID}},
				},
			},
		}
		steps = append(steps, step)
	}

	// Step 3: Compliance review gate
	gateStep := &ast.Step{
		Gate: &ast.Gate{
			ID:        "compliance-review",
			Condition: "all-kyc-complete AND all-aml-clear",
		},
	}
	steps = append(steps, gateStep)

	// Step 4: Setup products/resources
	for _, resource := range dslReq.Orchestrator.Resources {
		taskID := fmt.Sprintf("setup-%s", sanitizeID(resource.ID))
		step := &ast.Step{
			Task: &ast.Task{
				ID: taskID,
				On: resource.ID,
				Op: g.getSetupOperation(resource.Typ),
				Args: []*ast.KVPair{
					{Key: "resource-id", Value: &ast.Value{String: &resource.ID}},
				},
			},
		}
		steps = append(steps, step)
	}

	// Create main flow
	mainFlow := &ast.Flow{
		ID:    "main",
		Steps: steps,
	}

	dslReq.Orchestrator.Flows = append(dslReq.Orchestrator.Flows, mainFlow)
}

// getSetupOperation returns the appropriate setup operation for a resource type
func (g *Generator) getSetupOperation(resourceType string) string {
	switch resourceType {
	case "CustodySafekeeping", "custody":
		return "create-account"
	case "investment-management":
		return "setup-mandate"
	case "reporting":
		return "configure-reporting"
	default:
		return "initialize"
	}
}

// sanitizeID removes problematic characters from IDs for use in task names
func sanitizeID(id string) string {
	// Simple sanitization: replace : with -
	result := ""
	for _, ch := range id {
		if ch == ':' {
			result += "-"
		} else if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' {
			result += string(ch)
		}
	}
	return result
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}
