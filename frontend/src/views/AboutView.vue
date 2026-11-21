<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { Check, ShieldCheck } from 'lucide-vue-next';
import { useGsapScope } from '@/composables/useGsapScope';
import { assetUrl } from '@/utils/assetUrl';

const router = useRouter();
const pageRef = ref<HTMLElement | null>(null);

useGsapScope(pageRef, ({ gsap }) => {
  gsap.from('.about-hero-anim', { y: 26, opacity: 0, duration: 0.7, ease: 'power3.out' });
  gsap.from('.about-stat-anim', {
    y: 18,
    opacity: 0,
    duration: 0.5,
    stagger: 0.06,
    ease: 'power3.out',
    scrollTrigger: { trigger: '.about-story-grid', start: 'top 82%', once: true },
  });
  gsap.from('.about-pillar-anim', {
    y: 24,
    opacity: 0,
    duration: 0.55,
    stagger: 0.08,
    ease: 'power3.out',
    scrollTrigger: { trigger: '.about-pillars-grid', start: 'top 85%', once: true },
  });
});

function goToEstimate(): void {
  void router.push({ path: '/', hash: '#schedule' });
}

const pillars = [
  {
    title: 'Laser-verified substrates',
    body: 'Before any porcelain or limestone slab leaves its crate, we map every high spot, deflection joint, and moisture reading with digital instrumentation.',
  },
  {
    title: 'In-house licensed setters',
    body: 'We never hand your home off to random day-labor crews. Every Copely lead setter has at least seven years laying custom stone and bonded wet rooms.',
  },
  {
    title: 'Written 25-year guarantee',
    body: 'Your final walkthrough includes photos of every hidden waterproofing seam plus a transferable 25-year installation warranty in writing.',
  },
];

const stats = [
  { value: '2009', label: 'Founded in the Hill Country' },
  { value: '1,200+', label: 'Verified 5-star reviews' },
  { value: '25 yrs', label: 'Transferable warranty' },
  { value: '2–4 days', label: 'Average room turnaround' },
];
</script>

<template>
  <section ref="pageRef" class="pt-28 pb-16 px-4 md:px-10 bg-[#F7F7F3]">
    <div class="max-w-[1320px] mx-auto space-y-16 md:space-y-24">
      <!-- Full-Width Panoramic Sunset Custom Home Hero -->
      <div class="about-hero-anim relative rounded-3xl overflow-hidden bg-[#1A241B] aspect-[16/9] md:aspect-[21/9] shadow-xl">
        <img
          :src="assetUrl('/images/about-sunset-home.jpg')"
          alt="Custom Texas Hill Country limestone and glass residence at sunset"
          class="w-full h-full object-cover"
        />
        <div
          class="absolute inset-0 bg-gradient-to-t from-[#141A15]/75 via-transparent to-transparent"
        />
        <div class="absolute bottom-6 left-6 md:bottom-10 md:left-10 text-white">
          <span class="inline-flex items-center gap-2 text-xs font-semibold uppercase tracking-[0.18em] text-[#4ADE80]">
            <ShieldCheck class="w-4 h-4" />
            TEXAS HILL COUNTRY CRAFTSMEN
          </span>
          <p class="text-2xl md:text-4xl font-medium tracking-tight mt-2">
            Crafted to outlast the home itself.
          </p>
        </div>
      </div>

      <!-- Story Split Section -->
      <div class="about-story-grid grid grid-cols-1 lg:grid-cols-12 gap-10 lg:gap-16 items-start px-2">
        <div class="lg:col-span-4 space-y-4">
          <span class="block text-xs font-semibold uppercase tracking-[0.18em] text-[#22C55E]">
            ABOUT COPELY
          </span>
          <div class="grid grid-cols-2 gap-4 pt-2">
            <div
              v-for="stat in stats"
              :key="stat.label"
              class="about-stat-anim rounded-2xl bg-white border border-[#E2E6E0] p-4"
            >
              <span class="block text-2xl font-semibold text-[#141A15]">{{ stat.value }}</span>
              <span class="block text-xs text-[#5A655C] mt-1">{{ stat.label }}</span>
            </div>
          </div>
        </div>

        <div class="lg:col-span-8 space-y-6">
          <h1 class="text-3xl sm:text-4xl md:text-5xl font-medium tracking-[-0.03em] leading-[1.12] text-[#141A15]">
            Copely started in 2009 with two tile saws, one truck, and a simple promise:
            <span class="font-serif italic font-normal text-[#1A241B]">
              do the job right and put it in writing.
            </span>
          </h1>

          <div class="space-y-4 text-base md:text-lg text-[#5A655C] leading-relaxed">
            <p>
              Too many Central Texas homeowners were hiring remodelers who buried shortcuts under
              pretty tile—skipping crack-isolation membranes over shifting clay slabs, grouting
              change-of-plane corners rigid until they split, or tiling shower walls directly over
              unsealed drywall.
            </p>
            <p>
              Founder Sam Calloway and Operations Director Erin Vasquez built Copely to treat tile as
              an architectural system rather than a cosmetic finish. Every floor we level, every
              curbless shower pan we flood-test for 24 hours, and every patio paver terrace we set
              is documented for your records.
            </p>
          </div>

          <div class="about-pillars-grid pt-4 grid grid-cols-1 md:grid-cols-3 gap-5">
            <div
              v-for="pillar in pillars"
              :key="pillar.title"
              class="about-pillar-anim rounded-2xl bg-[#E8EFE6] p-6 space-y-2.5"
            >
              <span
                class="w-7 h-7 rounded-full bg-[#1A241B] text-[#4ADE80] inline-flex items-center justify-center"
              >
                <Check class="w-4 h-4" />
              </span>
              <h2 class="text-lg font-medium text-[#141A15]">{{ pillar.title }}</h2>
              <p class="text-sm text-[#5A655C] leading-relaxed">{{ pillar.body }}</p>
            </div>
          </div>

          <div class="pt-4">
            <button
              type="button"
              class="rounded-full bg-[#1A241B] text-white px-7 py-3.5 text-sm font-medium hover:bg-[#273528] transition shadow-md cursor-pointer"
              @click="goToEstimate"
            >
              Schedule your free estimate
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
