<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import {
  AccordionContent,
  AccordionHeader,
  AccordionItem,
  AccordionRoot,
  AccordionTrigger,
} from 'radix-vue';
import { Minus, Phone, Plus, Star } from 'lucide-vue-next';
import { useCopelyContent } from '@/composables/useCopelyContent';
import { useGsapScope } from '@/composables/useGsapScope';

const router = useRouter();
const { content, loadContent } = useCopelyContent();
const pageRef = ref<HTMLElement | null>(null);

useGsapScope(pageRef, ({ gsap }) => {
  gsap.from('.faq-hero-anim', { y: 22, opacity: 0, duration: 0.6, ease: 'power3.out' });
  gsap.from('.faq-review-card', {
    y: 24,
    opacity: 0,
    duration: 0.55,
    stagger: 0.07,
    ease: 'power3.out',
    scrollTrigger: { trigger: '.faq-reviews-masonry', start: 'top 82%', once: true },
  });
  gsap.from('.glossary-card-anim', {
    y: 20,
    opacity: 0,
    duration: 0.5,
    stagger: 0.05,
    ease: 'power3.out',
    scrollTrigger: { trigger: '.glossary-grid-wrap', start: 'top 84%', once: true },
  });
});

const openByCategory = reactive<Record<string, string>>({
  estimates: 'faq-1',
  'during-job': 'faq-2',
  system: 'faq-sys-1',
  waterproofing: 'faq-5',
  warranty: 'faq-4',
});

onMounted(() => {
  void loadContent();
});

function goToSchedule(): void {
  void router.push({ path: '/', hash: '#schedule' });
}

function scrollToSection(id: string): void {
  const el = document.getElementById(`faq-cat-${id}`);
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }
}
</script>

