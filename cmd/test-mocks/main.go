package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/example/dsl-go/internal/generator"
	"github.com/example/dsl-go/internal/manager"
	"github.com/example/dsl-go/internal/mocks"
)

func main() {
	fmt.Println("=== Mock Data Loader Test ===")

	// Create mock data loader
	loader := mocks.NewDefaultLoader()

	// Test 1: List available mock data
	fmt.Println("1. Listing available mock data:")
	testListMocks(loader)

	// Test 2: Load a single entity
	fmt.Println("\n2. Loading single entity:")
	testLoadEntity(loader)

	// Test 3: Load a single product
	fmt.Println("\n3. Loading single product:")
	testLoadProduct(loader)

	// Test 4: Load complete scenario
	fmt.Println("\n4. Loading complete scenario:")
	testLoadScenario(loader)

	// Test 5: Build custom scenario
	fmt.Println("\n5. Building custom scenario:")
	testBuildCustom(loader)

	// Test 6: Load entities by role
	fmt.Println("\n6. Loading entities by role:")
	testLoadByRole(loader)

	// Test 7: Generate DSL from scenario
	fmt.Println("\n7. Generating DSL from scenario:")
	testGenerateDSL(loader)

	fmt.Println("\n=== All Tests Complete ===")
}

func testListMocks(loader *mocks.Loader) {
	entities, err := loader.ListEntities()
	if err != nil {
		log.Printf("  ❌ Error listing entities: %v", err)
		return
	}
	fmt.Printf("  ✅ Found %d entity files:\n", len(entities))
	for _, e := range entities {
		fmt.Printf("     - %s\n", e)
	}

	products, err := loader.ListProducts()
	if err != nil {
		log.Printf("  ❌ Error listing products: %v", err)
		return
	}
	fmt.Printf("  ✅ Found %d product files:\n", len(products))
	for _, p := range products {
		fmt.Printf("     - %s\n", p)
	}

	scenarios, err := loader.ListScenarios()
	if err != nil {
		log.Printf("  ❌ Error listing scenarios: %v", err)
		return
	}
	fmt.Printf("  ✅ Found %d scenario files:\n", len(scenarios))
	for _, s := range scenarios {
		fmt.Printf("     - %s\n", s)
	}
}

func testLoadEntity(loader *mocks.Loader) {
	entity, err := loader.LoadEntity("investment-manager-001.json")
	if err != nil {
		log.Printf("  ❌ Error: %v", err)
		return
	}

	fmt.Printf("  ✅ Loaded entity successfully:\n")
	fmt.Printf("     ID: %s\n", entity.ID)
	fmt.Printf("     Name: %s\n", entity.Name)
	fmt.Printf("     Role: %s\n", entity.Role)
	fmt.Printf("     Type: %s\n", entity.EntityType)
	fmt.Printf("     Country: %s\n", entity.Country)
	fmt.Printf("     LEI: %s\n", entity.LEI)
	fmt.Printf("     Attributes: %d fields\n", len(entity.Attributes))

	if aum, ok := entity.Attributes["aum"]; ok {
		fmt.Printf("     AUM: %v\n", aum)
	}
	if regulated, ok := entity.Attributes["regulated"]; ok {
		fmt.Printf("     Regulated: %v\n", regulated)
	}
}

func testLoadProduct(loader *mocks.Loader) {
	product, err := loader.LoadProduct("custody-safekeeping-eur.json")
	if err != nil {
		log.Printf("  ❌ Error: %v", err)
		return
	}

	fmt.Printf("  ✅ Loaded product successfully:\n")
	fmt.Printf("     ID: %s\n", product.ID)
	fmt.Printf("     Type: %s\n", product.ProductType)
	fmt.Printf("     Currency: %s\n", product.Currency)
	fmt.Printf("     Config fields: %d\n", len(product.Config))

	if accountType, ok := product.Config["account_type"]; ok {
		fmt.Printf("     Account Type: %v\n", accountType)
	}
}

func testLoadScenario(loader *mocks.Loader) {
	scenario, err := loader.LoadScenario("institutional-onboarding-001.json")
	if err != nil {
		log.Printf("  ❌ Error: %v", err)
		return
	}

	fmt.Printf("  ✅ Loaded scenario successfully:\n")
	fmt.Printf("     Request ID: %s\n", scenario.RequestID)
	fmt.Printf("     Tenant ID: %s\n", scenario.TenantID)
	fmt.Printf("     Entities: %d\n", len(scenario.Entities))
	fmt.Printf("     Products: %d\n", len(scenario.Products))
	fmt.Printf("     Resources: %d\n", len(scenario.Resources))

	fmt.Printf("\n     Entity details:\n")
	for i, entity := range scenario.Entities {
		fmt.Printf("       %d. %s (%s) - %s\n", i+1, entity.Name, entity.Role, entity.Country)
	}

	fmt.Printf("\n     Product details:\n")
	for i, product := range scenario.Products {
		fmt.Printf("       %d. %s (%s)\n", i+1, product.ID, product.ProductType)
	}
}

