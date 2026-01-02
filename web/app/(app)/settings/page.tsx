
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { User, Key, Globe } from 'lucide-react';

export default function SettingsPage() {
  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <header style={{ marginBottom: '32px' }}>
        <h1 style={{ fontSize: '34px', fontWeight: 700, marginBottom: '4px', letterSpacing: '0.01em', color: 'var(--text-primary)' }}>Settings</h1>
        <p style={{ color: 'var(--text-secondary)' }}>Manage your account and preferences.</p>
      </header>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
        <Card title="Account Profile">
          <div style={{ display: 'flex', alignItems: 'center', gap: '20px' }}>
            <div style={{
              width: '64px', height: '64px',
              borderRadius: '50%',
              background: 'var(--system-gray-5)',
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              color: 'var(--text-secondary)'
            }}>
              <User size={32} />
            </div>
            <div>
              <div style={{ fontSize: '18px', fontWeight: 600 }}>Root User</div>
              <div style={{ color: 'var(--text-secondary)' }}>root@thecloud.local</div>
            </div>
            <Button variant="secondary" style={{ marginLeft: 'auto' }}>Edit</Button>
          </div>
        </Card>

        <Card title="API Access">
          <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
            <p style={{ color: 'var(--text-secondary)', fontSize: '14px', lineHeight: '1.5' }}>
              Use these keys to access The Cloud via the CLI or SDKs. Do not share your secret key.
            </p>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '12px', background: 'var(--system-gray-6)', borderRadius: '8px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                <Key size={16} color="var(--text-secondary)" />
                <code style={{ fontSize: '13px', fontFamily: 'Menlo, monospace' }}>AKIA-THE-CLOUD-DEMO-KEY</code>
              </div>
              <span style={{ fontSize: '12px', color: 'var(--accent-green)', fontWeight: 500 }}>Active</span>
            </div>
            <div>
              <Button>Generate New Key</Button>
            </div>
          </div>
        </Card>

        <Card title="Region">
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <Globe size={18} color="var(--text-secondary)" />
              <div>
                <div style={{ fontWeight: 500 }}>US East (N. Virginia)</div>
                <div style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>us-east-1</div>
              </div>
            </div>
            <Button variant="secondary">Change</Button>
          </div>
        </Card>
      </div>
    </div>
  );
}