<template>
  <div ref="pageRef" class="pt-32 pb-16 bg-[#F7F7F3]">
    <!-- Header Split -->
    <section class="faq-hero-anim px-6 md:px-10 max-w-[1320px] mx-auto pb-14 border-b border-[#E2E6E0]">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-end">
        <div class="lg:col-span-7">
          <h1 class="text-4xl sm:text-5xl md:text-6xl font-medium tracking-[-0.03em] leading-[1.06] text-[#141A15]">
            Common questions<br />
            <span class="font-serif italic font-normal text-[#1A241B]">from homeowners.</span>
          </h1>
        </div>
        <div class="lg:col-span-5 space-y-5">
          <p class="text-base text-[#5A655C] leading-relaxed">
            Costs, timelines, waterproofing, warranties, and what to expect on install day —
            answered the way we’d explain it to a neighbor.
          </p>
          <div class="flex flex-wrap items-center gap-3">
            <button
              type="button"
              class="rounded-full bg-[#1A241B] text-white px-5 py-3 text-sm font-medium hover:bg-[#273528] transition cursor-pointer"
              @click="goToSchedule"
            >
              Get a free estimate
            </button>
            <a
              href="tel:+15125550142"
              class="inline-flex items-center gap-2 rounded-full bg-white border border-[#E2E6E0] px-5 py-3 text-sm font-medium text-[#141A15] hover:border-[#1A241B] transition"
            >
              <Phone class="w-4 h-4 text-[#22C55E]" />
              <span>(512) 555-0142</span>
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- Categorized FAQ Layout with Sticky Left Sidebar -->
    <section class="px-6 md:px-10 max-w-[1320px] mx-auto py-16 md:py-20">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 items-start">
        <!-- Left Sticky Sidebar Navigation -->
        <aside class="lg:col-span-4 lg:sticky lg:top-28 space-y-6">
          <nav aria-label="FAQ categories" class="space-y-1.5">
            <button
              v-for="cat in content.faqCategories"
              :key="cat.id"
              type="button"
              class="w-full text-left rounded-xl px-4 py-3 text-sm font-medium text-[#141A15] hover:bg-[#E8EFE6] transition-colors cursor-pointer"
              @click="scrollToSection(cat.id)"
            >
              {{ cat.title }}
            </button>
          </nav>

          <!-- Floating Google Rating Card at bottom of sidebar -->
          <div class="rounded-2xl bg-white border border-[#E2E6E0] p-5 shadow-sm space-y-2">
            <div class="flex items-center gap-2">
              <span class="text-2xl font-semibold text-[#141A15]">4.9</span>
              <div class="flex items-center gap-0.5 text-[#22C55E]">
                <Star v-for="n in 5" :key="n" class="w-4 h-4 fill-current" />
              </div>
            </div>
            <p class="text-xs text-[#5A655C]">from 1,200+ Google reviews</p>
          </div>
        </aside>

        <!-- Right Column: 2-Column Accordion Grids per Category -->
        <div class="lg:col-span-8 space-y-14">
          <div
            v-for="cat in content.faqCategories"
            :id="`faq-cat-${cat.id}`"
            :key="cat.id"
            class="scroll-mt-28 space-y-5"
          >
            <h2 class="text-2xl md:text-3xl font-medium tracking-tight text-[#141A15]">
              {{ cat.title }}
            </h2>

            <AccordionRoot
              v-model="openByCategory[cat.id]"
              type="single"
              collapsible
              class="grid grid-cols-1 md:grid-cols-2 gap-4 items-start"
            >
              <AccordionItem
                v-for="item in cat.items"
                :key="item.id"
                :value="item.id"
                class="rounded-2xl border transition-colors overflow-hidden"
                :class="
                  openByCategory[cat.id] === item.id
                    ? 'bg-[#DCE8D9] border-[#C4D7BF]'
                    : 'bg-white border-[#E2E6E0]'
                "
              >
                <AccordionHeader>
                  <AccordionTrigger
                    class="w-full p-5 flex items-start justify-between gap-3 text-left cursor-pointer focus:outline-none"
                  >
                    <span class="text-sm md:text-base font-medium text-[#141A15] leading-snug">
                      {{ item.question }}
                    </span>
                    <span
                      class="w-6 h-6 rounded-full flex items-center justify-center shrink-0 mt-0.5"
                      :class="
                        openByCategory[cat.id] === item.id
                          ? 'bg-[#1A241B] text-white'
                          : 'bg-[#F7F7F3] text-[#141A15]'
                      "
                      aria-hidden="true"
                    >
                      <Minus v-if="openByCategory[cat.id] === item.id" class="w-3.5 h-3.5" />
                      <Plus v-else class="w-3.5 h-3.5" />
                    </span>
                  </AccordionTrigger>
                </AccordionHeader>
                <AccordionContent class="px-5 pb-5 text-sm text-[#364137] leading-relaxed">
                  {{ item.answer }}
                </AccordionContent>
              </AccordionItem>
            </AccordionRoot>
          </div>
        </div>
      </div>
    </section>

    <!-- Google Reviews Masonry Section -->
    <section class="px-3 md:px-4">
      <div class="rounded-3xl bg-[#1A241B] text-white py-20 px-6 md:px-12">
        <div class="max-w-[1320px] mx-auto grid grid-cols-1 lg:grid-cols-12 gap-12 items-start">
          <!-- Left Summary -->
          <div class="lg:col-span-4 lg:sticky lg:top-28 space-y-5">
            <h2 class="text-3xl sm:text-4xl font-medium tracking-[-0.03em] leading-[1.1]">
              Don’t take our word.<br />
              <span class="font-serif italic font-normal text-white/90">Take theirs.</span>
            </h2>

            <div class="rounded-2xl bg-[#253227] border border-white/10 p-5 space-y-2">
              <div class="flex items-center gap-2.5">
                <span class="text-3xl font-semibold text-white">4.9</span>
                <div class="flex items-center gap-0.5 text-[#4ADE80]">
                  <Star v-for="n in 5" :key="n" class="w-4 h-4 fill-current" />
                </div>
              </div>
              <p class="text-xs text-white/70">Verified Hill Country homeowner reviews</p>
            </div>

            <a
              href="https://www.google.com/maps"
              target="_blank"
              rel="noreferrer"
              class="inline-flex items-center gap-2.5 rounded-full bg-white text-[#141A15] px-5 py-3 text-sm font-semibold hover:bg-[#E8EFE6] transition"
            >
              <span class="font-bold text-[#22C55E]">G</span>
              <span>Write a review</span>
            </a>
          </div>

          <!-- Right 2-Column Masonry Grid -->
          <div class="faq-reviews-masonry lg:col-span-8 columns-1 md:columns-2 gap-5 space-y-5">
            <article
              v-for="review in content.googleReviews"
              :key="review.id"
              class="faq-review-card break-inside-avoid rounded-2xl bg-[#253227]/85 border border-white/10 p-6 space-y-4"
            >
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <span
                    class="w-10 h-10 rounded-full bg-[#34D361] text-[#141A15] font-semibold text-xs flex items-center justify-center"
                  >
                    {{ review.initials }}
                  </span>
                  <div>
                    <p class="text-sm font-semibold text-white">{{ review.author }}</p>
                    <p class="text-xs text-white/55">{{ review.projectTag }}</p>
                  </div>
                </div>
                <span class="text-xs text-white/45">{{ review.dateLabel }}</span>
              </div>

              <div class="flex items-center gap-0.5 text-[#4ADE80]">
                <Star v-for="n in review.rating" :key="n" class="w-3.5 h-3.5 fill-current" />
              </div>

              <p class="text-sm text-white/85 leading-relaxed">“{{ review.body }}”</p>
            </article>
          </div>
        </div>
      </div>
    </section>

    <!-- Glossary Section: Commonly-used tiling terms -->
    <section class="px-6 md:px-10 max-w-[1320px] mx-auto pt-20 md:pt-24">
      <div class="max-w-xl space-y-3 mb-10">
        <h2 class="text-3xl sm:text-4xl font-medium tracking-[-0.03em] text-[#141A15]">
          Commonly-used tiling terms.
        </h2>
        <p class="text-base text-[#5A655C]">
          The words you’ll see on your estimate, each explained in one simple sentence.
        </p>
      </div>

      <div class="glossary-grid-wrap grid grid-cols-1 md:grid-cols-3 gap-5">
        <div
          v-for="item in content.glossaryTerms"
          :key="item.id"
          class="glossary-card-anim rounded-2xl bg-white border border-[#E2E6E0] p-6 space-y-2"
        >
          <h3 class="text-lg font-medium text-[#141A15]">{{ item.term }}</h3>
          <p class="text-sm text-[#5A655C] leading-relaxed">{{ item.definition }}</p>
        </div>
      </div>
    </section>
  </div>
</template>
