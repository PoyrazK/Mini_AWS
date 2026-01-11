// Package domain defines the core business entities and value objects for TheCloud platform.
//
// This package contains all domain models representing cloud resources, including:
//
//   - Compute: Instance, Database, Function, Container
//   - Networking: VPC, Subnet, SecurityGroup, LoadBalancer, Gateway
//   - Storage: Volume, Snapshot, Object
//   - Platform Services: Queue, Cache, Notification, CronJob
//   - Infrastructure-as-Code: Stack, StackResource
//   - Access Control: User, Role, Permission, APIKey
//   - Monitoring: Event, AuditLog, MetricPoint
//
// # Design Principles
//
// Domain models follow these principles:
//
//  1. Rich domain models with behavior, not anemic data structures
//  2. Immutable value objects where appropriate
//  3. Clear ownership and lifecycle management
//  4. Resource isolation per user/tenant
//  5. Status tracking for async operations
//
// # Resource Lifecycle
//
// Most resources follow this lifecycle:
//
//	Creating → Active → (Updating) → Deleting → Deleted
//
// With error states:
//
//	Creating → Failed
//	Active → Error
//
// # Multi-tenancy
//
// All resources include a UserID field for tenant isolation:
//
//	type Instance struct {
//	    ID     uuid.UUID
//	    UserID uuid.UUID // Tenant isolation
//	    // ...
//	}
//
// # Validation
//
// Domain entities include validation logic:
//
//   - ARN format validation
//   - CIDR block validation for networks
//   - Port range validation
//   - Resource name constraints
//
// See individual type documentation for specific validation rules.
package domain
