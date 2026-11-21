<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { ref } from 'vue';
import CopelyLogo from '@/components/ui/CopelyLogo.vue';
import { FOOTER_POLAROIDS } from '@/data/copelyData';
import { useGsapScope } from '@/composables/useGsapScope';

const route = useRoute();
const router = useRouter();
const footerRef = ref<HTMLElement | null>(null);
const polaroidRowRef = ref<HTMLElement | null>(null);
const watermarkRef = ref<HTMLElement | null>(null);

useGsapScope(footerRef, ({ gsap }) => {
  const fanEls = polaroidRowRef.value?.querySelectorAll<HTMLElement>('.polaroid-fan-item') ?? [];
  const targetRotations = [-8, -4, 0, 4, 8];
  const targetY = [12, 4, -8, 4, 12];
  const startX = [120, 60, 0, -60, -120];

  if (fanEls.length > 0 && polaroidRowRef.value) {
    fanEls.forEach((el, idx) => {
      gsap.fromTo(
        el,
        {
          x: startX[idx] ?? 0,
          y: 28,
          rotation: 0,
          scale: 0.9,
          opacity: 0,
        },
        {
          x: 0,
          y: targetY[idx] ?? 0,
          rotation: targetRotations[idx] ?? 0,
          scale: 1,
          opacity: 1,
          duration: 0.75,
          delay: Math.abs(idx - 2) * 0.06,
          ease: 'power3.out',
          scrollTrigger: {
            trigger: polaroidRowRef.value,
            start: 'top 82%',
            once: true,
          },
        },
      );
    });
  }

  if (watermarkRef.value) {
    gsap.fromTo(
      watermarkRef.value,
      { y: 42 },
      {
        y: 0,
        ease: 'none',
        scrollTrigger: {
          trigger: '.footer-dark-base',
          start: 'top 95%',
          end: 'bottom bottom',
          scrub: true,
        },
      },
    );
  }
});

function goToEstimate(): void {
  if (route.path === '/') {
    const el = document.getElementById('schedule');
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' });
      return;
    }
  }
  void router.push({ path: '/', hash: '#schedule' });
}
</script>

