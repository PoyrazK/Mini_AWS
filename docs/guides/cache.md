# CloudCache Feature

CloudCache provides managed Redis instances for your applications. It abstracts away the complexity of managing Redis containers, networking, and configuration.

## Features

- **Managed Redis**: Automatically provision and manage Redis instances.
- **Version Support**: Supports Redis 7.2 (default) and other versions via provided image tags.
- **Custom Configuration**: Configure memory limits and password authentication details are handled automatically.
- **VPC Integration**: Deploy caches into specific VPCs for isolation (caches are currently accessible via host networking in this simulator version, but VPC IDs are tracked).
- **Persistence**: AOF (Append Only File) persistence is enabled by default.
- **Monitoring**: Basic stats (memory usage) available.

## Architecture

CloudCache follows the `thecloud` 3-tier architecture:
- **API**: Handles requests validation, DB state management, and orchestration.
- **Service**: `CacheService` manages the business logic and orchestrates Docker containers.
- **Repository**: `CacheRepository` (Postgres) stores metadata.
- **Infrastructure**: Uses Docker to spawn `redis:alpine` containers.

### Docker Configuration
Containers are launched with:
- `redis-server --appendonly yes --requirepass <generated_password> --maxmemory <limit>mb`
- Exposed on a random host port mapped to 6379.

## Usage

### CLI

1. **Create a Cache**
   ```bash
   cloud cache create --name my-cache --memory 256 --wait
   ```

2. **List Caches**
   ```bash
   cloud cache list
   ```

3. **Get Connection String**
   ```bash
   cloud cache connection my-cache
   ```
   Output: `redis://:password@localhost:32768`

4. **View Details**
   ```bash
   cloud cache show my-cache
   ```

5. **Delete Cache**
   ```bash
   cloud cache rm my-cache
   ```

### SDK

```go
import "github.com/poyrazk/thecloud/pkg/sdk"

client := sdk.NewClient(apiKey)

// Create
cache, err := client.CreateCache("app-cache", "7.2", 512, nil)

// Connect
connStr, err := client.GetCacheConnectionString(cache.ID)

// Use with go-redis
rdb := redis.NewClient(&redis.Options{Addr: connStr})
```

## Limitations (v1)

- **TLS**: Not enabled by default. Traffic is unencrypted.
- **Clustering**: Single node only.
- **Public Access**: Exposed on localhost/host-ip. Use Security Groups (future) or VPC peering for restrictions.
- **Flush**: `flush` command is currently a placeholder (requires Exec implementation).

## Troubleshooting

- **Memory 0MB**: Ensure you specify valid memory limits > 0.
- **Connection Refused**: Ensure the instance status is `RUNNING` and you are using the correct mapped port from `cloud cache show`.
