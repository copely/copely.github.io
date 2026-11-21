<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import gsap from 'gsap';
import { ArrowUpRight, ChevronDown, Compass, Menu, Phone, X } from 'lucide-vue-next';
import { useCopelyContent } from '@/composables/useCopelyContent';
import { allowsMotion } from '@/composables/useGsapScope';
import { HOME_SECTIONS, useSiteExplorer } from '@/composables/useSiteExplorer';

const route = useRoute();
const router = useRouter();
const { content } = useCopelyContent();
const {
  activeSection,
  activeSectionId,
  forceExpandTemporarily,
  isExploring,
  jumpToSection,
  navbarMode,
  phasePulseTick,
  scrollProgress,
  setNavbarHovered,
  setNavMenuOpen,
} = useSiteExplorer();

const isDropdownOpen = ref<boolean>(false);
const isExplorerOpen = ref<boolean>(false);
const isMobileOpen = ref<boolean>(false);
let closeTimer: ReturnType<typeof setTimeout> | null = null;

// DOM refs for reversible GSAP navbar morphing
const pillRef = ref<HTMLElement | null>(null);
const navRef = ref<HTMLElement | null>(null);
const pillBackdropRef = ref<HTMLElement | null>(null);
const wordmarkRef = ref<HTMLElement | null>(null);
const telemetryWrapRef = ref<HTMLElement | null>(null);
const centerNavRef = ref<HTMLElement | null>(null);
const rightActionsRef = ref<HTMLElement | null>(null);
const mobileToggleRef = ref<HTMLElement | null>(null);
const laserLineRef = ref<HTMLElement | null>(null);
const ringWrapRef = ref<HTMLElement | null>(null);
const servicesMenuRef = ref<HTMLElement | null>(null);
const explorerMenuRef = ref<HTMLElement | null>(null);
const navPhaseBarRef = ref<HTMLElement | null>(null);
const telemetryReadoutRef = ref<HTMLElement | null>(null);

let compactTl: gsap.core.Timeline | null = null;
let quickLaserScale: ((v: number) => gsap.core.Tween) | null = null;
let prefersReducedMotion = false;

const isLightTopHeader = computed<boolean>(() => {
  return navbarMode.value === 'top' && route.meta.headerTone === 'light';
});

const servicesList = computed(() => content.value.services);

const routeContextLabel = computed<{ code: string; label: string }>(() => {
  if (route.path === '/') {
    return {
      code: activeSection.value.code,
      label: activeSection.value.shortLabel,
    };
  }
  if (route.path.startsWith('/services/')) {
    const slug = String(route.params.slug ?? '');
    const svc = servicesList.value.find((item) => item.slug === slug);
    return { code: 'SVC', label: svc?.title ?? 'Service Specs' };
  }
  if (route.path === '/services') {
    return { code: 'SVC', label: 'All Services' };
  }
  if (route.path === '/about') {
    return { code: 'ABT', label: 'Our Craft' };
  }
  if (route.path.startsWith('/blog')) {
    return { code: 'LOG', label: 'Field Journal' };
  }
  if (route.path === '/faq') {
    return { code: 'FAQ', label: 'Answers & Reviews' };
  }
  if (route.path === '/contact') {
    return { code: 'EST', label: 'Free Estimate' };
  }
  return { code: 'NAV', label: 'Copely' };
});

const ringCircumference = 2 * Math.PI * 24;
const ringDashOffset = computed<number>(() => {
  return ringCircumference * (1 - scrollProgress.value);
});

watch(
  () => isDropdownOpen.value || isExplorerOpen.value || isMobileOpen.value,
  (anyOpen) => {
    setNavMenuOpen(anyOpen);
  },
);

watch(
  () => route.fullPath,
  () => {
    isDropdownOpen.value = false;
    isExplorerOpen.value = false;
    isMobileOpen.value = false;
  },
);

