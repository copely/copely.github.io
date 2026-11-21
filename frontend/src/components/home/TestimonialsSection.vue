<script setup lang="ts">
import { computed, ref } from 'vue';
import { ArrowUpRight } from 'lucide-vue-next';
import { useGsapScope } from '@/composables/useGsapScope';
import type { Testimonial } from '@/types/copely';

interface Props {
  testimonials: Testimonial[];
}

const props = defineProps<Props>();

const activeIndex = ref<number>(0);
const sectionRef = ref<HTMLElement | null>(null);

const activeTestimonial = computed<Testimonial | undefined>(() => {
  return props.testimonials[activeIndex.value] ?? props.testimonials[0];
});

useGsapScope(sectionRef, ({ gsap }) => {
  gsap.from('.story-avatar-btn', {
    y: 16,
    scale: 0.9,
    opacity: 0,
    duration: 0.45,
    stagger: 0.05,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: sectionRef.value,
      start: 'top 82%',
      once: true,
    },
  });

  gsap.from('.story-case-card', {
    y: 28,
    opacity: 0,
    duration: 0.65,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.story-case-card',
      start: 'top 84%',
      once: true,
    },
  });
});
</script>

<template>
  <section ref="sectionRef" class="py-20 md:py-24 px-6 md:px-10 bg-[#F7F7F3]">
    <div v-if="activeTestimonial" class="max-w-[1120px] mx-auto">
      <!-- Avatar Selector Bar -->
      <div class="flex items-center justify-center gap-3 sm:gap-4 mb-10 flex-wrap">
        <button
          v-for="(item, idx) in testimonials"
          :key="item.id"
          type="button"
          :aria-label="`View story from ${item.author}`"
          :aria-pressed="idx === activeIndex"
          class="story-avatar-btn relative w-14 h-14 sm:w-16 sm:h-16 rounded-2xl overflow-hidden transition-all duration-200 cursor-pointer focus:outline-none"
          :class="
            idx === activeIndex
              ? 'ring-2 ring-[#34D361] ring-offset-2 ring-offset-[#F7F7F3] scale-105 shadow-md'
              : 'grayscale opacity-70 hover:opacity-100 hover:grayscale-0'
          "
          @click="activeIndex = idx"
        >
          <img
            :src="item.avatarUrl"
            :alt="item.author"
            loading="lazy"
            class="w-full h-full object-cover"
          />
        </button>
      </div>

      <!-- Quote Display -->
      <div class="text-center max-w-3xl mx-auto min-h-[150px] flex flex-col justify-center">
        <Transition
          mode="out-in"
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-0 translate-y-1.5"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100 translate-y-0"
          leave-to-class="opacity-0 -translate-y-1.5"
        >
          <div :key="activeTestimonial.id" class="space-y-5">
            <blockquote
              class="text-2xl sm:text-3xl md:text-[32px] font-medium tracking-[-0.02em] leading-[1.28] text-[#141A15]"
            >
              “{{ activeTestimonial.quote }}”
            </blockquote>
            <p class="text-sm text-[#5A655C]">
              <span class="font-semibold text-[#141A15]">{{ activeTestimonial.author }}</span>
              <span class="mx-1.5" aria-hidden="true">·</span>
              <span class="italic">
                {{ activeTestimonial.projectType }} · {{ activeTestimonial.location }}
              </span>
            </p>
          </div>
        </Transition>
      </div>

      <!-- Featured Case Study Visual Card -->
      <div class="story-case-card mt-10 relative rounded-3xl overflow-hidden bg-[#1A241B] aspect-[16/9] md:aspect-[21/9] shadow-xl">
        <Transition
          mode="out-in"
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-85"
          enter-to-class="opacity-100"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-85"
        >
          <img
            :key="activeTestimonial.heroImageUrl"
            :src="activeTestimonial.heroImageUrl"
            :alt="`${activeTestimonial.projectType} in ${activeTestimonial.location}`"
            class="w-full h-full object-cover"
          />
        </Transition>

        <div
          class="absolute inset-0 bg-gradient-to-t from-[#141A15]/90 via-[#141A15]/20 to-transparent"
        />

        <!-- Bottom Overlay Metrics Bar -->
        <div
          class="absolute inset-x-4 bottom-4 md:inset-x-6 md:bottom-6 rounded-2xl bg-[#1A241B]/90 backdrop-blur-md border border-white/15 px-5 py-4 flex flex-wrap items-center justify-between gap-4 text-white"
        >
          <div class="flex flex-wrap items-center gap-6 md:gap-10">
            <div
              v-for="metric in activeTestimonial.metrics"
              :key="metric.label"
              class="flex flex-col"
            >
              <span class="text-base md:text-lg font-semibold text-white">
                {{ metric.value }}
              </span>
              <span class="text-[11px] uppercase tracking-wider text-white/60">
                {{ metric.label }}
              </span>
            </div>
          </div>

          <RouterLink
            to="/services"
            class="inline-flex items-center gap-1.5 rounded-full bg-[#34D361] text-[#141A15] px-4 py-2 text-xs md:text-sm font-semibold hover:bg-[#4ADE80] transition-colors"
          >
            <span>See our services</span>
            <ArrowUpRight class="w-4 h-4" />
          </RouterLink>
        </div>
      </div>
    </div>
  </section>
</template>
