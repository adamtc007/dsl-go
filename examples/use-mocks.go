package main

import (
	"fmt"
	"log"

	"github.com/example/dsl-go/internal/generator"
	"github.com/example/dsl-go/internal/manager"
	"github.com/example/dsl-go/internal/mocks"
)

// Example demonstrating how to use the mock data loader
// to generate onboarding DSL instances without a database

func main() {
	fmt.Println("=== Mock Data Loader Examples ===\n")

	// Create a mock data loader
	loader := mocks.NewDefaultLoader()

	// Example 1: List available mock data
	fmt.Println("1. Listing available mock data:")
	listAvailableMocks(loader)

	// Example 2: Load and inspect a single entity
	fmt.Println("\n2. Loading a single entity:")
	loadSingleEntity(loader)

	// Example 3: Load a complete scenario
	fmt.Println("\n3. Loading a complete scenario:")
	loadCompleteScenario(loader)

	// Example 4: Build a custom scenario
	fmt.Println("\n4. Building a custom scenario:")
	buildCustomScenario(loader)

	// Example 5: Generate DSL from scenario
	fmt.Println("\n5. Generating DSL from scenario:")
	generateDSLFromScenario(loader)

	fmt.Println("\n=== Examples Complete ===")
}

func listAvailableMocks(loader *mocks.Loader) {
	entities, err := loader.ListEntities()
	if err != nil {
		log.Printf("Error listing entities: %v", err)
		return
	}
	fmt.Printf("  Available entities: %v\n", entities)

	products, err := loader.ListProducts()
	if err != nil {
		log.Printf("Error listing products: %v", err)
		return
	}
	fmt.Printf("  Available products: %v\n", products)

	scenarios, err := loader.ListScenarios()
	if err != nil {
		log.Printf("Error listing scenarios: %v", err)
		return
	}
	fmt.Printf("  Available scenarios: %v\n", scenarios)
}

func loadSingleEntity(loader *mocks.Loader) {
	entity, err := loader.LoadEntity("investment-manager-001.json")
	if err != nil {
		log.Printf("Error loading entity: %v", err)
		return
	}

	fmt.Printf("  Loaded entity:\n")
	fmt.Printf("    ID: %s\n", entity.ID)
	fmt.Printf("    Name: %s\n", entity.Name)
	fmt.Printf("    Role: %s\n", entity.Role)
	fmt.Printf("    Country: %s\n", entity.Country)
	fmt.Printf("    LEI: %s\n", entity.LEI)

	if aum, ok := entity.Attributes["aum"]; ok {
		fmt.Printf("    AUM: %v\n", aum)
	}
}

func loadCompleteScenario(loader *mocks.Loader) {
	scenario, err := loader.LoadScenario("institutional-onboarding-001.json")
	if err != nil {
		log.Printf("Error loading scenario: %v", err)
		return
	}

	fmt.Printf("  Loaded scenario:\n")
	fmt.Printf("    Request ID: %s\n", scenario.RequestID)
	fmt.Printf("    Tenant ID: %s\n", scenario.TenantID)
	fmt.Printf("    Entities: %d\n", len(scenario.Entities))
	fmt.Printf("    Products: %d\n", len(scenario.Products))
	fmt.Printf("    Resources: %d\n", len(scenario.Resources))

	fmt.Printf("\n  Entity roles:\n")
	for _, entity := range scenario.Entities {
		fmt.Printf("    - %s (%s)\n", entity.Name, entity.Role)
	}

	fmt.Printf("\n  Products:\n")
	for _, product := range scenario.Products {
		fmt.Printf("    - %s (%s)\n", product.ID, product.ProductType)
	}
}

