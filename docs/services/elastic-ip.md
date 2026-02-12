# Elastic IP Service

Elastic IP (EIP) is a managed service that provides static public IPv4 addresses for dynamic cloud computing. It allows you to mask instance failures by rapidly remaping your public IP address to another instance in your account.

## Service Architecture

### Components
1. **API Handler**: Manages REST requests for allocation, association, and release.
2. **ElasticIP Service**: Core business logic ensuring unique allocation and valid associations.
3. **PostgreSQL Repository**: Stores the state of IP addresses, including their allocation timeframe and association.
4. **Compute Backend**: The service integrates with the compute layer (Docker/Libvirt) to configure network interfaces on the target instances.

### Data Flow
1. **Allocation**: User requests a new IP. The system reserves a unique address from the pool (e.g., `100.64.0.0/10`).
2. **Association**: User requests to associate an EIP with an Instance.
   - The system verifies the instance exists and is running.
   - It updates the database record to link the IP to the Instance UUID.
   - It triggers a network reconfiguration on the instance (if supported by backend) to alias the IP.
3. **Release**: User returns the IP to the pool. It must be disassociated first.

## Usage Guide

### Allocating an IP
You can allocate an IP to your project to reserve it for future use.

```bash
# Allocate
curl -X POST http://localhost:8080/elastic-ips \
  -H "X-API-Key: $API_KEY"

# Response
{
  "id": "eip-uuid...",
  "public_ip": "100.64.12.15"
}
```

### Associating with an Instance
Once allocated, you can attach it to any instance in your VPC.

```bash
curl -X POST http://localhost:8080/elastic-ips/$EIP_ID/associate \
  -H "X-API-Key: $API_KEY" \
  -d '{"instance_id": "inst-uuid..."}'
```

### Disassociating
To move an IP or release it, disassociate it first.

```bash
curl -X POST http://localhost:8080/elastic-ips/$EIP_ID/disassociate \
  -H "X-API-Key: $API_KEY"
```

## Implementation Details

### IP Pool
For this implementation, we simulate a public IP pool using the Carrier-Grade NAT (CGNAT) range `100.64.0.0/10`.
- **First Octet**: 100
- **Second Octet**: 64 + (derived from UUID)
- **Deterministic Generation**: The IP is generated deterministically from the UUID to ensuring consistency without a complex IPAM State Machine for this demo.

### Infrastructure Constraints
- **One EIP per Instance**: Currently, an instance can only have one Elastic IP associated with it.
- **VPC Scoping**: EIPs are currently region/VPC agnostic in allocation but associated instances must have valid networking.
