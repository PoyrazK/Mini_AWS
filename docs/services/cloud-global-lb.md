# Global Load Balancer (GLB)

## Overview
The **Global Load Balancer** (GLB) provides multi-region traffic distribution at the DNS level. Unlike regional load balancers that operate at the network layer within a VPC, GLB utilizes GeoDNS to steer users to the optimal regional endpoint based on policies like latency, health, and weight.

## Features
- **Global Traffic Steering**: Route traffic across multiple regions.
- **Health-Aware Routing**: Automatically removes unhealthy regional endpoints from DNS resolution.
- **Policy-Based Distribution**:
    - **Latency**: Directs users to the region with the lowest network latency.
    - **Geolocation**: Routes traffic based on the user's geographic location.
    - **Weighted**: Distributes traffic proportionally across regions.
    - **Failover**: Priority-based failover for disaster recovery.
- **Unified Hostname**: Provide a single global hostname for your application (e.g., `api.global.example.com`).

## Architecture
GLB follows the standard hexadecimal architecture of the platform:

1.  **Service Layer**: Handles the business logic of GLB creation, endpoint management, and synchronization with DNS.
2.  **GeoDNS Adapter**: Communicates with the authoritative DNS server (PowerDNS) to update records in real-time.
3.  **Repository**: Stores metadata in PostgreSQL about GLB configurations and endpoint health status.

## Configuration
### Health Checks
GLB performs synthesized health checks from multiple points of presence. Configuration includes:
- **Protocol**: HTTP, HTTPS, or TCP.
- **Port**: The destination port to probe.
- **Interval**: Frequency of probes (default 30s).
- **Thresholds**: Number of consecutive successes/failures to change health status.

### Endpoints
Endpoints can be:
- **Regional Load Balancers**: Linked by ID to existing platform resource.
- **External IPs**: Arbitrary static IPs for hybrid-cloud scenarios.

## API Usage
Create a GLB:
```bash
cloud global-lb create --name "prod-api" --hostname "api.global.com" --policy "LATENCY"
```

Add a regional endpoint:
```bash
cloud global-lb add-endpoint --id <glb-id> --region "us-east-1" --target-ip "1.2.3.4" --weight 100
```