watch(scrollProgress, (val) => {
  if (quickLaserScale) {
    quickLaserScale(val);
  } else if (laserLineRef.value) {
    laserLineRef.value.style.transform = `scaleX(${val})`;
  }
});

watch(navbarMode, (nextMode, prevMode) => {
  applyNavbarModeAnimation(nextMode, prevMode);
});

watch(isDropdownOpen, async (open) => {
  if (!open || prefersReducedMotion) {
    return;
  }
  await nextTick();
  if (servicesMenuRef.value) {
    const items = servicesMenuRef.value.querySelectorAll('.nav-stagger-item');
    gsap.fromTo(
      items,
      { opacity: 0, y: 8 },
      { opacity: 1, y: 0, duration: 0.22, stagger: 0.035, ease: 'power3.out' },
    );
  }
});

watch(isExplorerOpen, async (open) => {
  if (!open || prefersReducedMotion) {
    return;
  }
  await nextTick();
  if (explorerMenuRef.value) {
    const items = explorerMenuRef.value.querySelectorAll('.explorer-stagger-item');
    gsap.fromTo(
      items,
      { opacity: 0, y: 8 },
      { opacity: 1, y: 0, duration: 0.22, stagger: 0.025, ease: 'power3.out' },
    );
  }
});

function buildCompactTimeline(): void {
  if (!pillRef.value) {
    return;
  }

  const collapsibleSides = [
    wordmarkRef.value,
    telemetryWrapRef.value,
    centerNavRef.value,
    rightActionsRef.value,
    mobileToggleRef.value,
  ].filter((el): el is HTMLElement => el !== null);

  compactTl?.kill();
  compactTl = gsap.timeline({
    paused: true,
    onReverseComplete: () => {
      if (navbarMode.value === 'expanded') {
        gsap.set(collapsibleSides, {
          clearProps: 'maxWidth,marginLeft,marginRight,paddingLeft,paddingRight,flexGrow,display',
        });
        if (navRef.value) {
          gsap.set(navRef.value, { clearProps: 'gap,justifyContent' });
        }
        if (pillRef.value) {
          gsap.set(pillRef.value, { clearProps: 'width,height' });
        }
      }
    },
  });

  compactTl
    .to(
      collapsibleSides,
      {
        opacity: 0,
        scale: 0.92,
        duration: 0.14,
        ease: 'power2.in',
      },
      0,
    )
    .to(
      collapsibleSides,
      {
        maxWidth: 0,
        marginLeft: 0,
        marginRight: 0,
        paddingLeft: 0,
        paddingRight: 0,
        flexGrow: 0,
        duration: 0.24,
        ease: 'power3.inOut',
      },
      0.08,
    )
    .set(collapsibleSides, { display: 'none' }, 0.32)
    .to(
      pillRef.value,
      {
        maxWidth: 56,
        height: 56,
        paddingLeft: 16,
        paddingRight: 16,
        paddingTop: 16,
        paddingBottom: 16,
        duration: 0.32,
        ease: 'power3.inOut',
      },
      0.08,
    );

  if (navRef.value) {
    compactTl.to(navRef.value, { gap: 0, justifyContent: 'center', duration: 0.24, ease: 'power3.inOut' }, 0.08);
  }

  if (laserLineRef.value) {
    compactTl.to(laserLineRef.value, { opacity: 0, duration: 0.15 }, 0);
  }
  if (ringWrapRef.value) {
    compactTl.to(ringWrapRef.value, { opacity: 1, scale: 1, duration: 0.24, ease: 'power2.out' }, 0.22);
  }
}

