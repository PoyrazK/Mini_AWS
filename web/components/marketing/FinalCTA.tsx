
'use client';

import React, { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { registerGSAP } from '@/animations/gsap';

export const FinalCTA = () => {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    registerGSAP();
    const ctx = gsap.context(() => {
      gsap.from(containerRef.current, {
        scrollTrigger: {
          trigger: containerRef.current,
          start: 'top 80%'
        },
        y: 30,
        autoAlpha: 0,
        duration: 0.8,
        ease: 'power2.out',
        id: 'tlFinalCTA'
      });
    }, containerRef);
    
    return () => ctx.revert();
  }, []);

  return (
    <section ref={containerRef} style={{ 
      padding: '120px 20px', 
      textAlign: 'center',
      background: 'linear-gradient(180deg, transparent 0%, rgba(0, 122, 255, 0.05) 100%)'
    }}>
      <h2 style={{ fontSize: '48px', fontWeight: 700, marginBottom: '24px', color: '#1D1D1F' }}>
        Ship faster. Sleep better.
      </h2>
      <p style={{ fontSize: '20px', color: '#86868B', marginBottom: '40px' }}>
        Start monitoring your edge infrastructure today.
      </p>

       <div style={{ display: 'flex', gap: '24px', justifyContent: 'center' }}>
          <button style={{
            background: '#000000',
            color: '#FFFFFF',
            padding: '16px 32px',
            borderRadius: '99px',
            fontSize: '18px',
            fontWeight: 600,
            border: 'none',
            cursor: 'pointer',
            boxShadow: '0 0 30px rgba(0,0,0,0.1)'
          }}>
            Request access
          </button>
      </div>
      
      <footer style={{ marginTop: '120px', color: '#424245', fontSize: '14px' }}>
        © 2026 EdgePulse. All rights reserved.
      </footer>
    </section>
  );
};
