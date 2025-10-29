package mocks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/example/dsl-go/internal/generator"
)

// Loader provides access to mock data from JSON files
type Loader struct {
	basePath string
}

// NewLoader creates a new mock data loader with the specified base path
func NewLoader(basePath string) *Loader {
	return &Loader{
		basePath: basePath,
	}
}

// NewDefaultLoader creates a loader using the default data-mocks directory
func NewDefaultLoader() *Loader {
	return &Loader{
		basePath: "",
	}
}

// LoadEntity loads a single entity from a JSON file
func (l *Loader) LoadEntity(filename string) (*generator.ClientEntity, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read entity file %s: %w", filename, err)
	}

	var entity generator.ClientEntity
	if err := json.Unmarshal(data, &entity); err != nil {
		return nil, fmt.Errorf("failed to parse entity JSON from %s: %w", filename, err)
	}

	return &entity, nil
}

// LoadProduct loads a single product from a JSON file
func (l *Loader) LoadProduct(filename string) (*generator.ProductSpec, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read product file %s: %w", filename, err)
	}

	var product generator.ProductSpec
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, fmt.Errorf("failed to parse product JSON from %s: %w", filename, err)
	}

	return &product, nil
}

// LoadScenario loads a complete scenario from a JSON file
func (l *Loader) LoadScenario(filename string) (*generator.GenerateRequest, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read scenario file %s: %w", filename, err)
	}

	var scenario generator.GenerateRequest
	if err := json.Unmarshal(data, &scenario); err != nil {
		return nil, fmt.Errorf("failed to parse scenario JSON from %s: %w", filename, err)
	}

	return &scenario, nil
}

// LoadAllEntities loads all entity JSON files from the entities directory
func (l *Loader) LoadAllEntities() ([]generator.ClientEntity, error) {
	entitiesPath := filepath.Join(l.basePath, "entities")
	files, err := os.ReadDir(entitiesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read entities directory: %w", err)
	}

	entities := make([]generator.ClientEntity, 0, len(files))
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		entity, err := l.LoadEntity(file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to load entity %s: %w", file.Name(), err)
		}
		entities = append(entities, *entity)
	}

	return entities, nil
}

// LoadAllProducts loads all product JSON files from the products directory
func (l *Loader) LoadAllProducts() ([]generator.ProductSpec, error) {
	productsPath := filepath.Join(l.basePath, "products")
	files, err := os.ReadDir(productsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read products directory: %w", err)
	}

	products := make([]generator.ProductSpec, 0, len(files))
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		product, err := l.LoadProduct(file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to load product %s: %w", file.Name(), err)
		}
		products = append(products, *product)
	}

	return products, nil
}

// ListEntities returns a list of available entity mock files
func (l *Loader) ListEntities() ([]string, error) {
	entitiesPath := filepath.Join(l.basePath, "entities")
	files, err := os.ReadDir(entitiesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read entities directory: %w", err)
	}

	var names []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			names = append(names, file.Name())
		}
	}

	return names, nil
}

// ListProducts returns a list of available product mock files
func (l *Loader) ListProducts() ([]string, error) {
	productsPath := filepath.Join(l.basePath, "products")
	files, err := os.ReadDir(productsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read products directory: %w", err)
	}

	var names []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			names = append(names, file.Name())
		}
	}

	return names, nil
}

// ListScenarios returns a list of available scenario mock files
func (l *Loader) ListScenarios() ([]string, error) {
	scenariosPath := filepath.Join(l.basePath, "scenarios")
	files, err := os.ReadDir(scenariosPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read scenarios directory: %w", err)
	}

	var names []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			names = append(names, file.Name())
		}
	}

	return names, nil
}

// LoadEntitiesByRole loads all entities with a specific role
func (l *Loader) LoadEntitiesByRole(role generator.ClientRole) ([]generator.ClientEntity, error) {
	allEntities, err := l.LoadAllEntities()
	if err != nil {
		return nil, err
	}

	var filtered []generator.ClientEntity
	for _, entity := range allEntities {
		if entity.Role == role {
			filtered = append(filtered, entity)
		}
	}

	return filtered, nil
}

// BuildCustomScenario builds a custom scenario by selecting specific entities and products
func (l *Loader) BuildCustomScenario(requestID string, entityFiles []string, productFiles []string) (*generator.GenerateRequest, error) {
	entities := make([]generator.ClientEntity, 0, len(entityFiles))
	for _, filename := range entityFiles {
		entity, err := l.LoadEntity(filename)
		if err != nil {
			return nil, err
		}
		entities = append(entities, *entity)
	}

	products := make([]generator.ProductSpec, 0, len(productFiles))
	for _, filename := range productFiles {
		product, err := l.LoadProduct(filename)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}

	return &generator.GenerateRequest{
		RequestID: requestID,
		TenantID:  "default",
		Entities:  entities,
		Products:  products,
		Metadata:  make(map[string]interface{}),
	}, nil
}

// SaveEntity saves an entity to a JSON file
func (l *Loader) SaveEntity(entity *generator.ClientEntity, filename string) error {
	path := filepath.Join(l.basePath, "entities", filename)
	data, err := json.MarshalIndent(entity, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal entity: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write entity file: %w", err)
	}

	return nil
}

// SaveProduct saves a product to a JSON file
func (l *Loader) SaveProduct(product *generator.ProductSpec, filename string) error {
	path := filepath.Join(l.basePath, "products", filename)
	data, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write product file: %w", err)
	}

	return nil
}

// SaveScenario saves a scenario to a JSON file
func (l *Loader) SaveScenario(scenario *generator.GenerateRequest, filename string) error {
	path := filepath.Join(l.basePath, "scenarios", filename)
	data, err := json.MarshalIndent(scenario, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal scenario: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write scenario file: %w", err)
	}

	return nil
}
