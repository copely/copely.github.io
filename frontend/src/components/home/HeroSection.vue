<script setup lang="ts">
import { ref } from 'vue';
import { Calendar, ChevronDown, Grid2x2, Star } from 'lucide-vue-next';
import { SERVICE_OPTIONS, TIMELINE_OPTIONS } from '@/data/copelyData';
import { useEstimateDraft } from '@/composables/useEstimateDraft';
import { useGsapScope } from '@/composables/useGsapScope';
import { assetUrl } from '@/utils/assetUrl';
import type { ProjectTimeline, TileServiceCategory } from '@/types/copely';

const { estimateDraft, setDraftSelection } = useEstimateDraft();

const selectedService = ref<TileServiceCategory>(estimateDraft.service);
const selectedTimeline = ref<ProjectTimeline>(estimateDraft.timeline);
const heroImageSrc = assetUrl('/images/hero-craftsman.jpg');
const heroVideoSrc = assetUrl('/videos/hero-loop.mp4');
const videoReady = ref<boolean>(false);
const videoUnavailable = ref<boolean>(false);

const sectionRef = ref<HTMLElement | null>(null);
const backdropRef = ref<HTMLElement | null>(null);
const eyebrowRuleRef = ref<HTMLElement | null>(null);
const cardLaserRef = ref<HTMLElement | null>(null);

useGsapScope(sectionRef, ({ gsap }) => {
  const tl = gsap.timeline({ defaults: { ease: 'power3.out' } });

  if (eyebrowRuleRef.value) {
    tl.from(eyebrowRuleRef.value, { scaleX: 0, transformOrigin: 'left center', duration: 0.45 }, 0.05);
  }
  tl.from('.hero-eyebrow-text', { x: -12, opacity: 0, duration: 0.45 }, 0.12)
    .from('.hero-title-line', { y: 28, opacity: 0, duration: 0.65, stagger: 0.1 }, 0.18)
    .from('.hero-fade-up', { y: 20, opacity: 0, duration: 0.55, stagger: 0.1 }, 0.42)
    .from('.hero-estimate-card', { y: 32, scale: 0.98, opacity: 0, duration: 0.7 }, 0.34);

  if (cardLaserRef.value) {
    tl.fromTo(
      cardLaserRef.value,
      { xPercent: -100, opacity: 0 },
      { xPercent: 120, opacity: 1, duration: 1.0, ease: 'power2.inOut' },
      0.85,
    ).to(cardLaserRef.value, { opacity: 0, duration: 0.25 });
  }

  if (backdropRef.value && sectionRef.value) {
    gsap.to(backdropRef.value, {
      yPercent: 8,
      ease: 'none',
      scrollTrigger: {
        trigger: sectionRef.value,
        start: 'top top',
        end: 'bottom top',
        scrub: true,
      },
    });
  }
});

function handleGetEstimate(): void {
  setDraftSelection(selectedService.value, selectedTimeline.value);
  const target = document.getElementById('schedule');
  if (target) {
    target.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }
}
</script>

