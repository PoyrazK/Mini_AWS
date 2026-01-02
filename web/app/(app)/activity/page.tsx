'use client';

import { Table, Column } from '@/components/ui/Table';
import { Button } from '@/components/ui/Button';
import { Download, RefreshCw } from 'lucide-react';

interface Event {
  id: string;
  action: string;
  resource: string;
  user: string;
  status: 'success' | 'failure';
  timestamp: string;
}

const DUMMY_EVENTS: Event[] = [
  { id: 'evt-1001', action: 'RunInstances', resource: 'i-0x8231', user: 'root', status: 'success', timestamp: '2025-01-14 10:42:01' },
  { id: 'evt-1002', action: 'CreateBucket', resource: 'logs-archive', user: 'admin', status: 'success', timestamp: '2025-01-14 09:15:33' },
  { id: 'evt-1003', action: 'StopInstances', resource: 'i-0x11b2', user: 'root', status: 'success', timestamp: '2025-01-13 18:20:00' },
  { id: 'evt-1004', action: 'AttachVolume', resource: 'vol-0x555', user: 'system', status: 'failure', timestamp: '2025-01-13 18:19:45' },
];

export default function ActivityPage() {
  const columns: Column<Event>[] = [
    { header: 'Event Name', accessorKey: 'action', width: '25%' },
    { header: 'Resource', accessorKey: 'resource', width: '20%' },
    { header: 'User', accessorKey: 'user' },
    { 
      header: 'Status', 
      cell: (item) => (
        <span style={{ 
          color: item.status === 'success' ? 'var(--accent-green)' : 'var(--accent-red)',
          fontWeight: 500
        }}>
          {item.status.toUpperCase()}
        </span>
      )
    },
    { header: 'Timestamp', accessorKey: 'timestamp' },
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
           <h1 style={{ fontSize: '34px', fontWeight: 700, marginBottom: '4px', letterSpacing: '0.01em', color: 'var(--text-primary)' }}>Activity</h1>
           <p style={{ color: 'var(--text-secondary)' }}>Audit logs and system events.</p>
        </div>
        <div style={{ display: 'flex', gap: '12px' }}>
          <Button variant="secondary"><RefreshCw size={16} /></Button>
          <Button variant="secondary"><Download size={16} style={{ marginRight: '8px' }} /> Export CSV</Button>
        </div>
      </header>

      <Table data={DUMMY_EVENTS} columns={columns} />
    </div>
  );
}
