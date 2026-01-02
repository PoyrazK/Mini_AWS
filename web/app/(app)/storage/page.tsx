'use client';

import { Table, Column } from '@/components/ui/Table';
import { Button } from '@/components/ui/Button';
import { Plus, RefreshCw, HardDrive } from 'lucide-react';

interface Bucket {
  name: string;
  region: string;
  objects: number;
  size: string;
  created_at: string;
}

const DUMMY_BUCKETS: Bucket[] = [
  { name: 'assets-prod-v1', region: 'us-east-1', objects: 1240, size: '4.2 GB', created_at: '2024-12-01' },
  { name: 'user-uploads', region: 'us-east-1', objects: 8502, size: '156 GB', created_at: '2024-12-15' },
  { name: 'logs-archive', region: 'us-west-2', objects: 450, size: '240 MB', created_at: '2025-01-02' },
];

export default function StoragePage() {
  const columns: Column<Bucket>[] = [
    { 
      header: 'Name', 
      cell: (item) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontWeight: 500 }}>
          <HardDrive size={16} color="var(--accent-blue)" />
          {item.name}
        </div>
      )
    },
    { header: 'Region', accessorKey: 'region' },
    { header: 'Objects', accessorKey: 'objects' },
    { header: 'Size', accessorKey: 'size' },
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
           <h1 style={{ fontSize: '34px', fontWeight: 700, marginBottom: '4px', letterSpacing: '0.01em', color: 'var(--text-primary)' }}>Storage</h1>
           <p style={{ color: 'var(--text-secondary)' }}>S3-compatible object storage.</p>
        </div>
        <div style={{ display: 'flex', gap: '12px' }}>
          <Button variant="secondary"><RefreshCw size={16} /></Button>
          <Button><Plus size={16} style={{ marginRight: '8px' }} /> Create Bucket</Button>
        </div>
      </header>

      <Table data={DUMMY_BUCKETS} columns={columns} />
    </div>
  );
}
