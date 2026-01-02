
'use client';

import React, { useEffect, useRef, useState } from 'react';
import gsap from 'gsap';
import { registerGSAP } from '@/animations/gsap';

const TIERS = [
  {
    name: "Starter",
    price: { monthly: "0", annual: "0" },
    desc: "For hobbyists and side projects.",
    features: ["500MB Telemetry / mo", "1 Edge Location", "Community Support", "1 User"]
  },
  {
    name: "Pro",
    price: { monthly: "29", annual: "24" },
    desc: "For scaling applications.",
    features: ["50GB Telemetry / mo", "Global Edge Network", "Email Support", "5 Users", "AI Triage (Beta)"]
  },
  {
    name: "Enterprise",
    price: { monthly: "Custom", annual: "Custom" },
    desc: "For mission-critical workloads.",
    features: ["Unlimited Telemetry", "Private Edge Nodes", "24/7 Phone Support", "SSO & Audit Logs", "Dedicated Solution Architect"]
  }
];

export const Pricing = () => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [annual, setAnnual] = useState(true);

  useEffect(() => {
    registerGSAP();
    const ctx = gsap.context(() => {
      const cards = gsap.utils.toArray('.pricing-card');
      
      gsap.from(cards, {
        y: 30,
        autoAlpha: 0,
        stagger: 0.1,
        duration: 0.8,
        ease: 'power2.out',
        id: 'tlPricingIn'
      });
    }, containerRef);
    
    return () => ctx.revert();
  }, []);

  return (
    <section ref={containerRef} style={{ maxWidth: '1200px', margin: '0 auto', padding: '120px 20px', textAlign: 'center' }}>
      <h1 style={{ fontSize: '56px', fontWeight: 700, marginBottom: '24px', color: '#1D1D1F' }}>Simple pricing.</h1>
      
      {/* Toggle */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '16px', marginBottom: '60px' }}>
        <span style={{ color: !annual ? '#1D1D1F' : '#86868B' }}>Monthly</span>
        <button 
          onClick={() => setAnnual(!annual)}
          style={{
            width: '48px',
            height: '28px',
            background: 'rgba(0,0,0,0.05)',
            borderRadius: '99px',
            position: 'relative',
            border: 'none',
            cursor: 'pointer'
          }}
        >
          <div style={{
            width: '24px',
            height: '24px',
            background: '#007AFF',
            borderRadius: '50%',
            position: 'absolute',
            top: '2px',
            left: annual ? '22px' : '2px',
            transition: 'left 0.3s cubic-bezier(0.2, 0, 0, 1)'
          }} />
        </button>
        <span style={{ color: annual ? '#1D1D1F' : '#86868B' }}>Annual <span style={{ color: '#34C759', fontSize: '12px' }}>(-20%)</span></span>
      </div>

      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', 
        gap: '24px' 
      }}>
        {TIERS.map((tier, i) => (
          <div key={i} className="pricing-card" style={{
            background: i === 1 ? 'rgba(0, 122, 255, 0.04)' : 'rgba(255,255,255,0.5)',
            border: i === 1 ? '1px solid rgba(0, 122, 255, 0.3)' : '1px solid rgba(0,0,0,0.05)',
            borderRadius: '24px',
            padding: '32px',
            textAlign: 'left',
            position: 'relative',
            display: 'flex',
            flexDirection: 'column',
            boxShadow: '0 4px 20px rgba(0,0,0,0.02)'
          }}>
            <h3 style={{ fontSize: '20px', fontWeight: 600, color: '#1D1D1F', marginBottom: '8px' }}>{tier.name}</h3>
            <p style={{ color: '#86868B', fontSize: '14px', marginBottom: '32px' }}>{tier.desc}</p>
            
            <div style={{ marginBottom: '32px' }}>
              <span style={{ fontSize: '42px', fontWeight: 700, color: '#1D1D1F' }}>
                {tier.price.monthly === "Custom" ? "Custom" : `$${annual ? tier.price.annual : tier.price.monthly}`}
              </span>
              {tier.price.monthly !== "Custom" && <span style={{ color: '#86868B' }}>/mo</span>}
            </div>

            <ul style={{ listStyle: 'none', padding: 0, margin: '0 0 40px 0', flex: 1 }}>
              {tier.features.map((feat, j) => (
                <li key={j} style={{ 
                  display: 'flex', 
                  alignItems: 'center', 
                  gap: '12px', 
                  color: '#424245', 
                  marginBottom: '16px',
                  fontSize: '15px'
                }}>
                  <span style={{ color: '#007AFF' }}>✓</span> {feat}
                </li>
              ))}
            </ul>

            <button style={{
              width: '100%',
              padding: '12px',
              borderRadius: '12px',
              border: 'none',
              background: i === 1 ? '#007AFF' : 'rgba(0,0,0,0.05)',
              color: i === 1 ? '#fff' : '#1D1D1F',
              fontWeight: 600,
              cursor: 'pointer'
            }}>
              {i === 2 ? 'Contact Sales' : 'Get Started'}
            </button>
          </div>
        ))}
      </div>
    </section>
  );
};
