package generator

import "time"

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
	ID         string                 // Unique identifier (e.g., "le:ACME")
	Name       string                 // Legal name
	Role       ClientRole             // Role in the relationship
	EntityType string                 // LegalEntity, Individual, etc.
	LEI        string                 // Legal Entity Identifier (optional)
	Country    string                 // Jurisdiction/Country code
	Attributes map[string]interface{} // Additional attributes
}

// ProductSpec represents a product/service the client has contracted
type ProductSpec struct {
	ID          string                 // Product identifier (e.g., "prod:custody-eur")
	ProductType string                 // custody, investment-management, reporting, etc.
	Currency    string                 // Currency if applicable
	Config      map[string]interface{} // Product-specific configuration
}

// ResourceSpec represents a resource to be created during onboarding
type ResourceSpec struct {
	ID       string                 // Resource identifier
	Type     string                 // Resource type (CustodySafekeeping, Account, etc.)
	Requires []string               // IDs of entities/resources this depends on
	Config   map[string]interface{} // Resource configuration
}

// GenerateRequest contains all data needed to generate a populated DSL instance
type GenerateRequest struct {
	RequestID string            // Unique onboarding request ID
	TenantID  string            // Multi-tenant identifier
	Entities  []ClientEntity    // Client entities with their roles
	Products  []ProductSpec     // Products being onboarded
	Resources []ResourceSpec    // Resources to create
	Metadata  map[string]string // Additional metadata
}

// GenerateResponse contains the generated DSL and metadata
type GenerateResponse struct {
	RequestID      string    // The request ID
	DSL            string    // Generated S-expression DSL
	Version        uint64    // Version number (typically 1 for new)
	Hash           string    // Content hash
	GeneratedAt    time.Time // When it was generated
	EntitiesAdded  int       // Count of entities
	ResourcesAdded int       // Count of resources
	FlowsGenerated int       // Count of flows generated
}

// ValidationError represents an error during validation
type ValidationError struct {
	Field   string // Field that failed validation
	Message string // Error message
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