function applyNavbarModeAnimation(
  nextMode: 'top' | 'expanded' | 'compact-logo',
  prevMode?: 'top' | 'expanded' | 'compact-logo',
): void {
  const pill = pillRef.value;
  const backdrop = pillBackdropRef.value;
  if (!pill || !backdrop) {
    return;
  }

  if (prefersReducedMotion) {
    if (nextMode === 'top') {
      compactTl?.progress(0).pause();
      gsap.set(pill, { clearProps: 'all' });
      gsap.set(backdrop, { opacity: 0 });
    } else if (nextMode === 'expanded') {
      compactTl?.progress(0).pause();
      gsap.set(pill, { maxWidth: 1080, marginTop: 16, padding: '12px 22px' });
      gsap.set(backdrop, { opacity: 1 });
    } else {
      compactTl?.progress(1).pause();
      gsap.set(backdrop, { opacity: 1 });
    }
    return;
  }

  if (nextMode === 'top') {
    isExplorerOpen.value = false;
    compactTl?.progress(0).pause();
    gsap.set(
      [
        navRef.value,
        wordmarkRef.value,
        telemetryWrapRef.value,
        centerNavRef.value,
        rightActionsRef.value,
        mobileToggleRef.value,
      ].filter(Boolean),
      { clearProps: 'all' },
    );
    gsap.to(backdrop, { opacity: 0, duration: 0.24, ease: 'power2.out' });
    gsap.to(pill, {
      maxWidth: 1400,
      marginTop: 0,
      paddingLeft: window.innerWidth >= 768 ? 40 : 24,
      paddingRight: window.innerWidth >= 768 ? 40 : 24,
      paddingTop: 24,
      paddingBottom: 24,
      duration: 0.34,
      ease: 'power3.out',
      onComplete: () => {
        if (navbarMode.value === 'top') {
          gsap.set(pill, { clearProps: 'maxWidth,marginTop,paddingLeft,paddingRight,paddingTop,paddingBottom' });
        }
      },
    });
    return;
  }

  if (nextMode === 'expanded') {
    gsap.to(backdrop, { opacity: 1, duration: 0.24, ease: 'power2.out' });
    if (prevMode === 'compact-logo' && compactTl) {
      compactTl.timeScale(1).reverse();
    } else {
      compactTl?.progress(0).pause();
      gsap.to(pill, {
        maxWidth: 1080,
        marginTop: 16,
        paddingLeft: 22,
        paddingRight: 22,
        paddingTop: 11,
        paddingBottom: 11,
        duration: 0.38,
        ease: 'power3.out',
      });
    }
    return;
  }

  // Transition into compact-logo idle emblem
  isDropdownOpen.value = false;
  isExplorerOpen.value = false;
  gsap.set(backdrop, { opacity: 1 });
  if (prevMode === 'top') {
    gsap.set(pill, { maxWidth: 1080, marginTop: 16, padding: '11px 22px' });
  }
  compactTl?.timeScale(1).play();
}

function handleBrandClick(event: MouseEvent): void {
  if (navbarMode.value === 'compact-logo') {
    event.preventDefault();
    forceExpandTemporarily();
    return;
  }
  if (route.path === '/') {
    event.preventDefault();
    jumpToSection('hero', router, route.path);
  }
}

function openDropdown(): void {
  if (closeTimer) {
    clearTimeout(closeTimer);
    closeTimer = null;
  }
  isExplorerOpen.value = false;
  isDropdownOpen.value = true;
}

function scheduleCloseDropdown(): void {
  closeTimer = setTimeout(() => {
    isDropdownOpen.value = false;
  }, 150);
}

function toggleDropdown(): void {
  isExplorerOpen.value = false;
  isDropdownOpen.value = !isDropdownOpen.value;
}

function toggleExplorerPopover(): void {
  isDropdownOpen.value = false;
  isExplorerOpen.value = !isExplorerOpen.value;
}

function handleSectionJump(sectionId: string): void {
  isExplorerOpen.value = false;
  jumpToSection(sectionId, router, route.path);
}

function navigateToEstimate(): void {
  isMobileOpen.value = false;
  isExplorerOpen.value = false;
  isDropdownOpen.value = false;
  jumpToSection('schedule', router, route.path);
}

