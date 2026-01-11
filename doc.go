// Package thecloud provides a comprehensive cloud infrastructure platform built in Go.
//
// TheCloud is a multi-backend cloud infrastructure platform that supports multiple
// compute backends (Docker, Libvirt, Noop) and provides a wide range of cloud services
// including compute instances, networking (VPC, OVS), storage (volumes, snapshots),
// databases, caching, queuing, functions-as-a-service, and more.
//
// # Architecture
//
// The platform follows a hexagonal/ports-and-adapters architecture:
//
//   - cmd/: Command-line interfaces (API server, CLI client)
//   - internal/core/: Business logic and domain models
//   - internal/handlers/: HTTP handlers and API endpoints
//   - internal/repositories/: Infrastructure adapters (Docker, Libvirt, PostgreSQL, Redis)
//   - pkg/: Reusable packages (crypto, audit, ratelimit, httputil, SDK)
//
// # Key Features
//
//   - **Multi-Backend Compute**: Docker containers, Libvirt VMs, or No-op testing
//   - **Advanced Networking**: VPC, subnets, security groups, Open vSwitch integration
//   - **Storage Management**: Volumes, snapshots, LVM backend support
//   - **Managed Services**: PostgreSQL/MySQL databases, Redis caching, message queuing
//   - **Functions-as-a-Service**: Serverless function execution
//   - **Auto-scaling**: Policy-based horizontal scaling with load balancing
//   - **Infrastructure-as-Code**: Stack templates for declarative deployment
//   - **RBAC & Multi-tenancy**: Role-based access control with user isolation
//   - **Observability**: Prometheus metrics, audit logging, event tracking
//
// # Quick Start
//
// Start the API server:
//
//	go run cmd/api/main.go
//
// Use the CLI client:
//
//	cloud instance create -n my-instance -i ubuntu -p "8080:80"
//	cloud db create -n my-database -e postgres
//	cloud function deploy -n my-function -r go1.21
//
// # Configuration
//
// Configure via environment variables:
//
//   - DATABASE_URL: PostgreSQL connection string
//   - REDIS_ADDR: Redis server address
//   - COMPUTE_BACKEND: docker|libvirt|noop
//   - NETWORK_BACKEND: docker|ovs|noop
//   - JWT_SECRET: Secret for JWT token signing
//
// # Security
//
// TheCloud implements multiple security layers:
//
//   - JWT-based authentication
//   - API key support for programmatic access
//   - Role-based access control (RBAC)
//   - Resource isolation by user/tenant
//   - Rate limiting on API endpoints
//   - Audit logging of all operations
//   - Secrets encryption at rest
//
// # Development
//
// Run tests:
//
//	go test ./...
//	go test -tags=integration ./internal/repositories/postgres/...
//
// Build:
//
//	make build
//
// # API Documentation
//
// Swagger/OpenAPI documentation is available at /swagger/index.html when the
// API server is running.
//
// For more information, see the documentation in the docs/ directory.
package thecloud
