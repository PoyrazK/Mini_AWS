
'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { LayoutGrid, Server, HardDrive, Network, Settings, Activity } from 'lucide-react';
import styles from './Sidebar.module.css';

const MENU_ITEMS = [
  { name: 'Dashboard', icon: LayoutGrid, href: '/dashboard' },
  { name: 'Compute', icon: Server, href: '/compute' },
  { name: 'Storage', icon: HardDrive, href: '/storage' },
  { name: 'Network', icon: Network, href: '/network' },
  { name: 'Activity', icon: Activity, href: '/activity' },
  { name: 'Settings', icon: Settings, href: '/settings' },
];

export const Sidebar: React.FC = () => {
  const pathname = usePathname();

  return (
    <aside className={`${styles.sidebar} material-sidebar`}>
      <div className={styles.logo}>
        <div className={styles.logoIcon}>☁️</div>
        <span className={styles.logoText}>Mini AWS</span>
      </div>
      
      <nav className={styles.nav}>
        {MENU_ITEMS.map((item) => {
          const Icon = item.icon;
          const isActive = pathname === item.href;
          
          return (
            <Link 
              key={item.name} 
              href={item.href}
              className={`${styles.navItem} ${isActive ? styles.active : ''}`}
            >
              <Icon size={16} strokeWidth={2} style={{ opacity: isActive ? 1 : 0.7 }} />
              <span>{item.name}</span>
            </Link>
          );
        })}
      </nav>
      
      <div className={styles.footer}>
        <div className={styles.status}>
          <div className={styles.statusDot} />
          <span>us-east-1</span>
        </div>
      </div>
    </aside>
  );
};
