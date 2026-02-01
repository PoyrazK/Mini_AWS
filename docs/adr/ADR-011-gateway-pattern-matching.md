# ADR-011: CloudGateway Advanced Path Matching

## Status
Accepted

## Context
CloudGateway initially supported only simple prefix-based routing (e.g., `/api` matches `/api/*`). As the platform grows, we need more sophisticated routing capabilities to support:
- RESTful API patterns (e.g., `/users/{id}`)
- Multi-parameter extraction (e.g., `/orgs/{oid}/projects/{pid}`)
- Regex-constrained parameters (e.g., `/id/{id:[0-9]+}`)
- File extension patterns (e.g., `/assets/*.{ext}`)

The core requirement is to match these patterns dynamically and extract the variable parts to make them available to downstream services or for internal business logic.

## Decison
We will implement an **Advanced Path Matching Engine** using a custom regex-based compiler.

1. **Pattern Syntax**:
   - `{name}`: Matches any non-slash character (`[^/]+`).
   - `{name:regex}`: Matches based on a custom regular expression.
   - `*`: Greedy wildcard matching (`.*`).
   - `**`: Explicit recursive wildcard matching (mapped to `.*`).

2. **HTTP Method Matching**:
   - Routes can be restricted to specific HTTP methods (e.g., `["GET", "POST"]`).
   - If `methods` is empty or null, the route matches all methods.
   - This allows multiple services to share the same path but handle different verbs.

3. **Pre-compiled Matchers**:
   - Routes are compiled into `*regexp.Regexp` instances and cached in-memory.
   - We avoid per-request regex compilation to maintain sub-millisecond routing latency.

4. **Routing Priority (Specificity Scoring)**:
   - When multiple routes match, the one with the highest **Specificity Score** wins.
   - Score = `len(LiteralPrefix)` + `exactMatchBonus(100)` + `Priority * 1000`.
   - `LiteralPrefix` is the static part of the pattern before any variable `{}` or `*`.
   - `exactMatchBonus` is added if the pattern contains no variables/wildcards.
   - `Priority` is an explicit user-defined integer (default 0).

5. **Hexagonal Integration**:
   - **Domain**: `GatewayRoute` stores `Methods []string` and `Priority int`.
   - **Database**: `gateway_routes` uses a composite unique constraint on `(path_pattern, methods)` instead of a global unique `path_prefix`.
   - **Service**: Returns `(proxy, params, found)` to propagate extracted values.
   - **Handler**: Extracted parameters are injected into the Gin context (e.g., `c.Set("path_param_id", value)`).

## Architecture

```
User Request (Method, Path) ──▶ Gateway Handler ──▶ Gateway Service (GetProxy)
                                                          │
                                                          ▼
                                                ┌───────────────────┐
                                                │ Matching Engine   │
                                                ├───────────────────┤
                                                │ 1. Method Filter  │
                                                │ 2. Pattern Match  │
                                                │ 3. Score Selection│
                                                └─────────┬─────────┘
                                                          │
                                                          ▼
                                                (ReverseProxy, Params)
```

## Consequences

### Positive
- **Flexibility**: Supports standard RESTful patterns and method-level routing.
- **Performance**: Pre-compilation and scoring ensures low latency (~0.05ms per match).
- **Backward Compatibility**: Existing prefix routes migrate seamlessly to the new system.

### Negative
- **Regex Backtracking**: Potential for complex regexes to slow down matching (mitigated by pre-compilation).
- **Ambiguity**: Overlapping patterns require careful use of `priority`.

## Implementation Notes
- **Constraint Change**: The unique index on `path_prefix` was dropped to allow same-path-different-method routing.
- **Refactoring**: The `GatewayService` was refactored to separate proxy creation, matching, and selection logic, reducing cognitive complexity and improving testability.
- **Testing**: E2E tests verify pattern matching, priority tie-breaking, and method filtering.
