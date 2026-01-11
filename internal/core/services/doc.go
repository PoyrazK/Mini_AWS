// Package services implements the business logic layer of TheCloud platform.
//
// This package contains all service implementations that coordinate between
// domain models and infrastructure adapters through the ports interfaces.
//
// # Service Layer Responsibilities
//
// Services in this package handle:
//
//  1. **Business Logic**: Orchestrating domain models and enforcing business rules
//  2. **Transaction Management**: Coordinating multi-step operations
//  3. **Authorization**: Ensuring users can only access their own resources
//  4. **Event Publishing**: Recording audit logs and system events
//  5. **Error Handling**: Converting infrastructure errors to domain errors
//  6. **Validation**: Input validation and business rule enforcement
//
// # Service Categories
//
// ## Compute Services
//
//   - InstanceService: VM/container lifecycle management
//   - DatabaseService: Managed database provisioning
//   - FunctionService: Serverless function execution
//   - ContainerService: Container management
//
// ## Networking Services
//
//   - VpcService: Virtual private cloud management
//   - SubnetService: Subnet and IP management
//   - SecurityGroupService: Firewall rule management
//   - LBService: Load balancer management
//   - GatewayService: API gateway management
//
// ## Storage Services
//
//   - VolumeService: Block storage management
//   - SnapshotService: Volume snapshot management
//   - StorageService: Object storage (S3-like)
//
// ## Platform Services
//
//   - QueueService: Message queue management
//   - CacheService: Cache instance management
//   - NotifyService: Pub/sub notifications
//   - CronService: Scheduled job management
//   - SecretService: Secrets management
//
// ## Infrastructure Services
//
//   - StackService: Infrastructure-as-Code deployment
//   - AutoScalingService: Auto-scaling group management
//   - RBACService: Role-based access control
//   - AuditService: Audit logging
//   - EventService: Event tracking
//
// # Service Pattern
//
// All services follow a consistent pattern:
//
//	type ServiceImpl struct {
//	    repo         ports.Repository      // Data access
//	    backend      ports.Infrastructure  // External system
//	    eventSvc     ports.EventService    // Event publishing
//	    auditSvc     ports.AuditService    // Audit logging
//	}
//
//	func NewService(dependencies...) ports.Service {
//	    return &ServiceImpl{ ... }
//	}
//
// # Context Usage
//
// Services use context.Context for:
//
//   - Request cancellation
//   - User ID propagation (via appcontext package)
//   - Timeout enforcement
//   - Request tracing
//
// Example:
//
//	userID := appcontext.UserIDFromContext(ctx)
//	// User can only access their own resources
//
// # Error Handling
//
// Services return domain errors from the errors package:
//
//	if !found {
//	    return errors.New(errors.NotFound, "instance not found")
//	}
//	if !authorized {
//	    return errors.New(errors.Forbidden, "not authorized")
//	}
//
// # Testing
//
// Services are tested using:
//
//   - Mock implementations of ports (in shared_test.go)
//   - Table-driven tests for business logic
//   - Integration tests with real databases
//
// See *_test.go files for examples.
//
// # Worker Services
//
// Background workers handle async operations:
//
//   - AutoScalingWorker: Evaluates scaling policies
//   - LBWorker: Health checks for load balancers
//   - CronWorker: Executes scheduled jobs
//   - ContainerWorker: Container state reconciliation
//
// Workers run in goroutines and use context for graceful shutdown:
//
//	func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
//	    defer wg.Done()
//	    ticker := time.NewTicker(interval)
//	    defer ticker.Stop()
//
//	    for {
//	        select {
//	        case <-ctx.Done():
//	            return
//	        case <-ticker.C:
//	            w.DoWork(ctx)
//	        }
//	    }
//	}
package services
