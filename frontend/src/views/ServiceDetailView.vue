<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ArrowLeft, Check, Clock, ShieldCheck } from 'lucide-vue-next';
import { useCopelyContent } from '@/composables/useCopelyContent';
import { useEstimateDraft } from '@/composables/useEstimateDraft';
import { useGsapScope } from '@/composables/useGsapScope';
import { fetchServiceBySlug } from '@/services/api';
import type { ServiceItem } from '@/types/copely';

const route = useRoute();
const router = useRouter();
const { content, loadContent } = useCopelyContent();
const { setDraftSelection } = useEstimateDraft();
const pageRef = ref<HTMLElement | null>(null);

useGsapScope(pageRef, ({ gsap }) => {
  gsap.from('.svc-detail-hero', { y: 24, opacity: 0, duration: 0.6, stagger: 0.08, ease: 'power3.out' });
  gsap.from('.svc-detail-panel', {
    y: 24,
    opacity: 0,
    duration: 0.6,
    stagger: 0.1,
    ease: 'power3.out',
    scrollTrigger: { trigger: '.svc-detail-panels', start: 'top 84%', once: true },
  });
});

const fetchedService = ref<ServiceItem | null>(null);

const slug = computed<string>(() => String(route.params.slug ?? ''));

const service = computed<ServiceItem | undefined>(() => {
  if (fetchedService.value && fetchedService.value.slug === slug.value) {
    return fetchedService.value;
  }
  return content.value.services.find((s) => s.slug === slug.value) ?? content.value.services[0];
});

const otherServices = computed<ServiceItem[]>(() => {
  return content.value.services.filter((s) => s.slug !== service.value?.slug);
});

async function loadServiceDetail(): Promise<void> {
  await loadContent();
  try {
    fetchedService.value = await fetchServiceBySlug(slug.value);
  } catch {
    fetchedService.value = null;
  }
}

function scheduleForService(): void {
  if (service.value) {
    setDraftSelection(service.value.slug, 'this-week');
  }
  void router.push({ path: '/', hash: '#schedule' });
}

watch(slug, () => {
  void loadServiceDetail();
});

onMounted(() => {
  void loadServiceDetail();
});
</script>

