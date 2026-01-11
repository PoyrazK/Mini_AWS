// Package ports defines the interfaces (ports) for the hexagonal architecture.
//
// This package contains all the interfaces that define the contracts between
// the core business logic and external adapters (databases, compute backends,
// storage systems, etc.).
//
// # Architecture Pattern
//
// TheCloud uses the Ports and Adapters (Hexagonal) architecture:
//
//	┌─────────────────────────────────────┐
//	│         Core Business Logic         │
//	│   (services, domain models)         │
//	│                                     │
//	│  Depends on → Ports (this package)  │
//	└─────────────────────────────────────┘
//	              ▲         ▲
//	              │         │
//	     ┌────────┘         └────────┐
//	     │                           │
//	┌────────────┐            ┌─────────────┐
//	│  Adapters  │            │  Adapters   │
//	│ (Inbound)  │            │ (Outbound)  │
//	│            │            │             │
//	│ HTTP API   │            │ PostgreSQL  │
//	│ gRPC       │            │ Docker      │
//	│ CLI        │            │ Redis       │
//	└────────────┘            └─────────────┘
//
// # Port Categories
//
// ## Service Ports (Primary/Inbound)
//
// Interfaces that define use cases exposed to external actors:
//
//   - InstanceService: Manage compute instances
//   - DatabaseService: Manage database instances
//   - VpcService: Manage virtual private clouds
//   - (etc.)
//
// ## Repository Ports (Secondary/Outbound)
//
// Interfaces for data persistence:
//
//   - InstanceRepository: Instance data access
//   - UserRepository: User data access
//   - (etc.)
//
// ## Infrastructure Ports (Secondary/Outbound)
//
// Interfaces for external systems:
//
//   - ComputeBackend: Docker, Libvirt, or Noop
//   - NetworkBackend: OVS, Docker networking
//   - StorageBackend: LVM, file-based storage
//   - Cache: Redis or in-memory
//   - Queue: Message queue abstraction
//
// # Dependency Rule
//
// Dependencies always point inward:
//
//	External → Adapters → Ports → Core
//
// The core business logic never depends on external frameworks or infrastructure.
//
// # Testing Benefits
//
// This architecture enables:
//
//   - Easy mocking of dependencies
//   - Testing business logic in isolation
//   - Swapping implementations (e.g., Docker ↔ Libvirt)
//   - Running without external dependencies (noop backends)
//
// # Example Usage
//
// Dependency injection in main.go:
//
//	// Create adapters
//	instanceRepo := postgres.NewInstanceRepository(db)
//	computeBackend := docker.NewDockerAdapter()
//
//	// Inject into service
//	instanceService := services.NewInstanceService(
//	    instanceRepo,     // InstanceRepository port
//	    computeBackend,   // ComputeBackend port
//	    // ...
//	)
//
// The service only knows about the ports (interfaces), not the concrete implementations.
package ports
