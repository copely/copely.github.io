import { computed, ref } from 'vue';
import type { Router } from 'vue-router';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import { allowsMotion, ensureGsapRegistered } from '@/composables/useGsapScope';

export interface ExplorerSectionMeta {
  id: string;
  code: string;
  shortLabel: string;
  title: string;
}

export type NavbarMode = 'top' | 'expanded' | 'compact-logo';

export const HOME_SECTIONS: ExplorerSectionMeta[] = [
  { id: 'hero', code: '01', shortLabel: 'Overview', title: 'Copely Tiling Hero' },
  { id: 'system', code: '02', shortLabel: 'Tile System', title: '5-Layer Engineered Assembly' },
  { id: 'services', code: '03', shortLabel: 'Services', title: 'Interior & Patio Services' },
  { id: 'recent-work', code: '04', shortLabel: 'Recent Work', title: 'Hill Country Projects' },
  { id: 'process', code: '05', shortLabel: 'Process', title: 'Five Steps Start to Finish' },
  { id: 'before-after', code: '06', shortLabel: 'Transformations', title: 'Before & After Sliders' },
  { id: 'testimonials', code: '07', shortLabel: 'Stories', title: 'Homeowner Case Studies' },
  { id: 'why-copely', code: '08', shortLabel: 'Why Copely', title: 'Compare & 25-Yr Warranty' },
  { id: 'faq', code: '09', shortLabel: 'FAQ', title: 'Common Questions' },
  { id: 'schedule', code: '10', shortLabel: 'Book Visit', title: 'Free On-Site Estimate' },
];

const scrollY = ref<number>(0);
const scrollProgress = ref<number>(0);
const scrollDirection = ref<'up' | 'down'>('up');
const isIdleDeep = ref<boolean>(false);
const isNavbarHovered = ref<boolean>(false);
const isNavMenuOpen = ref<boolean>(false);
const activeSectionId = ref<string>('hero');
const phasePulseTick = ref<number>(0);

export interface PendingRouteScroll {
  top?: number;
  left?: number;
  hash?: string;
}

let listenerAttached = false;
let lastY = 0;
let rafPending = false;
let idleTimer: ReturnType<typeof setTimeout> | null = null;
let sectionTriggers: ScrollTrigger[] = [];
let activePhaseTl: gsap.core.Timeline | null = null;
let programmaticScrollLockUntil = 0;
let pendingRouteScroll: PendingRouteScroll | null = null;

function setProgrammaticScrollLock(durationMs = 320): void {
  if (typeof performance !== 'undefined') {
    programmaticScrollLockUntil = performance.now() + durationMs;
  }
}

export function instantScrollTo(top: number, left = 0): void {
  if (typeof window === 'undefined') {
    return;
  }
  const htmlEl = document.documentElement;
  const prevBehavior = htmlEl.style.scrollBehavior;
  htmlEl.style.scrollBehavior = 'auto';
  window.scrollTo({ top, left, behavior: 'instant' });
  htmlEl.style.scrollBehavior = prevBehavior;
}

function isProgrammaticScrollLocked(): boolean {
  return typeof performance !== 'undefined' && performance.now() < programmaticScrollLockUntil;
}

function clearIdleTimer(): void {
  if (idleTimer) {
    clearTimeout(idleTimer);
    idleTimer = null;
  }
}

function scheduleIdleCollapse(): void {
  clearIdleTimer();
  if (scrollY.value <= 140) {
    isIdleDeep.value = false;
    return;
  }
  idleTimer = setTimeout(() => {
    if (scrollY.value > 140 && !isNavbarHovered.value && !isNavMenuOpen.value) {
      isIdleDeep.value = true;
    }
  }, 2800);
}

function updateScrollMetrics(): void {
  rafPending = false;
  const currentY = Math.max(0, window.scrollY || window.pageYOffset || 0);
  const delta = currentY - lastY;

  if (isProgrammaticScrollLocked()) {
    lastY = currentY;
    scrollDirection.value = 'up';
    isIdleDeep.value = false;
  } else if (Math.abs(delta) >= 6) {
    scrollDirection.value = delta > 0 ? 'down' : 'up';
    lastY = currentY;
    isIdleDeep.value = false;
  }

  scrollY.value = currentY;
  const docHeight = document.documentElement.scrollHeight - window.innerHeight;
  scrollProgress.value = docHeight > 0 ? Math.min(1, Math.max(0, currentY / docHeight)) : 0;

  if (currentY <= 40) {
    isIdleDeep.value = false;
    clearIdleTimer();
  } else {
    scheduleIdleCollapse();
  }
}

function attachWindowScrollListener(): void {
  if (listenerAttached || typeof window === 'undefined') {
    return;
  }
  listenerAttached = true;
  lastY = window.scrollY || 0;
  updateScrollMetrics();

  window.addEventListener(
    'scroll',
    () => {
      if (!rafPending) {
        rafPending = true;
        window.requestAnimationFrame(updateScrollMetrics);
      }
    },
    { passive: true },
  );
}