<template>
  <section ref="sectionRef" class="px-3 md:px-4 pt-3 md:pt-4">
    <div
      class="relative min-h-[680px] lg:min-h-[760px] rounded-3xl overflow-hidden bg-[#1A241B] text-white flex items-end"
    >
      <!-- Full-bleed video wrapper with crisp generated hero image + Ken Burns fallback -->
      <div ref="backdropRef" class="absolute inset-0 overflow-hidden will-change-transform">
        <img
          :src="heroImageSrc"
          alt="Copely tile craftsmen setting large-format porcelain slabs with a laser level"
          class="w-full h-full object-cover object-center animate-ken-burns"
        />
        <video
          v-if="!videoUnavailable"
          autoplay
          loop
          muted
          playsinline
          :poster="heroImageSrc"
          :src="heroVideoSrc"
          aria-hidden="true"
          class="absolute inset-0 w-full h-full object-cover transition-opacity duration-700"
          :class="videoReady ? 'opacity-100' : 'opacity-0 pointer-events-none'"
          @canplay="videoReady = true"
          @error="videoUnavailable = true"
        />

        <!-- Atmospheric gradient overlay matching Apex dark forest aesthetic -->
        <div
          class="absolute inset-0 bg-gradient-to-t from-[#141A15]/95 via-[#1A241B]/55 to-[#141A15]/45"
        />
        <div
          class="absolute inset-0 bg-gradient-to-r from-[#141A15]/85 via-[#141A15]/35 to-transparent"
        />
      </div>

      <!-- Foreground Split Layout -->
      <div
        class="relative z-10 w-full max-w-[1320px] mx-auto px-6 md:px-12 pt-32 pb-12 md:pb-16 grid grid-cols-1 lg:grid-cols-12 gap-10 lg:gap-12 items-end"
      >
        <!-- Left Content -->
        <div class="lg:col-span-7 space-y-6">
          <div class="flex items-center gap-3">
            <span ref="eyebrowRuleRef" class="w-8 h-[2px] bg-[#4ADE80]" aria-hidden="true" />
            <span class="hero-eyebrow-text text-xs font-semibold uppercase tracking-[0.18em] text-[#4ADE80]">
              SERVING THE TEXAS HILL COUNTRY
            </span>
          </div>

          <h1
            class="text-4xl sm:text-5xl md:text-6xl lg:text-[64px] font-medium tracking-[-0.03em] leading-[1.04] text-white"
          >
            <span class="hero-title-line block">Copely Tiling.</span>
            <span class="hero-title-line block">Done in Days.</span>
            <span class="hero-title-line block font-serif italic font-normal text-white/95">
              Built for Decades.
            </span>
          </h1>

          <p class="hero-fade-up text-base md:text-lg text-white/80 max-w-xl leading-relaxed">
            Installations and restorations backed by a 25-year warranty. Fast, careful work from a
            local crew. Get a free estimate anywhere in Texas Hill Country.
          </p>

          <!-- Social Proof Badge -->
          <div class="hero-fade-up inline-flex items-center gap-2.5 pt-1">
            <span
              class="inline-flex items-center gap-1.5 rounded-full bg-white/10 backdrop-blur-md border border-white/15 px-3.5 py-1.5 text-sm font-medium text-white"
            >
              <Star class="w-4 h-4 fill-[#4ADE80] text-[#4ADE80]" />
              <span>4.9</span>
            </span>
            <span class="text-sm text-white/75">from 1,200+ reviews</span>
          </div>
        </div>

        <!-- Right Floating Glassmorphic Estimate Card -->
        <div class="lg:col-span-5 lg:pl-6">
          <div
            class="hero-estimate-card relative overflow-hidden rounded-3xl bg-[#1A241B]/95 backdrop-blur-xl border border-white/15 p-6 md:p-7 shadow-2xl"
          >
            <!-- Subtle Emerald Laser Sweep -->
            <div
              ref="cardLaserRef"
              class="pointer-events-none absolute top-0 left-0 right-0 h-[2px] bg-gradient-to-r from-transparent via-[#4ADE80] to-transparent opacity-0"
              aria-hidden="true"
            />

            <h2 class="text-xl font-medium text-white tracking-tight">
              Get Your Free Tile Estimate
            </h2>

            <div class="mt-5 space-y-3.5">
              <!-- Dropdown 1: Service -->
              <div>
                <label for="hero-service" class="block text-xs text-white/60 mb-1.5 font-medium">
                  Service
                </label>
                <div class="relative">
                  <Grid2x2
                    class="w-4 h-4 text-[#4ADE80] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
                  />
                  <select
                    id="hero-service"
                    v-model="selectedService"
                    class="w-full appearance-none rounded-xl bg-[#253227] border border-white/15 pl-10 pr-10 py-3 text-sm font-medium text-white focus:outline-none focus:border-[#4ADE80] transition-colors cursor-pointer"
                  >
                    <option
                      v-for="opt in SERVICE_OPTIONS"
                      :key="opt.value"
                      :value="opt.value"
                      class="bg-[#1A241B] text-white"
                    >
                      {{ opt.label }}
                    </option>
                  </select>
                  <ChevronDown
                    class="w-4 h-4 text-white/60 absolute right-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
                  />
                </div>
              </div>

              <!-- Dropdown 2: Timeline -->
              <div>
                <label for="hero-timeline" class="block text-xs text-white/60 mb-1.5 font-medium">
                  Timeline
                </label>
                <div class="relative">
                  <Calendar
                    class="w-4 h-4 text-[#4ADE80] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
                  />
                  <select
                    id="hero-timeline"
                    v-model="selectedTimeline"
                    class="w-full appearance-none rounded-xl bg-[#253227] border border-white/15 pl-10 pr-10 py-3 text-sm font-medium text-white focus:outline-none focus:border-[#4ADE80] transition-colors cursor-pointer"
                  >
                    <option
                      v-for="opt in TIMELINE_OPTIONS"
                      :key="opt.value"
                      :value="opt.value"
                      class="bg-[#1A241B] text-white"
                    >
                      {{ opt.label }}
                    </option>
                  </select>
                  <ChevronDown
                    class="w-4 h-4 text-white/60 absolute right-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
                  />
                </div>
              </div>
            </div>

            <!-- 2-Column Info Row -->
            <div class="grid grid-cols-2 gap-3 mt-5 pt-4 border-t border-white/10">
              <div class="rounded-xl bg-white/5 px-3.5 py-2.5">
                <span class="block text-[11px] uppercase tracking-wider text-white/50">
                  Response
                </span>
                <span class="block text-sm font-semibold text-white mt-0.5">Within 24 hrs</span>
              </div>
              <div class="rounded-xl bg-white/5 px-3.5 py-2.5">
                <span class="block text-[11px] uppercase tracking-wider text-white/50">
                  Warranty
                </span>
                <span class="block text-sm font-semibold text-white mt-0.5">Up to 25 yrs</span>
              </div>
            </div>

            <!-- Footer Meta Row -->
            <div class="mt-4 flex items-center justify-between text-xs text-white/60 px-1">
              <span>Free estimate</span>
              <span class="text-white/25" aria-hidden="true">|</span>
              <span>Licensed &amp; insured</span>
            </div>

            <!-- Primary CTA Button -->
            <button
              type="button"
              class="mt-4 w-full rounded-full bg-white text-[#141A15] py-3.5 px-6 text-sm font-semibold hover:bg-[#E8EFE6] transition-colors shadow-lg cursor-pointer"
              @click="handleGetEstimate"
            >
              Get My Free Estimate
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