watch(phasePulseTick, () => {
  if (!allowsMotion() || !navPhaseBarRef.value) {
    return;
  }
  const bar = navPhaseBarRef.value;
  gsap.killTweensOf(bar);
  gsap
    .timeline()
    .fromTo(
      bar,
      { scaleX: 0, opacity: 1, transformOrigin: 'left center' },
      { scaleX: 1, duration: 0.24, ease: 'power3.out' },
    )
    .to(bar, {
      scaleX: 0,
      opacity: 0,
      transformOrigin: 'right center',
      duration: 0.34,
      ease: 'power2.in',
      clearProps: 'transform,opacity,transformOrigin',
    });
});

watch(
  () => `${routeContextLabel.value.code}:${routeContextLabel.value.label}`,
  async () => {
    if (!allowsMotion()) {
      return;
    }
    await nextTick();
    if (!telemetryReadoutRef.value) {
      return;
    }
    gsap.killTweensOf(telemetryReadoutRef.value);
    gsap.fromTo(
      telemetryReadoutRef.value,
      { y: 7, opacity: 0.15 },
      {
        y: 0,
        opacity: 1,
        duration: 0.28,
        ease: 'back.out(1.6)',
        clearProps: 'transform,opacity',
      },
    );
  },
);

onMounted(() => {
  prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  if (laserLineRef.value && !prefersReducedMotion) {
    quickLaserScale = gsap.quickTo(laserLineRef.value, 'scaleX', {
      duration: 0.14,
      ease: 'power2.out',
    });
  }
  buildCompactTimeline();
  applyNavbarModeAnimation(navbarMode.value);
});

onUnmounted(() => {
  compactTl?.kill();
  compactTl = null;
});
</script>

