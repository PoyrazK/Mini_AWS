# Future Plans & Contributing

This document outlines planned features and how you can contribute to The Cloud.

---

## 🎯 Active Development

### Now Accepting Contributions

| Feature | Difficulty | Good First Issue? | Description |
|---------|------------|-------------------|-------------|
| **Postgres Repo Tests** | Easy | ✅ Yes | Add tests to `internal/repositories/postgres/` |
| **SDK Tests** | Easy | ✅ Yes | Add tests to `pkg/sdk/` |
| **API Docs (OpenAPI)** | Medium | ✅ Yes | Generate Swagger spec from handlers |
| **Metrics Collection** | Medium | No | Populate `metrics_history` table |
| **RBAC** | Hard | No | Role-Based Access Control system |

### In Progress (Maintainers)

| Feature | Branch | Owner | ETA |
|---------|--------|-------|-----|
| Web Dashboard | `jack/main` | @jack | Q1 2026 |
| RDS (Database Service) | TBD | Unassigned | Q1 2026 |

---

## 📋 Feature Backlog

### High Priority
- [ ] **RBAC** - User roles (admin, developer, read-only)
- [ ] **RDS** - Managed PostgreSQL/MySQL containers
- [ ] **Secrets Manager** - Encrypted secret storage

### Medium Priority
- [ ] **CloudCache** - Managed Redis instances
- [ ] **CloudQueue** - SQS-like message queue
- [ ] **Snapshots** - Volume backup/restore

### Low Priority
- [ ] **CloudFormation Templates** - IaC YAML definitions
- [ ] **Multi-region** - Cluster support

---

## 🛠️ How to Contribute

### 1. Pick an Issue
Choose from "Good First Issue" items above or check [GitHub Issues](https://github.com/PoyrazK/thecloud/issues).

### 2. Fork & Clone
```bash
git clone https://github.com/YOUR_USERNAME/Mini_AWS.git
cd Mini_AWS
```

### 3. Create a Branch
```bash
git checkout -b feature/your-feature-name
```

### 4. Follow Project Structure
```
internal/
├── core/domain/    # Data structures
├── core/ports/     # Interfaces
├── core/services/  # Business logic
├── handlers/       # HTTP endpoints
└── repositories/   # Database/Docker adapters
```

### 5. Write Tests
- Place `_test.go` files next to the code
- Use `testify/mock` for mocking

### 6. Submit PR
- Reference any related issues
- Include test coverage
- Update docs if needed

---

## 📊 Current Test Coverage Goals

| Package | Current | Target |
|---------|---------|--------|
| `services/` | 19% | **60%** |
| `handlers/` | 12% | **50%** |
| `repositories/postgres/` | 66% | **40%** |
| `pkg/sdk/` | 51% | **50%** |

---

## 📞 Contact

- Open an issue for questions
- Tag maintainers for review

*Last updated: 2026-01-02*
