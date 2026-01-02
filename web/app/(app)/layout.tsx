
import { Sidebar } from '@/components/ui/Sidebar';

export default function AppLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <Sidebar />
      <main style={{ 
        flex: 1, 
        marginLeft: 'var(--sidebar-width)', 
        padding: '32px',
        overflowY: 'auto'
      }}>
        {children}
      </main>
    </div>
  );
}
