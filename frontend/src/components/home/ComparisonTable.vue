<script setup lang="ts">
import { ref } from 'vue';
import { Check, X } from 'lucide-vue-next';
import { COMPARISON_ROWS } from '@/data/copelyData';
import { useGsapScope } from '@/composables/useGsapScope';

const sectionRef = ref<HTMLElement | null>(null);

useGsapScope(sectionRef, ({ gsap }) => {
  gsap.from('.compare-matrix-wrap', {
    y: 26,
    opacity: 0,
    duration: 0.65,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.compare-matrix-wrap',
      start: 'top 82%',
      once: true,
    },
  });

  gsap.from('.compare-warranty-banner', {
    y: 24,
    opacity: 0,
    duration: 0.6,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.compare-warranty-banner',
      start: 'top 86%',
      once: true,
    },
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
    <div class="rounded-3xl bg-[#1A241B] text-white py-20 md:py-24 px-5 sm:px-8 md:px-12">
      <div class="max-w-[1160px] mx-auto">
        <h2 class="text-4xl md:text-5xl font-medium tracking-[-0.03em] text-center text-white mb-12">
          Why Copely?
        </h2>

        <!-- 5 Rows x 4 Columns Comparison Matrix -->
        <div class="compare-matrix-wrap overflow-x-auto pb-2">
          <div class="min-w-[720px] grid grid-cols-12 items-stretch">
            <!-- Header Row -->
            <div class="col-span-3 py-4 px-4 text-xs uppercase tracking-wider text-white/45 flex items-end">
              Compare options
            </div>
            <div
              class="col-span-3 bg-[#34D361] text-[#141A15] rounded-t-2xl py-5 px-5 flex items-center gap-2 font-semibold text-base"
            >
              <span aria-hidden="true">◆</span>
              <span>Copely</span>
            </div>
            <div
              class="col-span-2 py-5 px-4 flex items-center justify-center text-sm font-medium text-white/75 border-b border-[#2C3A2E]"
            >
              DIY
            </div>
            <div
              class="col-span-2 py-5 px-4 flex items-center justify-center text-sm font-medium text-white/75 border-b border-[#2C3A2E] text-center"
            >
              General Contractor
            </div>
            <div
              class="col-span-2 py-5 px-4 flex items-center justify-center text-sm font-medium text-white/75 border-b border-[#2C3A2E] text-center"
            >
              Big-Box Franchise
            </div>

            <!-- Matrix Body Rows -->
            <template v-for="(row, rIdx) in COMPARISON_ROWS" :key="row.label">
              <!-- Feature Label -->
              <div
                class="col-span-3 py-4 px-4 flex items-center text-sm font-medium text-white/85 border-b border-[#2C3A2E]"
                :class="rIdx === COMPARISON_ROWS.length - 1 ? 'border-b-0' : ''"
              >
                {{ row.label }}
              </div>

              <!-- Highlighted Copely Column -->
              <div
                class="col-span-3 bg-[#34D361] text-[#141A15] py-4 px-5 flex items-center gap-2 text-sm font-semibold border-b border-[#1A241B]/10"
                :class="rIdx === COMPARISON_ROWS.length - 1 ? 'rounded-b-2xl border-b-0' : ''"
              >
                <Check class="w-4 h-4 shrink-0 stroke-[2.5]" />
                <span>{{ row.copely.text }}</span>
              </div>

              <!-- DIY -->
              <div
                class="col-span-2 py-4 px-4 flex items-center justify-center gap-1.5 text-sm text-white/70 border-b border-[#2C3A2E] text-center"
                :class="rIdx === COMPARISON_ROWS.length - 1 ? 'border-b-0' : ''"
              >
                <X
                  v-if="row.diy.status === 'negative'"
                  class="w-3.5 h-3.5 text-rose-400 shrink-0"
                />
                <span :class="row.diy.status === 'negative' ? 'text-white/55' : ''">
                  {{ row.diy.text }}
                </span>
              </div>

              <!-- General Contractor -->
              <div
                class="col-span-2 py-4 px-4 flex items-center justify-center gap-1.5 text-sm text-white/70 border-b border-[#2C3A2E] text-center"
                :class="rIdx === COMPARISON_ROWS.length - 1 ? 'border-b-0' : ''"
              >
                <X
                  v-if="row.generalContractor.status === 'negative'"
                  class="w-3.5 h-3.5 text-rose-400 shrink-0"
                />
                <span :class="row.generalContractor.status === 'negative' ? 'text-white/55' : ''">
                  {{ row.generalContractor.text }}
                </span>
              </div>

              <!-- Big-Box Franchise -->
              <div
                class="col-span-2 py-4 px-4 flex items-center justify-center gap-1.5 text-sm text-white/70 border-b border-[#2C3A2E] text-center"
                :class="rIdx === COMPARISON_ROWS.length - 1 ? 'border-b-0' : ''"
              >
                <X
                  v-if="row.bigBox.status === 'negative'"
                  class="w-3.5 h-3.5 text-rose-400 shrink-0"
                />
                <span :class="row.bigBox.status === 'negative' ? 'text-white/55' : ''">
                  {{ row.bigBox.text }}
                </span>
              </div>
            </template>
          </div>
        </div>

        <!-- Bottom Guarantee Banner (Light Card inside Dark Container) -->
        <div
          class="compare-warranty-banner mt-12 rounded-2xl bg-[#F7F7F3] text-[#141A15] p-6 md:p-8 flex flex-col md:flex-row md:items-center justify-between gap-6 shadow-lg"
        >
          <div class="space-y-2 max-w-2xl">
            <h3 class="text-xl md:text-2xl font-medium tracking-tight text-[#141A15]">
              And every installation is backed by a
              <span class="bg-[#34D361]/30 px-2 py-0.5 rounded-lg font-semibold">
                25-year warranty
              </span>
            </h3>
            <p class="text-sm text-[#5A655C] leading-relaxed">
              With Copely you get a licensed local crew, quality materials, and a free same-day
              estimate — all backed in writing, with no surprise bills.
            </p>
          </div>

          <button
            type="button"
            class="shrink-0 rounded-full bg-[#34D361] text-[#141A15] px-6 py-3.5 text-sm font-semibold hover:bg-[#4ADE80] transition-colors cursor-pointer"
            @click="scrollToSchedule"
          >
            Get a free estimate
          </button>
        </div>
      </div>
    </div>
  </section>
</template>
