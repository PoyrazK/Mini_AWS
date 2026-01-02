import { Features } from '@/components/marketing/Features';
import { ProductStory } from '@/components/marketing/ProductStory';
import { EdgeMap } from '@/components/marketing/EdgeMap';
import { FinalCTA } from '@/components/marketing/FinalCTA';

import { Hero } from '@/components/marketing/Hero';
import { SignalStrip } from '@/components/marketing/SignalStrip';

export const runtime = 'edge';

export default function EdgePulseHome() {
  return (
    <main>
      <Hero />
      <SignalStrip />
      <Features />
      <ProductStory />
      <EdgeMap />
      <FinalCTA />
    </main>
  );
}
