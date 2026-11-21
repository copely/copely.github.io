<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import {
  AccordionContent,
  AccordionHeader,
  AccordionItem,
  AccordionRoot,
  AccordionTrigger,
} from 'radix-vue';
import { ArrowRight, Minus, Plus } from 'lucide-vue-next';
import { allowsMotion, useGsapScope } from '@/composables/useGsapScope';
import type { FaqItem } from '@/types/copely';

interface Props {
  faqs: FaqItem[];
}

const props = defineProps<Props>();

const openId = ref<string>(props.faqs[0]?.id ?? 'faq-1');
const sectionRef = ref<HTMLElement | null>(null);
const hydrated = ref<boolean>(false);
const motionEnabled = ref<boolean>(true);
let interactiveCtx: gsap.Context | null = null;

function queryFaqPart<T extends HTMLElement = HTMLElement>(attr: string, id: string): T | null {
  return sectionRef.value?.querySelector<T>(`[${attr}="${id}"]`) ?? null;
}

function handleRowPointerMove(event: PointerEvent): void {
  if (!allowsMotion() || event.pointerType === 'touch') {
    return;
  }
  const card = event.currentTarget as HTMLElement | null;
  if (!card) {
    return;
  }
  const rect = card.getBoundingClientRect();
  card.style.setProperty('--faq-mx', `${event.clientX - rect.left}px`);
  card.style.setProperty('--faq-my', `${event.clientY - rect.top}px`);
}

function handleRowHover(itemId: string, entering: boolean): void {
  if (!allowsMotion()) {
    return;
  }
  const sheen = queryFaqPart('data-faq-sheen', itemId);
  const badge = queryFaqPart('data-faq-badge', itemId);
  if (sheen) {
    gsap.killTweensOf(sheen);
    gsap.to(sheen, {
      opacity: entering ? 1 : 0,
      duration: entering ? 0.24 : 0.32,
      ease: 'power2.out',
    });
  }
  if (badge && openId.value !== itemId) {
    gsap.killTweensOf(badge);
    gsap.to(badge, {
      scale: entering ? 1.1 : 1,
      x: entering ? 2 : 0,
      duration: 0.26,
      ease: entering ? 'back.out(2)' : 'power2.out',
    });
  }
}

function handleLinkHover(entering: boolean): void {
  if (!allowsMotion()) {
    return;
  }
  const arrow = sectionRef.value?.querySelector('.faq-bottom-arrow');
  if (arrow) {
    gsap.killTweensOf(arrow);
    gsap.to(arrow, {
      x: entering ? 6 : 0,
      duration: 0.28,
      ease: entering ? 'back.out(2.2)' : 'power2.out',
    });
  }
}

useGsapScope(sectionRef, ({ gsap: g }) => {
  const panels = sectionRef.value?.querySelectorAll<HTMLElement>('[data-faq-panel]');
  panels?.forEach((panel) => {
    const id = panel.dataset.faqPanel ?? '';
    const isInitialOpen = id === openId.value;
    g.set(panel, {
      height: isInitialOpen ? 'auto' : 0,
      autoAlpha: isInitialOpen ? 1 : 0,
    });
  });

  const rails = sectionRef.value?.querySelectorAll<HTMLElement>('[data-faq-rail]');
  rails?.forEach((rail) => {
    const id = rail.dataset.faqRail ?? '';
    const isInitialOpen = id === openId.value;
    g.set(rail, {
      scaleY: isInitialOpen ? 1 : 0,
      opacity: isInitialOpen ? 1 : 0,
    });
  });

  const badges = sectionRef.value?.querySelectorAll<HTMLElement>('[data-faq-badge]');
  badges?.forEach((badge) => {
    const id = badge.dataset.faqBadge ?? '';
    if (id === openId.value) {
      g.set(badge, { rotation: 180 });
    }
  });

  const tl = g.timeline({
    scrollTrigger: {
      trigger: sectionRef.value,
      start: 'top 78%',
      once: true,
    },
  });

  tl.fromTo(
    '.faq-head-line',
    { clipPath: 'inset(100% 0% 0% 0%)', y: 18 },
    {
      clipPath: 'inset(0% 0% 0% 0%)',
      y: 0,
      duration: 0.68,
      stagger: 0.08,
      ease: 'expo.out',
      clearProps: 'clipPath,transform',
    },
  )
    .fromTo(
      '.faq-laser-rule',
      { scaleX: 0 },
      { scaleX: 1, duration: 0.5, ease: 'power3.out' },
      '-=0.42',
    )
    .fromTo(
      '.faq-subtext',
      { y: 12, opacity: 0 },
      { y: 0, opacity: 1, duration: 0.45, ease: 'power2.out', clearProps: 'transform,opacity' },
      '-=0.35',
    )
    .fromTo(
      '.home-faq-item',
      {
        x: (i: number) => (i % 2 === 0 ? -24 : 24),
        y: 22,
        rotationX: -9,
        rotationY: (i: number) => (i % 2 === 0 ? 3 : -3),
        transformPerspective: 900,
        opacity: 0,
      },
      {
        x: 0,
        y: 0,
        rotationX: 0,
        rotationY: 0,
        opacity: 1,
        duration: 0.64,
        stagger: 0.06,
        ease: 'expo.out',
        clearProps: 'transform,opacity',
      },
      '-=0.3',
    );

  const activeRail = queryFaqPart('data-faq-rail', openId.value);
  const activeBadge = queryFaqPart('data-faq-badge', openId.value);
  if (activeRail) {
    tl.fromTo(
      activeRail,
      { scaleY: 0 },
      { scaleY: 1, duration: 0.36, ease: 'back.out(1.8)' },
      '-=0.22',
    );
  }
  if (activeBadge) {
    tl.fromTo(
      activeBadge,
      { scale: 0.72 },
      { scale: 1, duration: 0.34, ease: 'back.out(2)' },
      '-=0.3',
    );
  }

  interactiveCtx = g.context(() => {}, sectionRef.value ?? undefined);
});

