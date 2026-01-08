---
description: Setup and verify Open vSwitch networking environment
---

This workflow automates the installation and configuration of Open vSwitch (OVS) required for Sprint 3 networking features.

// turbo-all
1. Make setup script executable
```bash
chmod +x /home/poyraz/dev/Cloud/scripts/setup-ovs.sh
```

2. Run the setup script
```bash
/home/poyraz/dev/Cloud/scripts/setup-ovs.sh
```

3. Update Docker Compose for OVS support
```bash
# This adds NET_ADMIN capabilities and mounts the OVS socket
```

4. Verify OVS integration via API health check
```bash
# Wait for API to restart and check /health/ovs
```