<template>
  <header class="fixed top-0 inset-x-0 z-50 pointer-events-none px-3 md:px-6">
    <!-- Architectural Viewport Phase Conduit Streak -->
    <div
      ref="navPhaseBarRef"
      aria-hidden="true"
      class="pointer-events-none fixed top-0 inset-x-0 h-[3px] z-[60] bg-gradient-to-r from-[#22C55E] via-[#4ADE80] to-[#34D361] origin-left scale-x-0 opacity-0"
    />
    <div
      ref="pillRef"
      class="pointer-events-auto relative max-w-[1400px] mx-auto px-6 md:px-10 py-6"
      :class="navbarMode === 'compact-logo' ? 'cursor-pointer' : ''"
      @mouseenter="setNavbarHovered(true)"
      @mouseleave="setNavbarHovered(false)"
      @focusin="setNavbarHovered(true)"
      @focusout="setNavbarHovered(false)"
    >
      <!-- Frosted Explorer Pill Surface & Laser Progress Hairline -->
      <div
        ref="pillBackdropRef"
        class="pointer-events-none absolute inset-0 rounded-full bg-[#141A15]/95 backdrop-blur-xl border border-[#4ADE80]/30 shadow-[0_18px_50px_rgba(20,26,21,0.42)] overflow-hidden opacity-0"
        aria-hidden="true"
      >
        <!-- Bottom Laser-Level Scroll Progress Hairline -->
        <div
          ref="laserLineRef"
          class="absolute bottom-0 inset-x-0 h-[2px] bg-gradient-to-r from-[#22C55E] via-[#4ADE80] to-[#34D361] origin-left scale-x-0"
        />
      </div>

      <!-- Compact Logo-Only Circular Progress Ring -->
      <svg
        ref="ringWrapRef"
        viewBox="0 0 56 56"
        class="pointer-events-none absolute inset-0 w-full h-full opacity-0 scale-90"
        aria-hidden="true"
      >
        <circle
          cx="28"
          cy="28"
          r="24"
          fill="none"
          stroke="#4ADE80"
          stroke-opacity="0.18"
          stroke-width="2"
        />
        <circle
          cx="28"
          cy="28"
          r="24"
          fill="none"
          stroke="#4ADE80"
          stroke-width="2"
          stroke-linecap="round"
          :stroke-dasharray="ringCircumference"
          :stroke-dashoffset="ringDashOffset"
          transform="rotate(-90 28 28)"
        />
      </svg>

      <nav ref="navRef" class="relative z-10 flex items-center justify-between gap-3" aria-label="Main navigation">
        <!-- Left: Brand Mark + Live Exploration Telemetry Chip -->
        <div class="flex items-center shrink-0">
          <RouterLink
            to="/"
            class="inline-flex items-center select-none focus:outline-none focus-visible:ring-2 focus-visible:ring-[#4ADE80] rounded-lg"
            :aria-label="navbarMode === 'compact-logo' ? 'Expand Copely navigation' : 'Copely home'"
            @click="handleBrandClick"
          >
            <svg
              width="24"
              height="24"
              viewBox="0 0 28 28"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              aria-hidden="true"
              class="shrink-0"
            >
              <path
                d="M3.5 19.5L14 9L24.5 19.5L21.5 22.5L14 15L6.5 22.5L3.5 19.5Z"
                fill="#4ADE80"
              />
              <path
                d="M7.5 11.5L14 5L20.5 11.5L17.8 14.2L14 10.4L10.2 14.2L7.5 11.5Z"
                fill="#22C55E"
              />
            </svg>
            <span
              ref="wordmarkRef"
              class="overflow-hidden whitespace-nowrap font-sans font-semibold tracking-tight text-xl leading-none ml-2.5 transition-colors"
              :class="isLightTopHeader ? 'text-[#141A15]' : 'text-white'"
            >
              Copely
            </span>
          </RouterLink>

          <!-- Live Exploration Section Telemetry Chip (visible when exploring) -->
          <div
            ref="telemetryWrapRef"
            class="hidden md:flex items-center overflow-hidden whitespace-nowrap"
            :class="isExploring ? 'max-w-[230px] opacity-100 ml-3 pl-3 border-l border-white/15' : 'max-w-0 opacity-0 ml-0 pl-0 border-l-0'"
          >
            <div class="relative">
              <button
                type="button"
                class="inline-flex items-center gap-2 rounded-full bg-white/10 hover:bg-white/15 border border-white/10 px-3 py-1 text-xs font-medium text-white cursor-pointer"
                :aria-expanded="isExplorerOpen"
                @click="toggleExplorerPopover"
              >
                <span class="w-1.5 h-1.5 rounded-full bg-[#4ADE80] animate-pulse" aria-hidden="true" />
                <span ref="telemetryReadoutRef" class="inline-flex items-center gap-2 overflow-hidden">
                  <span class="text-[#4ADE80] font-semibold">{{ routeContextLabel.code }}</span>
                  <span class="text-white/85 truncate max-w-[105px]">
                    {{ routeContextLabel.label }}
                  </span>
                </span>
                <ChevronDown
                  class="w-3.5 h-3.5 text-white/60 transition-transform duration-150"
                  :class="isExplorerOpen ? 'rotate-180' : ''"
                />
              </button>
            </div>
          </div>
        </div>

        <!-- Center Navigation Links (Desktop) -->
        <div
          ref="centerNavRef"
          class="hidden lg:flex items-center gap-6 text-sm font-medium overflow-visible whitespace-nowrap transition-colors"
          :class="isLightTopHeader ? 'text-[#141A15]/85' : 'text-white/90'"
        >
          <!-- Services Dropdown -->
          <div
            class="relative"
            @mouseenter="openDropdown"
            @mouseleave="scheduleCloseDropdown"
          >
            <button
              type="button"
              class="inline-flex items-center gap-1 py-1.5 transition-colors hover:text-[#4ADE80] focus:outline-none cursor-pointer"
              :aria-expanded="isDropdownOpen"
              aria-haspopup="true"
              @click="toggleDropdown"
            >
              <span>Services</span>
              <ChevronDown
                class="w-4 h-4 transition-transform duration-200"
                :class="isDropdownOpen ? 'rotate-180' : ''"
              />
            </button>

            <div
              v-if="isDropdownOpen"
              ref="servicesMenuRef"
              class="absolute left-0 top-full mt-3 w-72 rounded-2xl bg-[#1A241B] text-white shadow-2xl border border-white/15 p-2"
            >
              <RouterLink
                v-for="svc in servicesList"
                :key="svc.id"
                :to="`/services/${svc.slug}`"
                class="nav-stagger-item flex flex-col rounded-xl px-3.5 py-2.5 transition-colors hover:bg-white/10"
              >
                <span class="text-sm font-medium text-white">{{ svc.title }}</span>
                <span class="text-xs text-white/60">{{ svc.shortSubtitle }}</span>
              </RouterLink>
              <div class="my-1.5 border-t border-white/10" />
              <RouterLink
                to="/services"
                class="nav-stagger-item flex items-center justify-between rounded-xl px-3.5 py-2 text-xs font-medium text-[#4ADE80] hover:bg-white/5"
              >
                <span>View all tiling services</span>
                <span aria-hidden="true">→</span>
              </RouterLink>
            </div>
          </div>

          <RouterLink
            to="/about"
            class="py-1.5 transition-colors hover:text-[#4ADE80]"
            active-class="text-[#4ADE80]"
          >
            About
          </RouterLink>
          <RouterLink
            to="/blog"
            class="py-1.5 transition-colors hover:text-[#4ADE80]"
            active-class="text-[#4ADE80]"
          >
            Blog
          </RouterLink>
          <RouterLink
            to="/faq"
            class="py-1.5 transition-colors hover:text-[#4ADE80]"
            active-class="text-[#4ADE80]"
          >
            FAQ
          </RouterLink>
          <RouterLink
            to="/contact"
            class="py-1.5 transition-colors hover:text-[#4ADE80]"
            active-class="text-[#4ADE80]"
          >
            Contact
          </RouterLink>
        </div>

        <!-- Right Actions -->
        <div ref="rightActionsRef" class="hidden md:flex items-center gap-4 shrink-0 whitespace-nowrap">
          <a
            href="tel:+15125550142"
            class="inline-flex items-center gap-1.5 text-sm font-medium transition-colors hover:text-[#4ADE80]"
            :class="isLightTopHeader ? 'text-[#141A15]' : 'text-white/95'"
          >
            <Phone class="w-4 h-4 text-[#4ADE80]" />
            <span>(512) 555-0142</span>
          </a>

          <button
            type="button"
            class="rounded-full px-5 py-2.5 text-sm font-semibold transition-colors duration-200 shadow-sm cursor-pointer"
            :class="
              isExploring
                ? 'bg-[#34D361] text-[#141A15] hover:bg-[#4ADE80]'
                : isLightTopHeader
                  ? 'bg-[#1A241B] text-white hover:bg-[#273528]'
                  : 'bg-white text-[#141A15] hover:bg-[#E8EFE6]'
            "
            @click="navigateToEstimate"
          >
            Free Estimate
          </button>
        </div>

        <!-- Mobile Menu Toggle -->
        <button
          ref="mobileToggleRef"
          type="button"
          class="lg:hidden inline-flex items-center justify-center w-10 h-10 rounded-full shrink-0 transition-colors"
          :class="
            isLightTopHeader
              ? 'bg-[#1A241B]/10 text-[#141A15]'
              : 'bg-white/10 text-white'
          "
          aria-label="Toggle navigation menu"
          :aria-expanded="isMobileOpen"
          @click="isMobileOpen = !isMobileOpen"
        >
          <X v-if="isMobileOpen" class="w-5 h-5" />
          <Menu v-else class="w-5 h-5" />
        </button>
      </nav>

      <!-- Site Explorer Quick-Jump Popover -->
      <div
        v-if="isExplorerOpen && navbarMode === 'expanded'"
        ref="explorerMenuRef"
        class="hidden md:block absolute left-6 top-full mt-3 w-[340px] rounded-2xl bg-[#141A15]/95 backdrop-blur-xl text-white border border-[#4ADE80]/25 shadow-2xl p-3"
      >
        <div class="flex items-center justify-between px-2 pb-2 mb-2 border-b border-white/10">
          <span class="inline-flex items-center gap-1.5 text-[11px] font-semibold uppercase tracking-[0.16em] text-[#4ADE80]">
            <Compass class="w-3.5 h-3.5" />
            Site Explorer
          </span>
          <span class="text-[11px] text-white/50">
            {{ Math.round(scrollProgress * 100) }}% explored
          </span>
        </div>

        <div class="grid grid-cols-1 gap-1 max-h-[min(420px,68vh)] overflow-y-auto pr-1">
          <button
            v-for="sec in HOME_SECTIONS"
            :key="sec.id"
            type="button"
            class="explorer-stagger-item w-full text-left flex items-center justify-between rounded-xl px-3 py-2 text-xs transition-colors cursor-pointer"
            :class="
              route.path === '/' && activeSectionId === sec.id
                ? 'bg-[#34D361]/15 text-[#4ADE80]'
                : 'text-white/80 hover:bg-white/10 hover:text-white'
            "
            @click="handleSectionJump(sec.id)"
          >
            <span class="flex items-center gap-2.5">
              <span class="font-mono text-[11px] opacity-65">{{ sec.code }}</span>
              <span class="font-medium">{{ sec.title }}</span>
            </span>
            <ArrowUpRight class="w-3.5 h-3.5 opacity-60" />
          </button>
        </div>
      </div>
    </div>

    <!-- Mobile Navigation Drawer -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-if="isMobileOpen"
        class="pointer-events-auto lg:hidden mx-1 mt-2 rounded-3xl bg-[#1A241B] text-white border border-white/10 shadow-2xl p-6 space-y-5"
      >
        <div>
          <p class="text-[11px] uppercase tracking-wider text-white/50 mb-2">Services</p>
          <div class="grid grid-cols-1 gap-1.5">
            <RouterLink
              v-for="svc in servicesList"
              :key="svc.id"
              :to="`/services/${svc.slug}`"
              class="rounded-xl px-3 py-2 text-sm text-white/90 hover:bg-white/10"
            >
              {{ svc.title }}
            </RouterLink>
          </div>
        </div>

        <div class="border-t border-white/10 pt-4 flex flex-col gap-2 text-sm font-medium">
          <RouterLink to="/about" class="px-3 py-2 rounded-xl hover:bg-white/10">About</RouterLink>
          <RouterLink to="/blog" class="px-3 py-2 rounded-xl hover:bg-white/10">Blog</RouterLink>
          <RouterLink to="/faq" class="px-3 py-2 rounded-xl hover:bg-white/10">FAQ</RouterLink>
          <RouterLink to="/contact" class="px-3 py-2 rounded-xl hover:bg-white/10">Contact</RouterLink>
        </div>

        <div class="border-t border-white/10 pt-4 flex flex-col gap-3">
          <a
            href="tel:+15125550142"
            class="inline-flex items-center justify-center gap-2 rounded-full border border-white/20 py-2.5 text-sm font-medium text-white"
          >
            <Phone class="w-4 h-4 text-[#4ADE80]" />
            <span>(512) 555-0142</span>
          </a>
          <button
            type="button"
            class="w-full rounded-full bg-[#34D361] text-[#141A15] py-3 text-sm font-semibold"
            @click="navigateToEstimate"
          >
            Free Estimate
          </button>
        </div>
      </div>
    </Transition>
  </header>
</template>
