
'use client';

import React, { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { registerGSAP } from '@/animations/gsap';

const METRICS = [
  "Sub-50ms insight",
  "99.99% edge uptime",
  "Auto-triage in minutes"
];

export const SignalStrip = () => {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    registerGSAP();
    const ctx = gsap.context(() => {
      const items = gsap.utils.toArray('.metric-item');
      
      gsap.from(items, {
        scrollTrigger: {
          trigger: containerRef.current,
          start: 'top 85%',
          toggleActions: 'play none none reverse'
        },
        y: 20,
        autoAlpha: 0,
        stagger: 0.1,
        duration: 0.6,
        ease: 'power2.out',
        id: 'tlMetricsReveal'
      });
    }, containerRef);
    
    return () => ctx.revert();
  }, []);

  return (
    <div ref={containerRef} style={{
      borderTop: '1px solid rgba(0,0,0,0.05)',
      borderBottom: '1px solid rgba(0,0,0,0.05)',
      padding: '24px 0',
      background: 'rgba(255,255,255,0.5)'
    }}>
      <div style={{
        maxWidth: '1000px',
        margin: '0 auto',
        display: 'flex',
        justifyContent: 'space-around',
        alignItems: 'center',
        flexWrap: 'wrap',
        gap: '20px'
      }}>
        {METRICS.map((text, i) => (
          <div key={i} className="metric-item" style={{
            color: '#1D1D1F',
            fontSize: '14px',
            fontWeight: 500,
            letterSpacing: '0.02em',
            display: 'flex',
            alignItems: 'center',
            gap: '8px'
          }}>
            <span style={{ 
              width: '6px', 
              height: '6px', 
              borderRadius: '50%', 
              background: '#007AFF',
              boxShadow: '0 0 8px #007AFF'
            }} />
            {text}
          </div>
        ))}
      </div>
    </div>
  );
};
