<script setup lang="ts">
import { ref } from 'vue';
import gsap from 'gsap';
import { ArrowLeft, ArrowRight, ArrowUpRight } from 'lucide-vue-next';
import { useGsapScope } from '@/composables/useGsapScope';
import type { ServiceItem } from '@/types/copely';

interface Props {
  services: ServiceItem[];
}

defineProps<Props>();

const sectionRef = ref<HTMLElement | null>(null);
const trackRef = ref<HTMLElement | null>(null);

useGsapScope(sectionRef, ({ gsap: g }) => {
  g.from('.carousel-header-anim', {
    y: 22,
    opacity: 0,
    duration: 0.6,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: sectionRef.value,
      start: 'top 82%',
      once: true,
    },
  });

  g.from('.service-card-anim', {
    y: 34,
    opacity: 0,
    duration: 0.65,
    stagger: 0.08,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: trackRef.value,
      start: 'top 82%',
      once: true,
    },
  });
});

function scrollCards(direction: 'prev' | 'next'): void {
  const el = trackRef.value;
  if (!el) {
    return;
  }
  const offset = direction === 'next' ? 360 : -360;
  const targetLeft = Math.max(0, Math.min(el.scrollWidth - el.clientWidth, el.scrollLeft + offset));
  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    el.scrollLeft = targetLeft;
    return;
  }
  gsap.to(el, {
    scrollLeft: targetLeft,
    duration: 0.55,
    ease: 'power3.out',
  });
}
</script>

<template>
  <section ref="sectionRef" class="px-3 md:px-4">
    <div class="rounded-3xl bg-[#E8EFE6] py-16 md:py-24 px-6 md:px-12 overflow-hidden">
      <div class="max-w-[1320px] mx-auto">
        <!-- Header Row -->
        <div
          class="carousel-header-anim flex flex-col md:flex-row md:items-end justify-between gap-6 mb-10 md:mb-12"
        >
          <h2 class="text-3xl sm:text-4xl md:text-5xl font-medium tracking-[-0.03em] leading-[1.08]">
            <span class="block text-[#141A15]">One local crew.</span>
            <span class="block text-[#22C55E]">The whole interior &amp; patio.</span>
          </h2>

          <div class="flex items-center gap-4 md:gap-5">
            <p class="text-sm text-[#5A655C] max-w-[230px] leading-snug">
              Licensed &amp; insured tiling across Texas Hill Country
            </p>
            <div class="flex items-center gap-2 shrink-0">
              <button
                type="button"
                aria-label="Previous service"
                class="w-11 h-11 rounded-full bg-white text-[#141A15] border border-[#E2E6E0] flex items-center justify-center hover:bg-[#1A241B] hover:text-white transition-colors shadow-sm cursor-pointer"
                @click="scrollCards('prev')"
              >
                <ArrowLeft class="w-4 h-4" />
              </button>
              <button
                type="button"
                aria-label="Next service"
                class="w-11 h-11 rounded-full bg-white text-[#141A15] border border-[#E2E6E0] flex items-center justify-center hover:bg-[#1A241B] hover:text-white transition-colors shadow-sm cursor-pointer"
                @click="scrollCards('next')"
              >
                <ArrowRight class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

        <!-- Interactive Card Track (5 Tall Rounded Cards, aspect-[3/4]) -->
        <div
          ref="trackRef"
          class="flex gap-5 overflow-x-auto snap-x snap-mandatory no-scrollbar pb-2"
        >
          <RouterLink
            v-for="service in services"
            :key="service.id"
            :to="`/services/${service.slug}`"
            class="service-card-anim group relative snap-start shrink-0 w-[275px] sm:w-[300px] lg:w-[calc((100%-60px)/4)] min-w-[270px] aspect-[3/4] rounded-3xl overflow-hidden bg-[#1A241B]"
          >
            <img
              :src="service.imageUrl"
              :alt="service.title"
              loading="lazy"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div
              class="absolute inset-0 bg-gradient-to-t from-[#141A15]/90 via-[#141A15]/20 to-transparent"
            />

            <div class="absolute inset-x-0 bottom-0 p-6 text-white">
              <div class="flex items-center justify-between gap-2">
                <h3 class="text-2xl font-medium tracking-tight text-white">
                  {{ service.title }}
                </h3>
                <ArrowUpRight
                  class="w-5 h-5 text-[#4ADE80] opacity-0 -translate-x-1 translate-y-1 transition-all duration-300 group-hover:opacity-100 group-hover:translate-x-0 group-hover:translate-y-0"
                />
              </div>
              <p
                class="text-sm text-white/80 mt-1.5 max-h-0 opacity-0 overflow-hidden transition-all duration-300 group-hover:max-h-16 group-hover:opacity-100"
              >
                {{ service.shortSubtitle }}
              </p>
            </div>
          </RouterLink>
        </div>
      </div>
    </div>
  </section>
</template>
