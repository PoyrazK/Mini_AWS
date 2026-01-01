
'use client';

import React, { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { registerGSAP } from '@/animations/gsap';

export const EdgeMap = () => {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    registerGSAP();
    const ctx = gsap.context(() => {
      gsap.fromTo('.map-grid', 
        { autoAlpha: 0, scale: 0.95 },
        { 
          autoAlpha: 1, 
          scale: 1, 
          duration: 1.5, 
          ease: 'power2.out',
          scrollTrigger: {
            trigger: containerRef.current,
            start: 'top 70%'
          },
          id: 'tlEdgeMap'
        }
      );

      gsap.to('.map-node', {
        scale: 1.5,
        opacity: 0,
        duration: 2.0,
        stagger: {
          each: 0.5,
          repeat: -1
        },
        ease: 'power1.out'
      });
    }, containerRef);
    
    return () => ctx.revert();
  }, []);

  return (
    <section ref={containerRef} style={{ 
      padding: '100px 20px', 
      textAlign: 'center', 
      position: 'relative',
      overflow: 'hidden' 
    }}>
      <h2 style={{ fontSize: '32px', fontWeight: 700, marginBottom: '60px', color: '#1D1D1F' }}>Global by default.</h2>
      
      <div className="map-grid" style={{
        maxWidth: '1000px',
        height: '500px',
        margin: '0 auto',
        background: 'rgba(255,255,255,0.5)',
        borderRadius: '24px',
        border: '1px solid rgba(0,0,0,0.05)',
        position: 'relative',
        backgroundImage: 'radial-gradient(rgba(0,0,0,0.1) 1px, transparent 1px)',
        backgroundSize: '40px 40px',
        boxShadow: 'inset 0 0 40px rgba(0,0,0,0.02)'
      }}>
        {/* Animated Nodes representing Edge Locations */}
        {[
          { top: '30%', left: '20%' }, // US West
          { top: '35%', left: '28%' }, // US East
          { top: '25%', left: '48%' }, // Europe
          { top: '45%', left: '75%' }, // Asia
          { top: '65%', left: '85%' }, // Australia
        ].map((pos, i) => (
          <div key={i} className="map-node" style={{
            position: 'absolute',
            top: pos.top,
            left: pos.left,
            width: '12px',
            height: '12px',
            background: '#007AFF',
            borderRadius: '50%',
            boxShadow: '0 0 10px #007AFF'
          }} />
        ))}
        {/* Static center dots */}
       {[
          { top: '30%', left: '20%' },
          { top: '35%', left: '28%' },
          { top: '25%', left: '48%' },
          { top: '45%', left: '75%' },
          { top: '65%', left: '85%' },
        ].map((pos, i) => (
          <div key={i} style={{
            position: 'absolute',
            top: pos.top,
            left: pos.left,
            width: '12px',
            height: '12px',
            background: '#007AFF',
            borderRadius: '50%',
          }} />
        ))}
      </div>
    </section>
  );
};