<template>
  <section v-if="service" ref="pageRef" class="pt-32 pb-16 px-6 md:px-10 bg-[#F7F7F3]">
    <div class="max-w-[1320px] mx-auto space-y-16">
      <!-- Breadcrumb -->
      <div>
        <RouterLink
          to="/services"
          class="inline-flex items-center gap-2 text-sm font-medium text-[#5A655C] hover:text-[#141A15] transition"
        >
          <ArrowLeft class="w-4 h-4" />
          <span>All Copely services</span>
        </RouterLink>
      </div>

      <!-- Hero Split -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-10 lg:gap-14 items-center">
        <div class="svc-detail-hero lg:col-span-6 space-y-6">
          <div class="inline-flex flex-wrap items-center gap-3">
            <span
              class="inline-flex items-center gap-1.5 rounded-full bg-[#DCE8D9] px-3.5 py-1 text-xs font-semibold text-[#141A15]"
            >
              <ShieldCheck class="w-4 h-4 text-[#22C55E]" />
              {{ service.warrantyYears }}-Year Transferable Warranty
            </span>
            <span
              class="inline-flex items-center gap-1.5 rounded-full bg-white border border-[#E2E6E0] px-3.5 py-1 text-xs font-medium text-[#5A655C]"
            >
              <Clock class="w-3.5 h-3.5 text-[#22C55E]" />
              {{ service.timelineLabel }}
            </span>
          </div>

          <h1 class="text-4xl sm:text-5xl md:text-6xl font-medium tracking-[-0.03em] leading-[1.06] text-[#141A15]">
            {{ service.title }}
            <span class="block font-serif italic font-normal text-3xl sm:text-4xl text-[#5A655C] mt-2">
              {{ service.shortSubtitle }}
            </span>
          </h1>

          <p class="text-base md:text-lg text-[#5A655C] leading-relaxed">
            {{ service.description }}
          </p>

          <div class="rounded-2xl bg-[#E8EFE6] p-5 border border-[#DCE8D9]">
            <p class="text-xs font-semibold uppercase tracking-wider text-[#5A655C]">Ideal for</p>
            <p class="text-sm font-medium text-[#141A15] mt-1">{{ service.idealFor }}</p>
          </div>

          <div class="flex flex-wrap items-center gap-4 pt-2">
            <button
              type="button"
              class="rounded-full bg-[#1A241B] text-white px-7 py-3.5 text-sm font-medium hover:bg-[#273528] transition shadow-md cursor-pointer"
              @click="scheduleForService"
            >
              Schedule a free measurement
            </button>
            <RouterLink
              to="/contact"
              class="rounded-full bg-white border border-[#E2E6E0] text-[#141A15] px-6 py-3.5 text-sm font-medium hover:border-[#1A241B] transition"
            >
              Request written quote
            </RouterLink>
          </div>
        </div>

        <div class="svc-detail-hero lg:col-span-6">
          <div class="rounded-3xl overflow-hidden bg-[#1A241B] aspect-[4/3] shadow-xl">
            <img
              :src="service.imageUrl"
              :alt="service.title"
              class="w-full h-full object-cover"
            />
          </div>
        </div>
      </div>

      <!-- Specs & Installation Highlights -->
      <div class="svc-detail-panels grid grid-cols-1 lg:grid-cols-12 gap-8">
        <div class="svc-detail-panel lg:col-span-7 rounded-3xl bg-[#1A241B] text-white p-8 md:p-10 space-y-6">
          <span class="text-xs font-semibold uppercase tracking-[0.18em] text-[#4ADE80]">
            SYSTEM SPECIFICATIONS
          </span>
          <h2 class="text-3xl font-medium tracking-tight">What goes beneath the tile</h2>
          <ul class="space-y-4">
            <li
              v-for="spec in service.specs"
              :key="spec"
              class="flex items-start gap-3 text-sm md:text-base text-white/85"
            >
              <span
                class="w-6 h-6 rounded-full bg-[#34D361] text-[#141A15] flex items-center justify-center shrink-0 mt-0.5"
              >
                <Check class="w-3.5 h-3.5 stroke-[2.5]" />
              </span>
              <span>{{ spec }}</span>
            </li>
          </ul>
        </div>

        <div class="svc-detail-panel lg:col-span-5 rounded-3xl bg-white border border-[#E2E6E0] p-8 md:p-10 space-y-6">
          <span class="text-xs font-semibold uppercase tracking-[0.18em] text-[#22C55E]">
            INSTALL PROTOCOL
          </span>
          <h2 class="text-3xl font-medium tracking-tight text-[#141A15]">On-site craftsmanship</h2>
          <ul class="space-y-4">
            <li
              v-for="(step, idx) in service.processHighlights"
              :key="step"
              class="flex items-start gap-3.5"
            >
              <span
                class="w-7 h-7 rounded-full bg-[#E8EFE6] text-[#141A15] text-xs font-semibold flex items-center justify-center shrink-0"
              >
                0{{ idx + 1 }}
              </span>
              <span class="text-sm md:text-base text-[#364137] pt-0.5">{{ step }}</span>
            </li>
          </ul>
        </div>
      </div>

      <!-- Explore Other Services -->
      <div class="pt-6 border-t border-[#E2E6E0]">
        <h3 class="text-xl font-medium text-[#141A15] mb-5">Explore other Copely services</h3>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <RouterLink
            v-for="item in otherServices"
            :key="item.id"
            :to="`/services/${item.slug}`"
            class="rounded-2xl bg-white border border-[#E2E6E0] p-5 hover:border-[#34D361] transition-colors"
          >
            <p class="font-medium text-[#141A15]">{{ item.title }}</p>
            <p class="text-xs text-[#5A655C] mt-1">{{ item.shortSubtitle }}</p>
          </RouterLink>
        </div>
      </div>
    </div>
  </section>
</template>
