
'use client';

import React, { useRef, useState } from 'react';
import gsap from 'gsap';
import { useGSAP } from '@gsap/react';
import { registerGSAP } from '@/animations/gsap';

const FEATURE_ITEMS = [
  {
    title: "Instant Edge Telemetry",
    desc: "Collect metrics from the edge with sub-50ms latency. No cold starts, no data loss.",
    span: "col-span-2"
  },
  {
    title: "AI Incident Triage",
    desc: "Automatically correlate anomalies with deploy events.",
    span: "col-span-1"
  },
  {
    title: "Global Reliability",
    desc: "Distributed across 300+ edge locations.",
    span: "col-span-1"
  },
  {
    title: "Zero Overhead",
    desc: "Lightweight agents that barely touch your CPU budget.",
    span: "col-span-2"
  }
];

const SpotlightCard = ({ children, className = "" }: { children: React.ReactNode, className?: string }) => {
  const divRef = useRef<HTMLDivElement>(null);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [opacity, setOpacity] = useState(0);

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!divRef.current) return;
    const rect = divRef.current.getBoundingClientRect();
    setPosition({ x: e.clientX - rect.left, y: e.clientY - rect.top });
  };

  const handleMouseEnter = () => setOpacity(1);
  const handleMouseLeave = () => setOpacity(0);

  return (
    <div
      ref={divRef}
      onMouseMove={handleMouseMove}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      className={className}
      style={{
        position: 'relative',
        overflow: 'hidden',
        borderRadius: '24px',
        border: '1px solid rgba(0, 0, 0, 0.05)',
        backgroundColor: 'rgba(255, 255, 255, 0.6)',
        padding: '32px',
        backdropFilter: 'blur(20px)',
        WebkitBackdropFilter: 'blur(20px)',
        boxShadow: '0 4px 20px rgba(0,0,0,0.02)'
      }}
    >
      <div
        style={{
          pointerEvents: 'none',
          position: 'absolute',
          top: -1, left: -1, right: -1, bottom: -1,
          opacity,
          transition: 'opacity 300ms',
          background: `radial-gradient(600px circle at ${position.x}px ${position.y}px, rgba(0, 122, 255, 0.1), transparent 40%)`
        }}
      />
      <div style={{ position: 'relative', zIndex: 10 }}>{children}</div>
    </div>
  );
};

export const Features = () => {
  const containerRef = useRef<HTMLDivElement>(null);

  useGSAP(() => {
    registerGSAP();
    const cards = gsap.utils.toArray('.bento-card');
    
    gsap.from(cards, {
      scrollTrigger: {
        trigger: containerRef.current,
        start: 'top 80%',
        toggleActions: 'play none none reverse'
      },
      y: 50,
      autoAlpha: 0,
      stagger: 0.1,
      duration: 0.8,
      ease: 'power2.out'
    });
  }, { scope: containerRef });

  return (
    <section ref={containerRef} style={{ padding: '120px 20px', maxWidth: '1000px', margin: '0 auto' }}>
      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(3, 1fr)', 
        gap: '24px' 
      }}>
        {FEATURE_ITEMS.map((item, i) => (
          <SpotlightCard 
            key={i} 
            className={`bento-card ${item.span === 'col-span-2' ? 'col-span-2-styles' : ''}`} 
          >
            {/* Hardcoded style injection for grid spans since we aren't using Tailwind */}
            <style jsx>{`
              .col-span-2-styles { grid-column: span 2; }
              @media (max-width: 768px) {
                .col-span-2-styles { grid-column: span 1 !important; }
                div[style*="display: grid"] { grid-template-columns: 1fr !important; }
              }
            `}</style>
            
            <h3 style={{ fontSize: '24px', fontWeight: 600, marginBottom: '12px', color: '#1D1D1F' }}>{item.title}</h3>
            <p style={{ fontSize: '16px', lineHeight: 1.6, color: '#86868B' }}>{item.desc}</p>
          </SpotlightCard>
        ))}
      </div>
    </section>
  );
};
