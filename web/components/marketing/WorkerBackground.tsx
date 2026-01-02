
'use client';

import React, { useEffect, useRef } from 'react';

export const WorkerBackground = () => {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    let w = (canvas.width = window.innerWidth);
    let h = (canvas.height = window.innerHeight);
    let frames = 0;

    const resize = () => {
      w = canvas.width = window.innerWidth;
      h = canvas.height = window.innerHeight;
    };

    window.addEventListener('resize', resize);

    // Grid configuration
    const gridSize = 40;

    const speed = 0.5;
    let offset = 0;

    const draw = () => {
      if (!ctx) return;
      frames++;
      offset = (offset + speed) % gridSize;

      // Clear with trail effect for "motion blur" feel
      ctx.fillStyle = 'rgba(255, 255, 255, 1)'; 
      ctx.fillRect(0, 0, w, h);

      ctx.save();
      // Center origin
      ctx.translate(w / 2, h / 2);

      // We will draw a "floor" grid that moves towards the viewer
      // Simple perspective projection
      
      ctx.strokeStyle = 'rgba(0, 122, 255, 0.1)'; // Subtle blue grid
      ctx.lineWidth = 1;

      // Vertical lines (diverging)
      for (let x = -w; x < w; x += gridSize * 2) {
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x * 4, h); // Perspective flare
        ctx.stroke();
      }

      // Horizontal lines (moving down)
      for (let y = 0; y < h; y += gridSize) {
        const perspectiveY = (y + offset) * 1.5; // Exponential spacing for 3D feel? Or just linear for retro
        const alpha = Math.max(0, 1 - perspectiveY / (h / 2)); // Fade out at bottom
        
        ctx.globalAlpha = alpha;
        ctx.beginPath();
        // Simple horizon line logic
        ctx.moveTo(-w, perspectiveY);
        ctx.lineTo(w, perspectiveY);
        ctx.stroke();
        ctx.globalAlpha = 1;
      }

      // Glowing "Packets"
      // Randomly spawn or move dots along the vertical lines
      const t = frames * 0.05;
      for (let i = 0; i < 5; i++) {
        const lineIndex = Math.floor(Math.sin(t + i) * 10);
        const x = lineIndex * gridSize * 2;
        const y = ((frames * 5 + i * 100) % h) * 1.5;
        
        if (y < h && y > 0) {
            ctx.shadowBlur = 10;
            ctx.shadowColor = '#007AFF';
            ctx.fillStyle = '#FFFFFF';
            ctx.beginPath();
            ctx.arc(x, y, 3, 0, Math.PI * 2);
            ctx.fill();
            ctx.shadowBlur = 0;
        }
      }

      ctx.restore();
      requestAnimationFrame(draw);
    };

    const animationId = requestAnimationFrame(draw);

    return () => {
        window.removeEventListener('resize', resize);
        cancelAnimationFrame(animationId);
    };
  }, []);

  return (
    <canvas 
      ref={canvasRef}
      style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
        pointerEvents: 'none',
        zIndex: 0,
        background: 'white' // Force white background immediately
      }}
    />
  );
};
