<script setup lang="ts">
import { ref } from 'vue';
import { PROCESS_STEPS } from '@/data/copelyData';
import { useGsapScope } from '@/composables/useGsapScope';

const activeIndex = ref<number>(0);
const sectionRef = ref<HTMLElement | null>(null);
const timelineColRef = ref<HTMLElement | null>(null);
const railFillRef = ref<HTMLElement | null>(null);

useGsapScope(sectionRef, ({ gsap, ScrollTrigger }) => {
  if (railFillRef.value && timelineColRef.value) {
    gsap.fromTo(
      railFillRef.value,
      { scaleY: 0 },
      {
        scaleY: 1,
        ease: 'none',
        transformOrigin: 'top center',
        scrollTrigger: {
          trigger: timelineColRef.value,
          start: 'top 62%',
          end: 'bottom 58%',
          scrub: 0.35,
        },
      },
    );
  }

  const cards = sectionRef.value?.querySelectorAll<HTMLElement>('.process-step-card') ?? [];
  cards.forEach((card, idx) => {
    gsap.from(card, {
      y: 26,
      opacity: 0,
      duration: 0.55,
      ease: 'power3.out',
      scrollTrigger: {
        trigger: card,
        start: 'top 86%',
        once: true,
      },
    });

    ScrollTrigger.create({
      trigger: card,
      start: 'top 56%',
      end: 'bottom 44%',
      onEnter: () => {
        activeIndex.value = idx;
      },
      onEnterBack: () => {
        activeIndex.value = idx;
      },
    });
  });
});

function scrollToSchedule(): void {
  const el = document.getElementById('schedule');
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }
}
</script>

<template>
  <section ref="sectionRef" class="px-3 md:px-4">
    <div class="rounded-3xl bg-[#1A241B] text-white py-20 md:py-28 px-6 md:px-12">
      <div class="max-w-[1320px] mx-auto grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16 items-start">
        <!-- Left Sticky Column -->
        <div class="lg:col-span-5 lg:sticky lg:top-28 space-y-6">
          <span class="block text-xs font-semibold uppercase tracking-[0.18em] text-[#4ADE80]">
            OUR PROCESS
          </span>

          <h2 class="text-4xl md:text-5xl font-medium tracking-[-0.03em] leading-[1.08] text-white">
            Five steps.<br />
            <span class="font-serif italic font-normal text-white/95">Start to finish.</span>
          </h2>

          <p class="text-base text-white/75 max-w-md leading-relaxed">
            We handle everything from the first laser measure to the final polish. You always know
            what happens next and what it costs.
          </p>

          <button
            type="button"
            class="mt-2 rounded-full bg-[#34D361] text-[#141A15] px-6 py-3.5 text-sm font-semibold hover:bg-[#4ADE80] transition-colors shadow-lg cursor-pointer"
            @click="scrollToSchedule"
          >
            Get a free estimate
          </button>
        </div>

        <!-- Right Vertical Timeline with GSAP Scrubbed Green Progress Line & Nodes -->
        <div ref="timelineColRef" class="lg:col-span-7 relative pl-8 md:pl-12">
          <!-- Static Track Line -->
          <div
            class="absolute left-[15px] md:left-[19px] top-3 bottom-3 w-[2px] bg-[#2C3A2E]"
            aria-hidden="true"
          />
          <!-- Active Scrubbed Progress Fill -->
          <div
            ref="railFillRef"
            class="absolute left-[15px] md:left-[19px] top-3 bottom-3 w-[2px] bg-[#34D361] origin-top"
            aria-hidden="true"
          />

          <div class="space-y-8">
            <div
              v-for="(step, idx) in PROCESS_STEPS"
              :key="step.number"
              class="process-step-card relative rounded-2xl border p-6 md:p-7 transition-colors duration-300 cursor-pointer"
              :class="
                idx <= activeIndex
                  ? 'bg-[#253227]/85 border-[#4ADE80]/35'
                  : 'bg-[#1F2C20]/60 border-white/10 hover:border-white/20'
              "
              @mouseenter="activeIndex = idx"
              @click="activeIndex = idx"
            >
              <!-- Node Marker -->
              <span
                class="absolute -left-[25px] md:-left-[37px] top-7 w-4 h-4 rounded-full border-2 transition-all duration-300"
                :class="
                  idx <= activeIndex
                    ? 'bg-[#34D361] border-[#1A241B] ring-4 ring-[#34D361]/25 scale-110'
                    : 'bg-[#1A241B] border-[#5A655C] scale-100'
                "
                aria-hidden="true"
              />

              <div class="flex items-center gap-3 mb-2">
                <span class="text-xs font-semibold tracking-wider text-[#4ADE80]">
                  STEP {{ step.number }}
                </span>
              </div>
              <h3 class="text-xl md:text-2xl font-medium text-white tracking-tight">
                {{ step.title }}
              </h3>
              <p class="mt-2 text-sm md:text-base text-white/75 leading-relaxed">
                {{ step.description }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
