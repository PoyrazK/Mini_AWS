
import '../globals.css';

export default function MarketingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div style={{ 
      backgroundColor: '#FFFFFF', 
      color: '#1D1D1F',
      minHeight: '100vh',
      fontFamily: '-apple-system, BlinkMacSystemFont, sans-serif'
    }}>
      {children}
    </div>
  );
}
