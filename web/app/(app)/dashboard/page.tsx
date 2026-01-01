
import { Card } from '@/components/ui/Card';
import { StatusIndicator } from '@/components/ui/StatusIndicator';
import { Button } from '@/components/ui/Button';
import Link from 'next/link';
import { Activity, Server, HardDrive, Cpu } from 'lucide-react';

export default function Home() {
  return (
    <div style={{ maxWidth: '1280px', margin: '0 auto' }}>
      <header style={{ 
        marginBottom: '32px', 
        paddingTop: '20px',
        position: 'sticky',
        top: 0,
        zIndex: 10
      }}>
        <h1 style={{ 
          fontSize: '34px', 
          fontWeight: 700, 
          marginBottom: '4px',
          letterSpacing: '0.01em',
          color: 'var(--text-primary)'
        }}>
          Dashboard
        </h1>
        <p style={{ color: 'var(--text-secondary)', fontSize: '15px' }}>Overview</p>
      </header>
  
      {/* Metrics Grid */}
      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', 
        gap: '20px',
        marginBottom: '32px'
      }}>
        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
            <div style={{ 
              padding: '12px', 
              borderRadius: '50%', /* Circular icon backgrounds */
              background: 'rgba(10, 132, 255, 0.15)',
              color: 'var(--accent-blue)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: '48px',
              height: '48px'
            }}>
              <Server size={22} strokeWidth={2.5} />
            </div>
            <div>
              <div style={{ fontSize: '26px', fontWeight: 600, lineHeight: 1 }}>12</div>
              <div style={{ color: 'var(--text-secondary)', fontSize: '13px', marginTop: '4px', fontWeight: 500 }}>Active Instances</div>
            </div>
          </div>
        </Card>
        
        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
             <div style={{ 
              padding: '12px', 
               borderRadius: '50%',
              background: 'rgba(48, 209, 88, 0.15)',
              color: 'var(--accent-green)',
               display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: '48px',
              height: '48px'
            }}>
              <Activity size={22} strokeWidth={2.5} />
            </div>
            <div>
              <div style={{ fontSize: '26px', fontWeight: 600, lineHeight: 1 }}>98.2%</div>
               <div style={{ color: 'var(--text-secondary)', fontSize: '13px', marginTop: '4px', fontWeight: 500 }}>Healthy Services</div>
            </div>
          </div>
        </Card>

        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
             <div style={{ 
              padding: '12px', 
               borderRadius: '50%',
              background: 'rgba(255, 159, 10, 0.15)',
              color: 'var(--accent-orange)',
               display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: '48px',
              height: '48px'
            }}>
              <HardDrive size={22} strokeWidth={2.5} />
            </div>
            <div>
              <div style={{ fontSize: '26px', fontWeight: 600, lineHeight: 1 }}>450 GB</div>
               <div style={{ color: 'var(--text-secondary)', fontSize: '13px', marginTop: '4px', fontWeight: 500 }}>Storage Used</div>
            </div>
          </div>
        </Card>
        
        <Card>
           <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
             <div style={{ 
              padding: '12px', 
               borderRadius: '50%',
              background: 'rgba(255, 69, 58, 0.15)',
              color: 'var(--accent-red)',
               display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: '48px',
              height: '48px'
            }}>
              <Cpu size={22} strokeWidth={2.5} />
            </div>
            <div>
              <div style={{ fontSize: '26px', fontWeight: 600, lineHeight: 1 }}>45%</div>
               <div style={{ color: 'var(--text-secondary)', fontSize: '13px', marginTop: '4px', fontWeight: 500 }}>CPU Load</div>
            </div>
          </div>
        </Card>
      </div>

      {/* Recent Activity & Resources */}
      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: '2fr 1fr', 
        gap: '24px' 
      }}>
        <Card title="Recent Activity" style={{ minHeight: '400px' }}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
            {[1, 2, 3].map((i) => (
              <div key={i} style={{ 
                display: 'flex', 
                alignItems: 'center', 
                justifyContent: 'space-between',
                padding: '12px 0',
                borderBottom: i < 3 ? '1px solid var(--glass-border)' : 'none'
              }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                   <StatusIndicator status="running" />
                   <div>
                     <div style={{ fontWeight: 500 }}>Instance i-0x823 launched</div>
                     <div style={{ fontSize: '12px', color: 'var(--text-secondary)' }}>2 minutes ago</div>
                   </div>
                </div>
                <Button variant="ghost" size="sm">View</Button>
              </div>
            ))}
          </div>
        </Card>

        <Card title="Quick Actions">
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <Link href="/compute">
              <Button variant="primary" style={{ width: '100%' }}>Launch Instance</Button>
            </Link>
            <Link href="/storage">
              <Button variant="secondary" style={{ width: '100%' }}>Create Bucket</Button>
            </Link>
            <Link href="/activity">
               <Button variant="ghost" style={{ width: '100%' }}>View Activity</Button>
            </Link>
          </div>
        </Card>
      </div>
    </div>
  );
}