<template>
  <footer ref="footerRef" class="mt-20">
    <!-- Part A: Pre-Footer Polaroid Photo Fan & CTA Card -->
    <section
      class="mx-3 mb-0 rounded-t-3xl bg-gradient-to-b from-[#E8EFE6] to-[#DDE8DA] pt-20 pb-16 px-4 text-center overflow-hidden"
    >
      <!-- 5 Overlapping Tilted Polaroid Cards -->
      <div
        ref="polaroidRowRef"
        class="flex justify-center items-center -space-x-6 md:-space-x-5 mb-12 pt-4"
      >
        <div
          v-for="(card, idx) in FOOTER_POLAROIDS"
          :key="card.id"
          class="polaroid-fan-item shrink-0 hover:z-30"
          :class="idx === 2 ? 'z-20' : idx === 1 || idx === 3 ? 'z-10' : 'z-0'"
        >
          <div
            class="w-32 sm:w-40 md:w-48 aspect-[4/5] rounded-2xl border-[6px] border-white bg-white shadow-xl overflow-hidden transition-transform duration-300 hover:scale-105"
          >
            <img
              :src="card.imageUrl"
              :alt="card.alt"
              loading="lazy"
              class="w-full h-full object-cover"
            />
          </div>
        </div>
      </div>

      <!-- Centered Typography Below Fan -->
      <div class="max-w-xl mx-auto">
        <h2 class="text-3xl md:text-4xl font-medium text-[#141A15] tracking-tight leading-tight">
          Done in days.
        </h2>
        <p class="font-serif italic text-3xl md:text-4xl text-[#141A15] mt-1 leading-tight">
          Built for decades.
        </p>
        <p class="text-sm text-[#5A655C] max-w-md mx-auto mt-3 leading-relaxed">
          A licensed local crew, honest answers, and a 25-year warranty in writing — anywhere in
          Texas Hill Country.
        </p>

        <button
          type="button"
          class="mt-6 bg-[#1A241B] text-white px-6 py-3 rounded-full text-sm font-medium inline-flex items-center gap-2.5 hover:bg-[#273528] transition-colors shadow-lg cursor-pointer"
          @click="goToEstimate"
        >
          <span class="w-2 h-2 rounded-full bg-[#4ADE80]" aria-hidden="true" />
          <span>Get a free estimate</span>
        </button>
      </div>
    </section>

    <!-- Part B: Infinite Scrolling Green Marquee Banner -->
    <div
      role="button"
      tabindex="0"
      aria-label="Get a free estimate"
      class="marquee-strip w-full bg-[#34D361] text-[#141A15] py-4 overflow-hidden flex whitespace-nowrap border-y border-[#1A241B]/10 cursor-pointer"
      @click="goToEstimate"
      @keydown.enter="goToEstimate"
      @keydown.space.prevent="goToEstimate"
    >
      <div class="marquee-track items-center">
        <div
          v-for="copy in 2"
          :key="copy"
          class="flex items-center shrink-0"
          :aria-hidden="copy === 2"
        >
          <div
            v-for="idx in 6"
            :key="`${copy}-${idx}`"
            class="inline-flex items-center gap-4 mx-6"
          >
            <svg
              width="28"
              height="28"
              viewBox="0 0 28 28"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              aria-hidden="true"
            >
              <path
                d="M3.5 19.5L14 9L24.5 19.5L21.5 22.5L14 15L6.5 22.5L3.5 19.5Z"
                fill="#1A241B"
              />
              <path
                d="M7.5 11.5L14 5L20.5 11.5L17.8 14.2L14 10.4L10.2 14.2L7.5 11.5Z"
                fill="#141A15"
              />
            </svg>
            <span class="text-3xl md:text-4xl font-medium tracking-tight text-[#141A1B]">
              Get a free estimate
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Part C: Deep Forest Main Footer with Giant Cropped Watermark -->
    <div class="footer-dark-base bg-[#1A241B] text-white pt-16 pb-0 px-8 md:px-16 relative overflow-hidden">
      <div class="max-w-[1320px] mx-auto grid grid-cols-1 md:grid-cols-12 gap-12 pb-32 md:pb-36 relative z-10">
        <!-- Left Column: Brand, Tagline, Socials, Legal -->
        <div class="md:col-span-6">
          <RouterLink to="/" class="inline-flex items-center">
            <CopelyLogo wordmark-class="text-white" :icon-size="24" />
          </RouterLink>

          <div class="mt-6">
            <p class="text-4xl md:text-5xl font-normal tracking-tight text-white leading-[1.08]">
              Done in days.
            </p>
            <p class="text-4xl md:text-5xl font-serif italic font-normal text-white/90 mt-1 leading-[1.08]">
              Built for decades.
            </p>
          </div>

          <!-- Social Buttons -->
          <div class="mt-8 flex items-center gap-3">
            <a
              href="https://instagram.com"
              target="_blank"
              rel="noreferrer"
              aria-label="Copely on Instagram"
              class="w-10 h-10 rounded-full bg-white/10 hover:bg-white/20 flex items-center justify-center transition text-white"
            >
              <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="2" y="2" width="20" height="20" rx="5" />
                <circle cx="12" cy="12" r="4" />
                <circle cx="17.5" cy="6.5" r="1" fill="currentColor" stroke="none" />
              </svg>
            </a>
            <a
              href="https://facebook.com"
              target="_blank"
              rel="noreferrer"
              aria-label="Copely on Facebook"
              class="w-10 h-10 rounded-full bg-white/10 hover:bg-white/20 flex items-center justify-center transition text-white"
            >
              <svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M14 13.5h2.5l1-4H14v-2c0-1.03 0-2 2-2h1.5V2.14c-.326-.043-1.557-.14-2.857-.14C11.928 2 10 3.657 10 6.7v2.8H7v4h3V22h4v-8.5Z"
                />
              </svg>
            </a>
          </div>

          <!-- Legal & Copyright Stack -->
          <div class="mt-14 space-y-2 text-xs text-white/50">
            <RouterLink to="/privacy" class="block hover:text-white/80 transition">
              Privacy policy
            </RouterLink>
            <RouterLink to="/terms" class="block hover:text-white/80 transition">
              Terms of service
            </RouterLink>
            <p class="pt-2 text-white/40">
              © Copely · Licensed &amp; Insured · Serving Texas Hill Country
            </p>
          </div>
        </div>

        <!-- Middle Column: Menu -->
        <div class="md:col-span-3">
          <h3 class="text-xs uppercase tracking-wider text-white/50 mb-4 font-medium">Menu</h3>
          <ul class="space-y-3 text-sm text-white/80">
            <li>
              <RouterLink to="/services" class="hover:text-white transition">Services</RouterLink>
            </li>
            <li>
              <RouterLink to="/about" class="hover:text-white transition">About</RouterLink>
            </li>
            <li>
              <RouterLink to="/blog" class="hover:text-white transition">Blog</RouterLink>
            </li>
            <li>
              <RouterLink to="/faq" class="hover:text-white transition">FAQ</RouterLink>
            </li>
            <li>
              <RouterLink to="/contact" class="hover:text-white transition">Contact</RouterLink>
            </li>
          </ul>
        </div>

        <!-- Right Column: Services -->
        <div class="md:col-span-3">
          <h3 class="text-xs uppercase tracking-wider text-white/50 mb-4 font-medium">Services</h3>
          <ul class="space-y-3 text-sm text-white/80">
            <li>
              <RouterLink to="/services/porcelain-ceramic" class="hover:text-white transition">
                Porcelain &amp; Tile
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/services/wet-rooms" class="hover:text-white transition">
                Wet Rooms
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/services/kitchen-backsplashes" class="hover:text-white transition">
                Backsplashes
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/services/natural-stone" class="hover:text-white transition">
                Natural Stone
              </RouterLink>
            </li>
            <li>
              <RouterLink to="/services/outdoor-patios" class="hover:text-white transition">
                Outdoor Patios
              </RouterLink>
            </li>
          </ul>
        </div>
      </div>

      <!-- Giant Bottom-Right Cropped Brand Watermark -->
      <div
        ref="watermarkRef"
        aria-hidden="true"
        class="pointer-events-none select-none absolute -bottom-6 right-4 md:right-8 leading-[0.75]"
      >
        <span class="block text-[22vw] md:text-[19vw] font-semibold tracking-tighter text-[#34D361]">
          Copely
        </span>
      </div>
    </div>
  </footer>
</template>
