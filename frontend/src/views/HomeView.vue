<script setup lang="ts">
import { nextTick, onMounted, onUnmounted } from 'vue';
import HeroSection from '@/components/home/HeroSection.vue';
import TileSystemCutaway from '@/components/home/TileSystemCutaway.vue';
import ServiceCarousel from '@/components/home/ServiceCarousel.vue';
import RecentWorkBento from '@/components/home/RecentWorkBento.vue';
import ProcessTimeline from '@/components/home/ProcessTimeline.vue';
import BeforeAfterSliders from '@/components/home/BeforeAfterSliders.vue';
import TestimonialsSection from '@/components/home/TestimonialsSection.vue';
import ComparisonTable from '@/components/home/ComparisonTable.vue';
import HomeFaqAccordion from '@/components/home/HomeFaqAccordion.vue';
import InlineBookingCalendar from '@/components/home/InlineBookingCalendar.vue';
import { useCopelyContent } from '@/composables/useCopelyContent';
import { useSiteExplorer } from '@/composables/useSiteExplorer';

const { content, loadContent } = useCopelyContent();
const { registerHomeSections, unregisterHomeSections } = useSiteExplorer();

onMounted(async () => {
  await loadContent();
  await nextTick();
  registerHomeSections();
});

onUnmounted(() => {
  unregisterHomeSections();
});
</script>

<template>
  <div>
    <div id="hero" class="scroll-mt-24">
      <HeroSection />
    </div>
    <div id="system" class="scroll-mt-20">
      <TileSystemCutaway />
    </div>
    <div id="services" class="scroll-mt-20">
      <ServiceCarousel :services="content.services" />
    </div>
    <div id="recent-work" class="scroll-mt-20">
      <RecentWorkBento :items="content.recentWork" />
    </div>
    <div id="process" class="scroll-mt-20">
      <ProcessTimeline />
    </div>
    <div id="before-after" class="scroll-mt-20">
      <BeforeAfterSliders :cases="content.beforeAfterCases" />
    </div>
    <div id="testimonials" class="scroll-mt-20">
      <TestimonialsSection :testimonials="content.testimonials" />
    </div>
    <div id="why-copely" class="scroll-mt-20">
      <ComparisonTable />
    </div>
    <div id="faq" class="scroll-mt-20">
      <HomeFaqAccordion :faqs="content.homeFaqs" />
    </div>
    <InlineBookingCalendar :towns="content.towns" />
  </div>
</template>
