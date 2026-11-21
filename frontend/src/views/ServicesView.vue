<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ArrowUpRight, ShieldCheck } from 'lucide-vue-next';
import { useCopelyContent } from '@/composables/useCopelyContent';
import { useEstimateDraft } from '@/composables/useEstimateDraft';
import { useGsapScope } from '@/composables/useGsapScope';
import { useRouter } from 'vue-router';
import type { TileServiceCategory } from '@/types/copely';

const router = useRouter();
const { content, loadContent } = useCopelyContent();
const { setDraftSelection } = useEstimateDraft();
const pageRef = ref<HTMLElement | null>(null);

useGsapScope(pageRef, ({ gsap }) => {
  gsap.from('.services-header-anim', { y: 24, opacity: 0, duration: 0.6, ease: 'power3.out' });
  gsap.from('.services-card-anim', {
    y: 28,
    opacity: 0,
    duration: 0.55,
    stagger: 0.08,
    ease: 'power3.out',
    delay: 0.12,
  });
});

onMounted(() => {
  void loadContent();
});

function bookService(slug: TileServiceCategory): void {
  setDraftSelection(slug, 'this-week');
  void router.push({ path: '/', hash: '#schedule' });
}
</script>

<template>
  <section ref="pageRef" class="pt-32 pb-16 px-6 md:px-10 bg-[#F7F7F3]">
    <div class="max-w-[1320px] mx-auto">
      <!-- Header -->
      <div class="services-header-anim grid grid-cols-1 lg:grid-cols-12 gap-8 items-end pb-12 border-b border-[#E2E6E0]">
        <div class="lg:col-span-7 space-y-4">
          <span class="text-xs font-semibold uppercase tracking-[0.18em] text-[#22C55E]">
            COPELY SERVICES
          </span>
          <h1 class="text-4xl sm:text-5xl md:text-6xl font-medium tracking-[-0.03em] leading-[1.06] text-[#141A15]">
            Built from the subfloor up.
            <span class="block font-serif italic font-normal text-[#1A241B]">
              Backed for 25 years.
            </span>
          </h1>
        </div>
        <div class="lg:col-span-5">
          <p class="text-base text-[#5A655C] leading-relaxed">
            Every surface system we install across the Texas Hill Country pairs ANSI-spec substrate
            prep, waterproof membranes, and kiln-fired porcelain or natural stone.
          </p>
        </div>
      </div>

      <!-- Services Grid -->
      <div class="mt-12 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-7">
        <article
          v-for="service in content.services"
          :key="service.id"
          class="services-card-anim group rounded-3xl bg-white border border-[#E2E6E0] overflow-hidden flex flex-col justify-between shadow-sm hover:shadow-xl transition-shadow duration-300"
        >
          <div>
            <RouterLink
              :to="`/services/${service.slug}`"
              class="block aspect-[4/3] overflow-hidden bg-[#1A241B]"
            >
              <img
                :src="service.imageUrl"
                :alt="service.title"
                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
              />
            </RouterLink>

            <div class="p-6 md:p-7 space-y-3">
              <div class="flex items-center justify-between text-xs text-[#5A655C]">
                <span class="font-medium text-[#22C55E] inline-flex items-center gap-1">
                  <ShieldCheck class="w-4 h-4" />
                  {{ service.warrantyYears }}-Year Warranty
                </span>
                <span>{{ service.timelineLabel }}</span>
              </div>

              <h2 class="text-2xl font-medium text-[#141A15] tracking-tight">
                <RouterLink :to="`/services/${service.slug}`" class="hover:text-[#22C55E] transition">
                  {{ service.title }}
                </RouterLink>
              </h2>

              <p class="text-sm text-[#5A655C] leading-relaxed">
                {{ service.description }}
              </p>
            </div>
          </div>

          <div class="px-6 md:px-7 pb-6 pt-2 flex items-center justify-between gap-3">
            <RouterLink
              :to="`/services/${service.slug}`"
              class="inline-flex items-center gap-1.5 text-sm font-semibold text-[#141A15] hover:text-[#22C55E] transition"
            >
              <span>Inspect specs</span>
              <ArrowUpRight class="w-4 h-4" />
            </RouterLink>
            <button
              type="button"
              class="rounded-full bg-[#1A241B] text-white px-4 py-2 text-xs font-medium hover:bg-[#273528] transition cursor-pointer"
              @click="bookService(service.slug)"
            >
              Estimate this
            </button>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>
