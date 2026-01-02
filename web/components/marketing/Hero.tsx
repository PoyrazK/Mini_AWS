
'use client';

import React, { useEffect, useRef } from 'react';
import Link from 'next/link';
import gsap from 'gsap';
import { registerGSAP } from '@/animations/gsap';
import { MagneticButton } from '@/components/marketing/MagneticButton';
import { WorkerBackground } from '@/components/marketing/WorkerBackground';

export const Hero = () => {
  const containerRef = useRef<HTMLDivElement>(null);
  const titleLine1Ref = useRef<HTMLHeadingElement>(null);
  const titleLine2Ref = useRef<HTMLHeadingElement>(null);
  const subRef = useRef<HTMLParagraphElement>(null);
  const ctaRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    registerGSAP();
    const ctx = gsap.context(() => {
      const tl = gsap.timeline({ id: 'tlHeroIntro' });
      
      // 1. Bg glow fade in
      tl.fromTo('.hero-glow', 
        { autoAlpha: 0, scale: 0.8 }, 
        { autoAlpha: 0.4, scale: 1, duration: 1.8, ease: 'power1.out' },
        0
      );

      // 2. Headline stagger
      tl.from([titleLine1Ref.current, titleLine2Ref.current], {
        y: 24,
        autoAlpha: 0,
        duration: 1.0,
        stagger: 0.1,
        ease: 'expo.out'
      }, 0.2);

      // 3. Subtitle
      tl.from(subRef.current, {
        y: 12,
        autoAlpha: 0,
        duration: 0.8,
        ease: 'power2.out'
      }, 0.5);

      // 4. CTA
      tl.from(ctaRef.current, {
        y: 8,
        autoAlpha: 0,
        scale: 0.98,
        duration: 0.6,
        ease: 'power2.out'
      }, 0.7);

    }, containerRef);
    
    return () => ctx.revert();
  }, []);

  return (
    <section ref={containerRef} style={{ 
      position: 'relative', 
      minHeight: '85vh', 
      display: 'flex', 
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center',
      textAlign: 'center',
      padding: '0 20px',
      overflow: 'hidden'
    }}>
      {/* Background Glow */}
      <WorkerBackground />

      <div style={{ position: 'relative', zIndex: 1, maxWidth: '800px' }}>
        <h1 style={{ 
          fontSize: 'clamp(40px, 8vw, 80px)', 
          fontWeight: 700, 
          letterSpacing: '-0.02em',
          lineHeight: 1.05,
          marginBottom: '24px',
          color: '#1D1D1F'
        }}>
          <div ref={titleLine1Ref}>Edge observability</div>
          <div ref={titleLine2Ref} style={{ color: '#007AFF' }}>that feels instant.</div>
        </h1>

        <p ref={subRef} style={{ 
          fontSize: '20px', 
          lineHeight: 1.5,
          color: '#86868B',
          maxWidth: '580px',
          margin: '0 auto 40px'
        }}>
          EdgePulse streams telemetry, detects anomalies, and guides incident response — before users notice.
        </p>

        <div ref={ctaRef} style={{ display: 'flex', gap: '24px', justifyContent: 'center', alignItems: 'center' }}>
          <MagneticButton style={{
            background: '#000000',
            color: '#FFFFFF',
            padding: '14px 28px',
            borderRadius: '99px',
            fontSize: '16px',
            fontWeight: 600,
            border: 'none'
          }}>
            Request access
          </MagneticButton>
          
          <Link href="/pricing" style={{ 
            color: '#1D1D1F', 
            textDecoration: 'none',
            fontSize: '16px',
            fontWeight: 500
          }}>
            View pricing →
          </Link>
        </div>
      </div>
    </section>
  );
};