func testBuildCustom(loader *mocks.Loader) {
	customScenario, err := loader.BuildCustomScenario(
		"test-custom-onboard-001",
		[]string{
			"investment-manager-001.json",
			"sicav-001.json",
		},
		[]string{
			"custody-safekeeping-eur.json",
		},
	)

	if err != nil {
		log.Printf("  ❌ Error: %v", err)
		return
	}

	fmt.Printf("  ✅ Built custom scenario:\n")
	fmt.Printf("     Request ID: %s\n", customScenario.RequestID)
	fmt.Printf("     Entities: %d\n", len(customScenario.Entities))
	fmt.Printf("     Products: %d\n", len(customScenario.Products))

	for _, entity := range customScenario.Entities {
		fmt.Printf("       - %s\n", entity.Name)
	}
}

func testLoadByRole(loader *mocks.Loader) {
	sicavs, err := loader.LoadEntitiesByRole(generator.RoleSicav)
	if err != nil {
		log.Printf("  ❌ Error: %v", err)
		return
	}

	fmt.Printf("  ✅ Found %d SICAV entities:\n", len(sicavs))
	for _, sicav := range sicavs {
		fmt.Printf("     - %s (ID: %s)\n", sicav.Name, sicav.ID)
	}

	investmentMgrs, err := loader.LoadEntitiesByRole(generator.RoleInvestmentManager)
	if err != nil {
		log.Printf("  ❌ Error: %v", err)
		return
	}

	fmt.Printf("  ✅ Found %d Investment Manager entities:\n", len(investmentMgrs))
	for _, mgr := range investmentMgrs {
		fmt.Printf("     - %s (ID: %s)\n", mgr.Name, mgr.ID)
	}
}

func testGenerateDSL(loader *mocks.Loader) {
	// Load scenario
	scenario, err := loader.LoadScenario("institutional-onboarding-001.json")
	if err != nil {
		log.Printf("  ❌ Error loading scenario: %v", err)
		return
	}

	// Generate DSL
	gen, err := generator.New()
	if err != nil {
		log.Printf("  ❌ Error creating generator: %v", err)
		return
	}
	response, err := gen.Generate(scenario)
	if err != nil {
		log.Printf("  ❌ Error generating DSL: %v", err)
		return
	}

	fmt.Printf("  ✅ DSL generation successful:\n")
	fmt.Printf("     Request ID: %s\n", response.RequestID)
	fmt.Printf("     Version: %d\n", response.Version)
	fmt.Printf("     Entities added: %d\n", response.EntitiesAdded)
	fmt.Printf("     Resources added: %d\n", response.ResourcesAdded)
	fmt.Printf("     Flows generated: %d\n", response.FlowsGenerated)
	fmt.Printf("     Generated at: %s\n", response.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("     DSL size: %d bytes\n", len(response.DSL))

	// Save to storage
	mgr, err := manager.New(manager.Config{
		DataDir:     "./data",
		RegistryDir: "./registry",
	})
	if err != nil {
		log.Printf("  ❌ Error creating manager: %v", err)
		return
	}

	version, hash, err := mgr.CreateRequest(response.RequestID, response.DSL)
	if err != nil {
		log.Printf("  ⚠️  Warning: Could not save to storage: %v", err)
	} else {
		fmt.Printf("\n  ✅ Saved to storage:\n")
		fmt.Printf("     Version: %d\n", version)
		fmt.Printf("     Hash: %s\n", hash)
		fmt.Printf("     Path: ./data/%s\n", response.RequestID)
	}

	// Show preview of DSL
	fmt.Printf("\n  DSL Preview (first 800 chars):\n")
	fmt.Println("  " + strings.Repeat("-", 70))
	if len(response.DSL) > 800 {
		fmt.Printf("%s\n  ...\n", response.DSL[:800])
	} else {
		fmt.Printf("%s\n", response.DSL)
	}
	fmt.Println("  " + strings.Repeat("-", 70))

	// Validate the generated DSL
	fmt.Printf("\n  Validating generated DSL...\n")
	issues, err := mgr.ValidateText(response.DSL)
	if err != nil {
		log.Printf("  ❌ Validation error: %v", err)
		return
	}

	if len(issues) == 0 {
		fmt.Printf("  ✅ DSL is valid!\n")
	} else {
		fmt.Printf("  ⚠️  Validation issues found:\n")
		for _, issue := range issues {
			fmt.Printf("     - %s\n", issue)
		}
	}

	// Offer to write DSL to file
	fmt.Printf("\n  Writing DSL to file for inspection...\n")
	filename := fmt.Sprintf("generated-%s.sexpr", response.RequestID)
	if err := os.WriteFile(filename, []byte(response.DSL), 0644); err != nil {
		log.Printf("  ⚠️  Could not write file: %v", err)
	} else {
		fmt.Printf("  ✅ Written to: %s\n", filename)
		fmt.Printf("     View with: cat %s\n", filename)
	}
}