watch(openId, async (newId, oldId) => {
  if (!allowsMotion() || newId === oldId) {
    return;
  }
  await nextTick();

  interactiveCtx?.add(() => {
    if (oldId) {
      const oldPanel = queryFaqPart('data-faq-panel', oldId);
      const oldCopy = queryFaqPart('data-faq-copy', oldId);
      const oldRail = queryFaqPart('data-faq-rail', oldId);
      const oldBadge = queryFaqPart('data-faq-badge', oldId);

      if (oldCopy) {
        gsap.killTweensOf(oldCopy);
        gsap.to(oldCopy, { y: -6, opacity: 0, duration: 0.2, ease: 'power2.in' });
      }
      if (oldPanel) {
        gsap.killTweensOf(oldPanel);
        gsap.set(oldPanel, { height: oldPanel.offsetHeight });
        gsap.to(oldPanel, {
          height: 0,
          autoAlpha: 0,
          duration: 0.34,
          ease: 'power3.inOut',
        });
      }
      if (oldRail) {
        gsap.killTweensOf(oldRail);
        gsap.to(oldRail, { scaleY: 0, opacity: 0, duration: 0.24, ease: 'power2.in' });
      }
      if (oldBadge) {
        gsap.killTweensOf(oldBadge);
        gsap.to(oldBadge, { rotation: 0, scale: 1, x: 0, duration: 0.34, ease: 'back.out(1.6)' });
      }
    }

    if (newId) {
      const newPanel = queryFaqPart('data-faq-panel', newId);
      const newCopy = queryFaqPart('data-faq-copy', newId);
      const newRail = queryFaqPart('data-faq-rail', newId);
      const newBadge = queryFaqPart('data-faq-badge', newId);

      if (newPanel) {
        gsap.killTweensOf(newPanel);
        const startHeight = newPanel.offsetHeight;
        const targetHeight = newPanel.scrollHeight;
        gsap.fromTo(
          newPanel,
          { height: startHeight, autoAlpha: 1 },
          {
            height: targetHeight,
            autoAlpha: 1,
            duration: 0.44,
            ease: 'expo.out',
            onComplete: () => {
              gsap.set(newPanel, { height: 'auto' });
              ScrollTrigger.refresh();
            },
          },
        );
      }
      if (newCopy) {
        gsap.killTweensOf(newCopy);
        gsap.fromTo(
          newCopy,
          { y: 10, opacity: 0 },
          {
            y: 0,
            opacity: 1,
            duration: 0.36,
            delay: 0.06,
            ease: 'power3.out',
            clearProps: 'transform,opacity',
          },
        );
      }
      if (newRail) {
        gsap.killTweensOf(newRail);
        gsap.fromTo(
          newRail,
          { scaleY: 0, opacity: 0.5 },
          { scaleY: 1, opacity: 1, duration: 0.42, ease: 'back.out(1.7)' },
        );
      }
      if (newBadge) {
        gsap.killTweensOf(newBadge);
        gsap.to(newBadge, {
          rotation: 180,
          scale: 1.08,
          x: 0,
          duration: 0.38,
          ease: 'back.out(1.8)',
          onComplete: () => {
            gsap.to(newBadge, { scale: 1, duration: 0.18, ease: 'power2.out' });
          },
        });
      }
    }
  });
});

