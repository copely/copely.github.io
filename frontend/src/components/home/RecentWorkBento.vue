<script setup lang="ts">
import { ref } from 'vue';
import { ArrowRight } from 'lucide-vue-next';
import { useGsapScope } from '@/composables/useGsapScope';
import type { RecentWorkItem } from '@/types/copely';

interface Props {
  items: RecentWorkItem[];
}

defineProps<Props>();

const sectionRef = ref<HTMLElement | null>(null);

useGsapScope(sectionRef, ({ gsap }) => {
  gsap.from('.bento-card-anim', {
    y: 28,
    scale: 0.97,
    opacity: 0,
    duration: 0.6,
    stagger: { amount: 0.35, from: 'start' },
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.bento-grid-wrap',
      start: 'top 80%',
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
  <section ref="sectionRef" class="py-20 md:py-28 px-6 md:px-10 bg-[#F7F7F3]">
    <div class="max-w-[1320px] mx-auto grid grid-cols-1 lg:grid-cols-12 gap-12 items-start">
      <!-- Sticky Left Column -->
      <div class="lg:col-span-4 lg:sticky lg:top-28 space-y-6">
        <h2 class="text-4xl md:text-5xl font-medium tracking-[-0.03em] leading-[1.08] text-[#141A15]">
          Recent work.<br />
          <span class="font-serif italic font-normal text-[#1A241B]">
            Across the Hill Country.
          </span>
        </h2>

        <p class="text-base text-[#5A655C] leading-relaxed">
          Every job gets the same care: a clean demo, laser-flat subfloor prep, and photos of the
          finished work for your records. See our work across Texas Hill Country.
        </p>

        <div class="flex flex-wrap items-center gap-4 pt-2">
          <button
            type="button"
            class="rounded-full bg-[#1A241B] text-white px-6 py-3 text-sm font-medium hover:bg-[#273528] transition-colors shadow-sm cursor-pointer"
            @click="scrollToSchedule"
          >
            Get a free estimate
          </button>

          <RouterLink
            to="/services"
            class="inline-flex items-center gap-1.5 text-sm font-medium text-[#141A15] hover:text-[#22C55E] transition-colors"
          >
            <span>View our work</span>
            <ArrowRight class="w-4 h-4" />
          </RouterLink>
        </div>
      </div>

      <!-- Right 3-Column Asymmetric Masonry / Bento Grid -->
      <div class="bento-grid-wrap lg:col-span-8 relative">
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
          <!-- Card 1 (Tall Top-Left) -->
          <div
            v-if="items[0]"
            class="bento-card-anim group relative rounded-2xl overflow-hidden bg-[#1A241B] aspect-[4/5] md:row-span-2 md:aspect-auto md:min-h-[380px]"
          >
            <img
              :src="items[0].imageUrl"
              :alt="items[0].title"
              loading="lazy"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent" />
            <div class="absolute bottom-3.5 left-3.5 right-3.5 text-white">
              <span class="text-[11px] uppercase tracking-wider text-[#4ADE80] font-medium">
                {{ items[0].category }}
              </span>
              <p class="text-sm font-medium leading-snug mt-0.5">{{ items[0].title }}</p>
            </div>
          </div>

          <!-- Card 2 -->
          <div
            v-if="items[1]"
            class="bento-card-anim group relative rounded-2xl overflow-hidden bg-[#1A241B] aspect-[4/3]"
          >
            <img
              :src="items[1].imageUrl"
              :alt="items[1].title"
              loading="lazy"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/65 via-transparent to-transparent" />
            <div class="absolute bottom-3 left-3.5 right-3.5 text-white">
              <p class="text-xs font-medium">{{ items[1].title }}</p>
            </div>
          </div>

          <!-- Card 3 -->
          <div
            v-if="items[2]"
            class="bento-card-anim group relative rounded-2xl overflow-hidden bg-[#1A241B] aspect-[4/3]"
          >
            <img
              :src="items[2].imageUrl"
              :alt="items[2].title"
              loading="lazy"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/65 via-transparent to-transparent" />
            <div class="absolute bottom-3 left-3.5 right-3.5 text-white">
              <p class="text-xs font-medium">{{ items[2].title }}</p>
            </div>
          </div>

          <!-- Featured Testimonial Overlay Card (Center-Right spanning 2 cols) -->
          <div
            class="bento-card-anim sm:col-span-2 rounded-2xl bg-[#1A241B] text-white p-6 md:p-7 flex flex-col justify-between shadow-xl border border-white/10"
          >
            <div class="inline-flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-[#4ADE80]">
              <span aria-hidden="true">◆</span>
              <span>Copely</span>
            </div>

            <blockquote class="text-lg md:text-xl font-medium leading-snug text-white mt-4">
              “They leveled our slab the same week we called, and 1,200 sq ft of new porcelain was
              down four days later.”
            </blockquote>

            <p class="text-xs text-white/70 mt-4 font-medium">
              Claire M. · Dripping Springs, TX
            </p>
          </div>

          <!-- Card 4 -->
          <div
            v-if="items[3]"
            class="bento-card-anim group relative rounded-2xl overflow-hidden bg-[#1A241B] aspect-[4/3]"
          >
            <img
              :src="items[3].imageUrl"
              :alt="items[3].title"
              loading="lazy"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/65 via-transparent to-transparent" />
            <div class="absolute bottom-3 left-3.5 right-3.5 text-white">
              <p class="text-xs font-medium">{{ items[3].title }}</p>
            </div>
          </div>

          <!-- Card 5 -->
          <div
            v-if="items[4]"
            class="bento-card-anim group relative rounded-2xl overflow-hidden bg-[#1A241B] aspect-[4/3]"
          >
            <img
              :src="items[4].imageUrl"
              :alt="items[4].title"
              loading="lazy"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/65 via-transparent to-transparent" />
            <div class="absolute bottom-3 left-3.5 right-3.5 text-white">
              <p class="text-xs font-medium">{{ items[4].title }}</p>
            </div>
          </div>

          <!-- Card 6 -->
          <div
            v-if="items[5]"
            class="bento-card-anim group relative rounded-2xl overflow-hidden bg-[#1A241B] aspect-[4/3]"
          >
            <img
              :src="items[5].imageUrl"
              :alt="items[5].title"
              loading="lazy"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/65 via-transparent to-transparent" />
            <div class="absolute bottom-3 left-3.5 right-3.5 text-white">
              <p class="text-xs font-medium">{{ items[5].title }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
