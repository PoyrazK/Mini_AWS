
'use client';

import React, { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { registerGSAP } from '@/animations/gsap';

const STEPS = [
  {
    id: 'capture',
    title: "Capture",
    body: "Telemetry is ingested at the edge. No data fees, no egress costs. Just pure, raw signal processed in main memory.",
  },
  {
    id: 'understand',
    title: "Understand",
    body: "Our AI models analyze patterns in real-time, distinguishing between noise and true anomalies with 99.9% accuracy.",
  },
  {
    id: 'resolve',
    title: "Resolve",
    body: "Automated playbooks trigger instantly. Reroute traffic, rollback deploys, or scale capacity without human intervention.",
  }
];

export const ProductStory = () => {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    registerGSAP();
    const ctx = gsap.context(() => {
      // SVG Line Animation
      const svgPath = document.querySelector('.story-line-path') as SVGPathElement;
      if (svgPath) {
        const length = svgPath.getTotalLength();
        gsap.set(svgPath, { strokeDasharray: length, strokeDashoffset: length });
        
        gsap.to(svgPath, {
          strokeDashoffset: 0,
          scrollTrigger: {
            trigger: containerRef.current,
            start: 'top center',
            end: 'bottom center',
            scrub: 1
          },
          ease: 'none'
        });
      }

      // Existing Stagger Logic
      STEPS.forEach((step) => {
        const section = document.getElementById(`story-${step.id}`);
        if (!section) return;

        const tl = gsap.timeline({
          scrollTrigger: {
            trigger: section,
            start: 'top 70%',
            end: 'bottom 70%',
            toggleActions: 'play reverse play reverse'
          },
          id: `tlStory${step.title}`
        });

        tl.from(section.querySelector('h2'), { y: 20, autoAlpha: 0, duration: 0.55, ease: 'power2.out' })
          .from(section.querySelector('p'), { y: 20, autoAlpha: 0, duration: 0.45, ease: 'power2.out' }, "-=0.3")
          .from(section.querySelector('.visual-node'), { scale: 0, autoAlpha: 0, duration: 0.9, ease: 'elastic.out(1, 0.5)' }, "-=0.4");
      });
    }, containerRef);
    
    return () => ctx.revert();
  }, []);

  return (
    <section ref={containerRef} style={{ padding: '100px 20px', maxWidth: '1000px', margin: '0 auto', position: 'relative' }}>
      
      {/* Connector Line SVG Overlay */}
      <svg style={{
        position: 'absolute',
        top: 0,
        left: '50%',
        transform: 'translateX(-50%)',
        width: '2px', // Thin strip down the middle
        height: '100%',
        overflow: 'visible',
        zIndex: 0,
        pointerEvents: 'none'
      }}>
         {/* Simple vertical line for mobile/desktop commonality */}
         <path 
           className="story-line-path"
           d="M 1 50 L 1 1200" 
           fill="none" 
           stroke="rgba(0, 0, 0, 0.1)" 
           strokeWidth="2"
           strokeLinecap="round"
           vectorEffect="non-scaling-stroke"
         />
      </svg>

      {STEPS.map((step, i) => (
        <div key={step.id} id={`story-${step.id}`} style={{
          display: 'flex',
          flexDirection: i % 2 === 0 ? 'row' : 'row-reverse',
          alignItems: 'center',
          gap: '80px',
          marginBottom: '160px',
          minHeight: '400px',
          position: 'relative',
          zIndex: 1
        }}>
          <div style={{ flex: 1 }}>
            <h2 style={{ fontSize: '32px', fontWeight: 700, marginBottom: '20px', color: '#1D1D1F' }}>{step.title}</h2>
            <p style={{ fontSize: '18px', lineHeight: 1.6, color: '#86868B' }}>{step.body}</p>
          </div>
          
          <div style={{ flex: 1, display: 'flex', justifyContent: 'center' }}>
            {/* Abstract Visual */}
            <div className="visual-node" style={{
              width: '280px',
              height: '280px',
              background: 'rgba(255,255,255,0.03)',
              borderRadius: '50%',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              position: 'relative'
            }}>
              <div style={{
                width: '120px',
                height: '120px',
                background: `radial-gradient(circle, ${i === 0 ? '#007AFF' : i === 1 ? '#30D158' : '#FF9500'} 0%, transparent 70%)`,
                opacity: 0.4,
                borderRadius: '50%',
                filter: 'blur(20px)'
              }} />
              <div style={{
                position: 'absolute',
                width: '100%',
                height: '1px',
                background: 'linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent)'
              }} />
              <div style={{
                position: 'absolute',
                height: '100%',
                width: '1px',
                background: 'linear-gradient(180deg, transparent, rgba(255,255,255,0.2), transparent)'
              }} />
            </div>
          </div>
        </div>
      ))}
    </section>
  );
};