onMounted(() => {
  motionEnabled.value = allowsMotion();
  hydrated.value = true;
});

onUnmounted(() => {
  interactiveCtx?.revert();
  interactiveCtx = null;
});
</script>

<template>
  <section ref="sectionRef" class="py-20 md:py-28 px-6 md:px-10 bg-[#F7F7F3] overflow-x-clip">
    <div class="max-w-[820px] mx-auto">
      <!-- Centered Header -->
      <div class="text-center space-y-3 mb-12">
        <h2 class="text-4xl md:text-5xl font-medium tracking-[-0.03em] leading-[1.08] text-[#141A15]">
          <span class="block overflow-hidden">
            <span class="faq-head-line block">Before you book.</span>
          </span>
          <span class="block overflow-hidden">
            <span class="faq-head-line block font-serif italic font-normal text-[#1A241B]">
              Here’s what to know.
            </span>
          </span>
        </h2>
        <div class="faq-laser-rule w-14 h-[3px] rounded-full bg-[#34D361] mx-auto origin-center" />
        <p class="faq-subtext text-base text-[#5A655C] max-w-lg mx-auto">
          Everything you need to know before you get your free estimate anywhere in the Texas Hill
          Country.
        </p>
      </div>

      <!-- 6 Expandable Accordion Items -->
      <AccordionRoot v-model="openId" type="single" collapsible class="space-y-3">
        <AccordionItem
          v-for="item in faqs"
          :key="item.id"
          :value="item.id"
          class="home-faq-item relative rounded-2xl transition-colors duration-300 overflow-hidden"
          :class="
            openId === item.id
              ? 'bg-[#DCE8D9] border border-[#BFD3BB]'
              : 'bg-white border border-[#E2E6E0] hover:border-[#B6C8B3]'
          "
          @pointermove="handleRowPointerMove"
          @mouseenter="handleRowHover(item.id, true)"
          @mouseleave="handleRowHover(item.id, false)"
        >
          <!-- Active Laser Rail -->
          <span
            :data-faq-rail="item.id"
            aria-hidden="true"
            class="absolute left-0 top-3 bottom-3 w-1 rounded-r-full bg-[#34D361] origin-center pointer-events-none"
          />

          <!-- Cursor-Tracked Radial Sheen -->
          <span
            :data-faq-sheen="item.id"
            aria-hidden="true"
            class="pointer-events-none absolute inset-0 opacity-0"
            style="
              background: radial-gradient(
                260px circle at var(--faq-mx, 50%) var(--faq-my, 50%),
                rgba(52, 211, 97, 0.11),
                transparent 72%
              );
            "
          />

          <AccordionHeader>
            <AccordionTrigger
              class="relative z-10 w-full px-6 py-5 flex items-center justify-between gap-4 text-left cursor-pointer focus:outline-none"
            >
              <span class="text-base md:text-lg font-medium text-[#141A15]">
                {{ item.question }}
              </span>
              <span
                :data-faq-badge="item.id"
                class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 transition-colors duration-200"
                :class="
                  openId === item.id
                    ? 'bg-[#1A241B] text-white'
                    : 'bg-[#F7F7F3] text-[#141A15]'
                "
                aria-hidden="true"
              >
                <Minus v-if="openId === item.id" class="w-4 h-4" />
                <Plus v-else class="w-4 h-4" />
              </span>
            </AccordionTrigger>
          </AccordionHeader>
          <AccordionContent
            force-mount
            :data-faq-panel="item.id"
            class="relative z-10 overflow-hidden"
            :class="[
              !hydrated && openId !== item.id ? 'h-0 opacity-0' : '',
              !motionEnabled && openId !== item.id ? 'hidden' : '',
            ]"
          >
            <div
              :data-faq-copy="item.id"
              class="px-6 pb-5 text-sm md:text-base text-[#364137] leading-relaxed"
            >
              {{ item.answer }}
            </div>
          </AccordionContent>
        </AccordionItem>
      </AccordionRoot>

      <!-- Bottom Link -->
      <div class="mt-8 text-center">
        <RouterLink
          to="/faq"
          class="inline-flex items-center gap-2 text-sm font-medium text-[#141A15] hover:text-[#22C55E] transition-colors"
          @mouseenter="handleLinkHover(true)"
          @mouseleave="handleLinkHover(false)"
        >
          <span>More questions? See all answers and reviews</span>
          <ArrowRight class="faq-bottom-arrow w-4 h-4" />
        </RouterLink>
      </div>
    </div>
  </section>
</template>
