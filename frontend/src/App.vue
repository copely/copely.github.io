<script setup lang="ts">
import { ref } from 'vue';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import Navbar from '@/components/layout/Navbar.vue';
import Footer from '@/components/layout/Footer.vue';
import { allowsMotion } from '@/composables/useGsapScope';
import { instantScrollTo, useSiteExplorer } from '@/composables/useSiteExplorer';

const mainRef = ref<HTMLElement | null>(null);
const footerStageRef = ref<HTMLElement | null>(null);

const {
  activeSectionId,
  consumePendingRouteScroll,
  forceExpandTemporarily,
  pulseSectionFocus,
  setProgrammaticScrollLock,
  syncScrollNow,
  triggerPhasePulse,
} = useSiteExplorer();

function onRouteBeforeLeave(): void {
  if (mainRef.value) {
    mainRef.value.style.minHeight = `${mainRef.value.offsetHeight}px`;
  }
  if (footerStageRef.value && allowsMotion()) {
    gsap.killTweensOf(footerStageRef.value);
    gsap.to(footerStageRef.value, {
      opacity: 0.32,
      duration: 0.18,
      ease: 'power2.out',
    });
  }
}

function onRouteLeave(el: Element, done: () => void): void {
  triggerPhasePulse();
  if (!allowsMotion()) {
    done();
    return;
  }
  gsap.killTweensOf(el);
  gsap.to(el, {
    opacity: 0,
    y: -16,
    scale: 0.988,
    duration: 0.21,
    ease: 'power2.in',
    onComplete: done,
  });
}

function onRouteBeforeEnter(el: Element): void {
  if (allowsMotion()) {
    gsap.set(el, { opacity: 0 });
  }
}

function onRouteEnter(el: Element, done: () => void): void {
  if (mainRef.value) {
    mainRef.value.style.minHeight = '';
  }

  const pending = consumePendingRouteScroll();
  let focusSectionId: string | null = null;

  if (pending?.hash) {
    const sectionId = pending.hash.replace(/^#/, '');
    const target = document.getElementById(sectionId);
    if (target) {
      const nextTop =
        sectionId === 'hero'
          ? 0
          : Math.max(0, target.getBoundingClientRect().top + window.scrollY - 88);
      setProgrammaticScrollLock(450);
      instantScrollTo(nextTop, 0);
      syncScrollNow();
      if (sectionId !== 'hero' && nextTop > 40) {
        forceExpandTemporarily();
      }
      ScrollTrigger.refresh();
      activeSectionId.value = sectionId;
      focusSectionId = sectionId;
    }
  } else if (pending) {
    setProgrammaticScrollLock(280);
    instantScrollTo(pending.top ?? 0, pending.left ?? 0);
    syncScrollNow();
    ScrollTrigger.refresh();
  }

  if (!allowsMotion()) {
    gsap.set(el, { clearProps: 'opacity' });
    if (footerStageRef.value) {
      gsap.set(footerStageRef.value, { clearProps: 'opacity' });
    }
    done();
    return;
  }

  if (footerStageRef.value) {
    gsap.killTweensOf(footerStageRef.value);
    gsap.to(footerStageRef.value, {
      opacity: 1,
      duration: 0.36,
      ease: 'power2.out',
      clearProps: 'opacity',
    });
  }

  gsap.killTweensOf(el);
  gsap.fromTo(
    el,
    { opacity: 0, y: 22, scale: 0.988 },
    {
      opacity: 1,
      y: 0,
      scale: 1,
      duration: 0.42,
      ease: 'expo.out',
      clearProps: 'transform,opacity',
      onStart: () => {
        if (focusSectionId) {
          pulseSectionFocus(document.getElementById(focusSectionId));
        }
      },
      onComplete: () => {
        ScrollTrigger.refresh();
        if (focusSectionId) {
          activeSectionId.value = focusSectionId;
        }
        done();
      },
    },
  );
}

function onRouteAfterEnter(): void {
  if (mainRef.value) {
    mainRef.value.style.minHeight = '';
  }
}
</script>

<template>
  <div class="min-h-screen flex flex-col bg-[#F7F7F3] text-[#141A15]">
    <Navbar />
    <div id="site-phase-stage" class="flex-1 flex flex-col">
      <main ref="mainRef" class="flex-1 min-h-screen">
        <RouterView v-slot="{ Component, route }">
          <Transition
            mode="out-in"
            :css="false"
            @before-leave="onRouteBeforeLeave"
            @leave="onRouteLeave"
            @before-enter="onRouteBeforeEnter"
            @enter="onRouteEnter"
            @after-enter="onRouteAfterEnter"
          >
            <div :key="route.path === '/' ? 'home' : route.path" class="route-phase-shell">
              <component :is="Component" />
            </div>
          </Transition>
        </RouterView>
      </main>
      <div ref="footerStageRef">
        <Footer />
      </div>
    </div>
  </div>
</template>
