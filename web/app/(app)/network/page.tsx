'use client';

import { Table, Column } from '@/components/ui/Table';
import { StatusIndicator } from '@/components/ui/StatusIndicator';
import { Button } from '@/components/ui/Button';
import { Plus, RefreshCw, Network } from 'lucide-react';

interface VPC {
  id: string;
  name: string;
  cidr: string;
  status: 'available' | 'pending';
  subnets: number;
}

const DUMMY_VPCS: VPC[] = [
  { id: 'vpc-0x12a', name: 'default-vpc', cidr: '172.31.0.0/16', status: 'available', subnets: 4 },
  { id: 'vpc-0x44b', name: 'prod-network', cidr: '10.0.0.0/16', status: 'available', subnets: 6 },
  { id: 'vpc-0x99c', name: 'dev-environment', cidr: '192.168.0.0/16', status: 'pending', subnets: 0 },
];

export default function NetworkPage() {
  const columns: Column<VPC>[] = [
    { 
      header: 'Name', 
      cell: (item) => (
         <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontWeight: 500 }}>
          <Network size={16} color="var(--accent-orange)" />
          {item.name}
        </div>
      ) 
    },
    { header: 'VPC ID', accessorKey: 'id' },
    { header: 'IPv4 CIDR', accessorKey: 'cidr' },
    { 
      header: 'Status', 
      cell: (item) => <StatusIndicator status={item.status === 'available' ? 'running' : 'pending'} label={item.status} /> 
    },
    { header: 'Subnets', accessorKey: 'subnets' },
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
           <h1 style={{ fontSize: '34px', fontWeight: 700, marginBottom: '4px', letterSpacing: '0.01em', color: 'var(--text-primary)' }}>Network</h1>
           <p style={{ color: 'var(--text-secondary)' }}>Virtual Private Clouds and subnets.</p>
        </div>
        <div style={{ display: 'flex', gap: '12px' }}>
          <Button variant="secondary"><RefreshCw size={16} /></Button>
          <Button><Plus size={16} style={{ marginRight: '8px' }} /> Create VPC</Button>
        </div>
      </header>

      <Table data={DUMMY_VPCS} columns={columns} />
    </div>
  );
}
