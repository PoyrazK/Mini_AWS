'use client';

import { useState } from 'react';
import { Table, Column } from '@/components/ui/Table';
import { StatusIndicator } from '@/components/ui/StatusIndicator';
import { Button } from '@/components/ui/Button';
import { LaunchInstanceModal } from '@/components/compute/LaunchInstanceModal';
import { Plus, RefreshCw } from 'lucide-react';

interface Instance {
  id: string;
  name: string;
  type: string;
  status: 'running' | 'stopped' | 'pending' | 'error';
  ip: string;
  created_at: string;
}

const DUMMY_INSTANCES: Instance[] = [
  { id: 'i-0x8231', name: 'Web Server 01', type: 't2.micro', status: 'running', ip: '10.0.1.12', created_at: '2025-01-10' },
  { id: 'i-0x992a', name: 'Worker Node', type: 't3.medium', status: 'running', ip: '10.0.1.15', created_at: '2025-01-11' },
  { id: 'i-0x11b2', name: 'DB Replica', type: 'm5.large', status: 'stopped', ip: '10.0.2.4', created_at: '2025-01-12' },
  { id: 'i-0x33c4', name: 'Cache Layer', type: 't2.small', status: 'error', ip: '-', created_at: '2025-01-14' },
];

export default function ComputePage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [instances, setInstances] = useState<Instance[]>(DUMMY_INSTANCES);

  const handleLaunch = (data: { name: string }) => {
    // Simulate instance launch
    const newInstance: Instance = {
      id: `i-0x${Math.floor(Math.random() * 10000).toString(16)}`,
      name: data.name,
      type: 't2.micro',
      status: 'pending',
      ip: '-',
      created_at: new Date().toISOString().split('T')[0]
    };
    setInstances([newInstance, ...instances]);
  };

  const columns: Column<Instance>[] = [
    { header: 'Name', accessorKey: 'name', width: '25%' },
    { header: 'Instance ID', accessorKey: 'id', width: '20%' },
    { header: 'Type', accessorKey: 'type', width: '15%' },
    { 
      header: 'Status', 
      cell: (item) => <StatusIndicator status={item.status} label={item.status} /> 
    },
    { header: 'Private IP', accessorKey: 'ip' },
    { header: 'Created', accessorKey: 'created_at' },
  ];

  return (
    <div style={{ maxWidth: '1280px', margin: '0 auto' }}>
      <header style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        marginBottom: '32px' 
      }}>
        <div>
           <h1 style={{ fontSize: '34px', fontWeight: 700, marginBottom: '4px', letterSpacing: '0.01em', color: 'var(--text-primary)' }}>Compute</h1>
           <p style={{ color: 'var(--text-secondary)' }}>Manage your virtual machines.</p>
        </div>
        <div style={{ display: 'flex', gap: '12px' }}>
          <Button variant="secondary"><RefreshCw size={16} /></Button>
          <Button onClick={() => setIsModalOpen(true)}>
            <Plus size={16} style={{ marginRight: '8px' }} /> Launch Instance
          </Button>
        </div>
      </header>

      <Table data={instances} columns={columns} />

      <LaunchInstanceModal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleLaunch}
      />
    </div>
  );
}
