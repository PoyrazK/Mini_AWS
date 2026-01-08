#!/bin/bash
set -e

# setup-ovs.sh - Automated OVS setup for TheCloud
# Requirements: Ubuntu/Debian

echo "--- Installing Open vSwitch ---"
sudo apt-get update
sudo apt-get install -y openvswitch-switch openvswitch-common

echo "--- Configuring Permissions ---"
# Allow current user to manage OVS without sudo
USER_NAME=$(whoami)
# Note: On some systems, adding to the 'ovs' or 'openvswitch' group works if it exists.
# Otherwise, we use ACLs or set permissions on the socket.
sudo setfacl -m u:$USER_NAME:rw /var/run/openvswitch/db.sock

echo "--- Verifying Service ---"
sudo systemctl enable openvswitch-switch
sudo systemctl start openvswitch-switch
ovs-vsctl show

echo "--- Done! OVS is ready for TheCloud ---"