export function useSiteExplorer() {
  attachWindowScrollListener();

  const isExploring = computed<boolean>(() => scrollY.value > 40);

  const navbarMode = computed<NavbarMode>(() => {
    if (scrollY.value <= 40) {
      return 'top';
    }
    if (isNavbarHovered.value || isNavMenuOpen.value) {
      return 'expanded';
    }
    if (scrollY.value > 140 && (scrollDirection.value === 'down' || isIdleDeep.value)) {
      return 'compact-logo';
    }
    return 'expanded';
  });

  const activeSection = computed<ExplorerSectionMeta>(() => {
    return HOME_SECTIONS.find((sec) => sec.id === activeSectionId.value) ?? HOME_SECTIONS[0]!;
  });

  function setNavbarHovered(hovered: boolean): void {
    isNavbarHovered.value = hovered;
    if (hovered) {
      isIdleDeep.value = false;
      clearIdleTimer();
    } else {
      scheduleIdleCollapse();
    }
  }

  function setNavMenuOpen(open: boolean): void {
    isNavMenuOpen.value = open;
    if (open) {
      isIdleDeep.value = false;
      clearIdleTimer();
    } else {
      scheduleIdleCollapse();
    }
  }

  function forceExpandTemporarily(): void {
    isIdleDeep.value = false;
    scrollDirection.value = 'up';
    scheduleIdleCollapse();
  }

  function registerHomeSections(): void {
    ensureGsapRegistered();
    unregisterHomeSections();

    for (const section of HOME_SECTIONS) {
      const el = document.getElementById(section.id);
      if (!el) {
        continue;
      }
      const st = ScrollTrigger.create({
        trigger: el,
        start: 'top 42%',
        end: 'bottom 42%',
        onEnter: () => {
          activeSectionId.value = section.id;
        },
        onEnterBack: () => {
          activeSectionId.value = section.id;
        },
      });
      sectionTriggers.push(st);
    }
    ScrollTrigger.refresh();
  }

  function unregisterHomeSections(): void {
    activePhaseTl?.kill();
    activePhaseTl = null;
    const stage = typeof document !== 'undefined' ? document.getElementById('site-phase-stage') : null;
    if (stage) {
      gsap.set(stage, { clearProps: 'transform,opacity' });
    }
    for (const st of sectionTriggers) {
      st.kill();
    }
    sectionTriggers = [];
  }

  function triggerPhasePulse(): void {
    phasePulseTick.value += 1;
  }

  function pulseSectionFocus(target: HTMLElement | null): void {
    if (!target || !allowsMotion()) {
      return;
    }
    gsap.fromTo(
      target,
      { boxShadow: 'inset 0 2px 0 0 rgba(52, 211, 97, 0.55)' },
      {
        boxShadow: 'inset 0 2px 0 0 rgba(52, 211, 97, 0)',
        duration: 0.85,
        ease: 'power2.out',
        clearProps: 'boxShadow',
      },
    );
  }

  function setPendingRouteScroll(scroll: PendingRouteScroll | null): void {
    pendingRouteScroll = scroll;
  }

  function consumePendingRouteScroll(): PendingRouteScroll | null {
    const val = pendingRouteScroll;
    pendingRouteScroll = null;
    return val;
  }

  function jumpToSection(sectionId: string, router: Router, currentPath: string): void {
    isNavMenuOpen.value = false;
    if (currentPath !== '/') {
      triggerPhasePulse();
      void router.push({ path: '/', hash: `#${sectionId}` });
      return;
    }

    const target = document.getElementById(sectionId);
    if (!target) {
      void router.push({ path: '/', hash: `#${sectionId}` });
      return;
    }

    const desiredOffset = sectionId === 'hero' ? 0 : 88;
    const rawTop =
      sectionId === 'hero' ? 0 : Math.max(0, target.getBoundingClientRect().top + window.scrollY - desiredOffset);

    if (!allowsMotion()) {
      setProgrammaticScrollLock(250);
      instantScrollTo(rawTop, 0);
      lastY = rawTop;
      updateScrollMetrics();
      ScrollTrigger.refresh();
      activeSectionId.value = sectionId;
      if (sectionId !== 'hero' && rawTop > 40) {
        forceExpandTemporarily();
      }
      return;
    }

    const stage = document.getElementById('site-phase-stage');
    const distFromCurrent = Math.abs(window.scrollY - rawTop);
    if (!stage || distFromCurrent < 90) {
      activeSectionId.value = sectionId;
      triggerPhasePulse();
      pulseSectionFocus(target);
      return;
    }

    const dir = rawTop >= window.scrollY ? 1 : -1;
    triggerPhasePulse();

    activePhaseTl?.kill();
    activePhaseTl = gsap.timeline({
      onComplete: () => {
        activePhaseTl = null;
      },
    });

    activePhaseTl
      .to(stage, {
        opacity: 0,
        y: -dir * 18,
        scale: 0.986,
        duration: 0.21,
        ease: 'power2.in',
      })
      .call(() => {
        gsap.set(stage, { y: 0, scale: 1 });
        const nextTop =
          sectionId === 'hero'
            ? 0
            : Math.max(0, target.getBoundingClientRect().top + window.scrollY - desiredOffset);
        setProgrammaticScrollLock(450);
        instantScrollTo(nextTop, 0);
        lastY = nextTop;
        updateScrollMetrics();
        if (sectionId !== 'hero' && nextTop > 40) {
          forceExpandTemporarily();
        }
        ScrollTrigger.refresh();
        activeSectionId.value = sectionId;
      })
      .fromTo(
        stage,
        { opacity: 0, y: dir * 24, scale: 0.986 },
        {
          opacity: 1,
          y: 0,
          scale: 1,
          duration: 0.42,
          ease: 'expo.out',
          clearProps: 'transform,opacity',
          onStart: () => {
            pulseSectionFocus(target);
          },
        },
      );
  }

  return {
    scrollY,
    scrollProgress,
    scrollDirection,
    isExploring,
    navbarMode,
    activeSectionId,
    activeSection,
    phasePulseTick,
    triggerPhasePulse,
    pulseSectionFocus,
    setPendingRouteScroll,
    consumePendingRouteScroll,
    setProgrammaticScrollLock,
    syncScrollNow: updateScrollMetrics,
    setNavbarHovered,
    setNavMenuOpen,
    forceExpandTemporarily,
    registerHomeSections,
    unregisterHomeSections,
    jumpToSection,
  };
}
