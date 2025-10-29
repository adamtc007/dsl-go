package generator

import (
	"time"

	"github.com/example/dsl-go/internal/manager"
)

// ClientRole represents the role of a client entity in the onboarding
type ClientRole string

const (
	RoleInvestmentManager ClientRole = "investment-manager"
	RoleAssetOwner        ClientRole = "asset-owner"
	RoleManagementCompany ClientRole = "management-company"
	RoleSicav             ClientRole = "sicav"
	RoleCustodian         ClientRole = "custodian"
	RoleAdministrator     ClientRole = "administrator"
)

// ClientEntity represents a legal entity being onboarded with their role
type ClientEntity struct {
	ID         string                 `json:"id"`          // Unique identifier (e.g., "le:ACME")
	Name       string                 `json:"name"`        // Legal name
	Role       ClientRole             `json:"role"`        // Role in the relationship
	EntityType string                 `json:"entity_type"` // LegalEntity, Individual, etc.
	LEI        string                 `json:"lei"`         // Legal Entity Identifier (optional)
	Country    string                 `json:"country"`     // Jurisdiction/Country code
	Attributes map[string]interface{} `json:"attributes"`  // Additional attributes
}

// ProductSpec represents a product/service the client has contracted
type ProductSpec struct {
	ID          string                 `json:"id"`           // Product identifier (e.g., "prod:custody-eur")
	ProductType string                 `json:"product_type"` // custody, investment-management, reporting, etc.
	Currency    string                 `json:"currency"`     // Currency if applicable
	Config      map[string]interface{} `json:"config"`       // Product-specific configuration
}

// ResourceSpec represents a resource to be created during onboarding
type ResourceSpec struct {
	ID       string                 `json:"id"`       // Resource identifier
	Type     string                 `json:"type"`     // Resource type (CustodySafekeeping, Account, etc.)
	Requires []string               `json:"requires"` // IDs of entities/resources this depends on
	Config   map[string]interface{} `json:"config"`   // Resource configuration
}

// GenerateRequest contains all data needed to generate a populated DSL instance
type GenerateRequest struct {
	RequestID      string                  `json:"request_id"` // Unique onboarding request ID
	TenantID       string                  `json:"tenant_id"`  // Multi-tenant identifier
	Entities       []ClientEntity          `json:"entities"`   // Client entities with their roles
	Products       []ProductSpec           `json:"products"`   // Products being onboarded
	Resources      []ResourceSpec          `json:"resources"`  // Resources to create
	Metadata       map[string]interface{}  `json:"metadata"`   // Additional metadata (supports nested objects)
	Now            time.Time               `json:"-"`          // The current time, for use in templates
	DataDictionary *manager.DataDictionary `json:"-"`          // The data dictionary
}

// ValidationError represents an error during validation
type ValidationError struct {
	Field   string // Field that failed validation
	Message string // Error message
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func (r *GenerateRequest) GetProduct(id string) *manager.Product {
	for _, p := range r.DataDictionary.Products {
		if p.ProductID == id {
			return &p
		}
	}
	return nil
}

func (r *GenerateRequest) GetService(id string) *manager.Service {
	for _, s := range r.DataDictionary.Services {
		if s.ServiceID == id {
			return &s
		}
	}
	return nil
}
