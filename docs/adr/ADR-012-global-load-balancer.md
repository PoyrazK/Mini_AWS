# ADR-012: Global Load Balancer Design

## Status
Accepted

## Context
As the platform scales to support multi-region deployments, users need a way to route traffic across different geographic locations to ensure high availability and low latency. Existing regional load balancers are confined to a single VPC/Region. We need a "Global" tier that operates at the DNS level to steer traffic effectively.

## Decision
We decided to implement a **Global Load Balancer (GLB)** service with the following architectural components:

1.  **GeoDNS Orchestration**:
    *   Instead of proxying traffic at the edge (which adds latency), we utilize **DNS-level steering**.
    *   The GLB service acts as an orchestrator for a **Geo-route capable DNS backend** (initially implemented via PowerDNS).
2.  **Hexagonal Architecture Integration**:
    *   **Domain**: Defined `GlobalLoadBalancer` and `GlobalEndpoint` entities.
    *   **Ports**: Created `GlobalLBRepository` for persistent state and `GeoDNSBackend` to abstract DNS providers.
    *   **Adapters**: Implemented a specialized PowerDNS adapter to manage resource records dynamically.
3.  **Routing Policies**:
    *   `LATENCY`: Route based on network proximity (implemented via GeoDNS regional logic).
    *   `GEOLOCATION`: Route based on specific country/continent rules.
    *   `WEIGHTED`: Distribute traffic across endpoints based on relative weights.
    *   `FAILOVER`: Priority-based steering with automatic health-tracking.
4.  **Endpoint Logic**:
    *   Endpoints can be **Regional Load Balancers** (integrated via `target_id`) or **Static IPs**.
    *   Health checks are centralized at the GLB level to ensure consistent steering decisions.

## Consequences
*   **Positive**: Reduced latency for global users by steering them to the nearest healthy region.
*   **Positive**: Increased resilience against regional outages.
*   **Neutral**: Increased complexity in DNS management and health-check synchronization.
*   **Negative**: Dependency on PowerDNS advanced features (LUA records) for granular steering accuracy in future phases.
