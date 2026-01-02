# ADR-008: Auto-Scaling Strategy

## Status
Proposed

## Context
We need to implement auto-scaling that responds to CPU metrics for the The Cloud simulator. 
We validated two main approaches:
1. **Event-driven**: React to metric events in real-time as they are ingested.
2. **Polling-based**: Periodically evaluate stored metrics and make scaling decisions.

## Decision
**Use polling-based evaluation** with a 10-second ticker interval.

## Consequences

### Positive
*   **Simplicity**: Follows the existing `LBWorker` pattern, reducing cognitive load and architectural complexity.
*   **Decoupling**: No need for a complex message broker or event stream processing for metrics.
*   **Testability**: Deterministic behavior is easier to unit test and verify in integration tests.

### Negative
*   **Latency**: Introduces up to 10 seconds of delay in responding to metric spikes.
*   **Overhead**: Polling creates a constant base load on the database (mitigated by efficient batch queries).

### Trade-off Justification
For a local cloud simulator, a 10-second response time is acceptable and mirrors real-world "warm-up" times. AWS Auto Scaling Groups also have reaction latencies and cooldown periods (often 60-300 seconds), making sub-10s precision unnecessary for our use case.
