
'use client';

import React, { useRef } from 'react';
import gsap from 'gsap';

interface MagneticButtonProps {
  children: React.ReactNode;
  onClick?: () => void;
  strength?: number; // 0.2 to 0.8 usually
  style?: React.CSSProperties;
}

export const MagneticButton: React.FC<MagneticButtonProps> = ({ 
  children, 
  onClick, 
  strength = 0.35,
  style 
}) => {
  const ref = useRef<HTMLButtonElement>(null);
  
  const handleMouseMove = (e: React.MouseEvent) => {
    if (!ref.current) return;
    
    const { clientX, clientY } = e;
    const { left, top, width, height } = ref.current.getBoundingClientRect();
    
    // Calculate center
    const centerX = left + width / 2;
    const centerY = top + height / 2;
    
    // Calculate distance from center
    const x = (clientX - centerX) * strength;
    const y = (clientY - centerY) * strength;
    
    // Animate button towards mouse
    gsap.to(ref.current, {
      x,
      y,
      duration: 1,
      ease: 'power4.out'
    });
    
    // Animate content/text slightly more for parallax effect (optional, maybe keep simple first)
  };
  
  const handleMouseLeave = () => {
    if (!ref.current) return;
    
    // Snap back
    gsap.to(ref.current, {
      x: 0,
      y: 0,
      duration: 1,
      ease: 'elastic.out(1, 0.3)'
    });
  };

  return (
    <button
      ref={ref}
      onClick={onClick}
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
      style={{
        cursor: 'pointer',
        ...style
      }}
    >
      {children}
    </button>
  );
};
