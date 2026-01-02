
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

export const registerGSAP = () => {
  if (typeof window !== 'undefined') {
    gsap.registerPlugin(ScrollTrigger);
    
    // Apple-style config
    gsap.defaults({
      ease: 'power2.out',
      duration: 0.6,
    });
  }
};

export const createScrollTrigger = (
  trigger: Element | string, 
  vars: ScrollTrigger.Vars
) => {
  return ScrollTrigger.create({
    trigger,
    ...vars
  });
};
