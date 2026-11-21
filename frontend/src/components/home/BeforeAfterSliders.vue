<script setup lang="ts">
import { computed, ref } from 'vue';
import ImageCompareSlider from '@/components/ui/ImageCompareSlider.vue';
import { useGsapScope } from '@/composables/useGsapScope';
import type { BeforeAfterCase } from '@/types/copely';

interface Props {
  cases: BeforeAfterCase[];
}

const props = defineProps<Props>();
const sectionRef = ref<HTMLElement | null>(null);

const heroCase = computed<BeforeAfterCase | undefined>(() => {
  return props.cases.find((item) => item.featured) ?? props.cases[0];
});

const secondaryCases = computed<BeforeAfterCase[]>(() => {
  const heroId = heroCase.value?.id;
  return props.cases.filter((item) => item.id !== heroId).slice(0, 3);
});

useGsapScope(sectionRef, ({ gsap }) => {
  gsap.from('.ba-header-anim', {
    y: 24,
    opacity: 0,
    duration: 0.6,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: sectionRef.value,
      start: 'top 80%',
      once: true,
    },
  });

  gsap.from('.ba-hero-wrap', {
    y: 30,
    opacity: 0,
    duration: 0.7,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.ba-hero-wrap',
      start: 'top 82%',
      once: true,
    },
  });

  gsap.from('.ba-secondary-card', {
    y: 28,
    opacity: 0,
    duration: 0.6,
    stagger: 0.1,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.ba-secondary-grid',
      start: 'top 82%',
      once: true,
    },
  });
});
</script>

<template>
  <section ref="sectionRef" class="py-20 md:py-28 px-6 md:px-10 bg-[#F7F7F3]">
    <div class="max-w-[1320px] mx-auto space-y-10 md:space-y-12">
      <!-- Header -->
      <div class="ba-header-anim grid grid-cols-1 lg:grid-cols-12 gap-6 items-end">
        <div class="lg:col-span-7">
          <h2 class="text-4xl md:text-5xl font-medium tracking-[-0.03em] leading-[1.08] text-[#141A15]">
            Before we showed up.<br />
            <span class="font-serif italic font-normal text-[#1A241B]">After we left.</span>
          </h2>
        </div>
        <div class="lg:col-span-5">
          <p class="text-base text-[#5A655C] leading-relaxed lg:max-w-md lg:ml-auto">
            Cracked grout, uneven tiles, and slow shower leaks — gone. Slide across to see what a
            few days with our crew changes.
          </p>
        </div>
      </div>

      <!-- Large Full-Width Hero Comparison Slider -->
      <div v-if="heroCase" class="ba-hero-wrap space-y-3">
        <ImageCompareSlider
          :before-image="heroCase.beforeImage"
          :after-image="heroCase.afterImage"
          :before-alt="`Before: ${heroCase.title}`"
          :after-alt="`After: ${heroCase.title}`"
          variant="hero"
          :initial-position="48"
          :auto-sweep-on-visible="true"
        />
        <div class="flex flex-wrap items-center justify-between gap-2 px-1 text-sm">
          <span class="font-medium text-[#141A15]">{{ heroCase.title }}</span>
          <span class="text-[#5A655C]">{{ heroCase.subtitle }} · {{ heroCase.location }}</span>
        </div>
      </div>

      <!-- 3 Smaller 3-Column Comparison Sliders -->
      <div class="ba-secondary-grid grid grid-cols-1 md:grid-cols-3 gap-6">
        <div v-for="item in secondaryCases" :key="item.id" class="ba-secondary-card space-y-3">
          <ImageCompareSlider
            :before-image="item.beforeImage"
            :after-image="item.afterImage"
            :before-alt="`Before: ${item.title}`"
            :after-alt="`After: ${item.title}`"
            variant="compact"
            :initial-position="50"
          />
          <div class="px-1">
            <h3 class="text-base font-medium text-[#141A15]">{{ item.title }}</h3>
            <p class="text-xs text-[#5A655C] mt-0.5">{{ item.subtitle }} · {{ item.location }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