func buildCustomScenario(loader *mocks.Loader) {
	// Build a custom scenario by selecting specific entities and products
	customScenario, err := loader.BuildCustomScenario(
		"custom-onboard-example-001",
		[]string{
			"investment-manager-001.json",
			"sicav-001.json",
		},
		[]string{
			"custody-safekeeping-eur.json",
		},
	)

	if err != nil {
		log.Printf("Error building custom scenario: %v", err)
		return
	}

	fmt.Printf("  Custom scenario created:\n")
	fmt.Printf("    Request ID: %s\n", customScenario.RequestID)
	fmt.Printf("    Selected entities: %d\n", len(customScenario.Entities))
	fmt.Printf("    Selected products: %d\n", len(customScenario.Products))

	for _, entity := range customScenario.Entities {
		fmt.Printf("    - Entity: %s\n", entity.Name)
	}
	for _, product := range customScenario.Products {
		fmt.Printf("    - Product: %s\n", product.ID)
	}
}

func generateDSLFromScenario(loader *mocks.Loader) {
	// Load a scenario
	scenario, err := loader.LoadScenario("institutional-onboarding-001.json")
	if err != nil {
		log.Printf("Error loading scenario: %v", err)
		return
	}

	// Generate DSL
	gen, err := generator.New()
	if err != nil {
		log.Printf("Error creating generator: %v", err)
		return
	}
	response, err := gen.Generate(scenario)
	if err != nil {
		log.Printf("Error generating DSL: %v", err)
		return
	}

	fmt.Printf("  DSL generation successful:\n")
	fmt.Printf("    Request ID: %s\n", response.RequestID)
	fmt.Printf("    Version: %d\n", response.Version)
	fmt.Printf("    Entities added: %d\n", response.EntitiesAdded)
	fmt.Printf("    Resources added: %d\n", response.ResourcesAdded)
	fmt.Printf("    Flows generated: %d\n", response.FlowsGenerated)
	fmt.Printf("    Generated at: %s\n", response.GeneratedAt)
	fmt.Printf("    DSL length: %d bytes\n", len(response.DSL))

	// Optionally save to storage
	mgr, err := manager.New(manager.Config{
		DataDir:     "./data",
		RegistryDir: "./registry",
	})
	if err != nil {
		log.Printf("Error creating manager: %v", err)
		return
	}

	version, hash, err := mgr.CreateRequest(response.RequestID, response.DSL)
	if err != nil {
		log.Printf("Error saving DSL: %v", err)
		return
	}

	fmt.Printf("\n  Saved to storage:\n")
	fmt.Printf("    Version: %d\n", version)
	fmt.Printf("    Hash: %s\n", hash)
	fmt.Printf("    Path: ./data/%s\n", response.RequestID)

	fmt.Printf("\n  First 500 chars of generated DSL:\n")
	if len(response.DSL) > 500 {
		fmt.Printf("    %s...\n", response.DSL[:500])
	} else {
		fmt.Printf("    %s\n", response.DSL)
	}
}

// Example: Load entities by role
func exampleLoadByRole(loader *mocks.Loader) {
	entities, err := loader.LoadEntitiesByRole(generator.RoleInvestmentManager)
	if err != nil {
		log.Printf("Error loading entities by role: %v", err)
		return
	}

	fmt.Printf("Found %d investment managers:\n", len(entities))
	for _, entity := range entities {
		fmt.Printf("  - %s\n", entity.Name)
	}
}

// Example: Save a new entity mock
func exampleSaveEntity(loader *mocks.Loader) {
	newEntity := &generator.ClientEntity{
		ID:         "le:new-example-001",
		Name:       "Example Asset Management",
		Role:       generator.RoleInvestmentManager,
		EntityType: "LegalEntity",
		Country:    "CH",
		Attributes: map[string]interface{}{
			"aum":       "5000000000",
			"regulated": true,
			"regulator": "FINMA",
		},
	}

	err := loader.SaveEntity(newEntity, "example-entity.json")
	if err != nil {
		log.Printf("Error saving entity: %v", err)
		return
	}

	fmt.Println("Entity saved successfully!")
}
