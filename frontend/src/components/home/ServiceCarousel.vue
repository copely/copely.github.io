<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue';
import { ArrowLeft, ArrowRight, ArrowUpRight, Pause, Play } from 'lucide-vue-next';
import { useGsapScope } from '@/composables/useGsapScope';
import type { ServiceItem } from '@/types/copely';

interface Props {
  services: ServiceItem[];
}

interface RenderedSlide {
  key: string;
  service: ServiceItem;
  serviceIdx: number;
  isPrimary: boolean;
}

const AUTOPLAY_INTERVAL_MS = 4500;
const SLIDE_DURATION_MS = 520;
const DRAG_THRESHOLD_PX = 6;

const props = defineProps<Props>();

const sectionRef = ref<HTMLElement | null>(null);
const viewportRef = ref<HTMLElement | null>(null);
const trackRef = ref<HTMLElement | null>(null);

const activeIndex = ref<number>(0);
const stepWidth = ref<number>(306);
const isDragging = ref<boolean>(false);
const isHovered = ref<boolean>(false);
const isFocusedWithin = ref<boolean>(false);
const isInViewport = ref<boolean>(true);
const isUserPaused = ref<boolean>(false);
const prefersReducedMotion = ref<boolean>(false);

const count = computed<number>(() => props.services.length);
const isLooping = computed<boolean>(() => count.value > 1);
const copiesEachSide = computed<number>(() => {
  if (!isLooping.value) {
    return 0;
  }
  return Math.max(2, Math.ceil(1600 / (count.value * 240)));
});

const renderedSlides = computed<RenderedSlide[]>(() => {
  if (count.value === 0) {
    return [];
  }
  const totalCopies = copiesEachSide.value * 2 + 1;
  const primaryCopy = copiesEachSide.value;
  const slides: RenderedSlide[] = [];

  for (let copyIdx = 0; copyIdx < totalCopies; copyIdx += 1) {
    for (let serviceIdx = 0; serviceIdx < count.value; serviceIdx += 1) {
      const service = props.services[serviceIdx];
      slides.push({
        key: `${copyIdx}-${service.id}`,
        service,
        serviceIdx,
        isPrimary: copyIdx === primaryCopy,
      });
    }
  }
  return slides;
});

const isAutoPlayActive = computed<boolean>(() => {
  return (
    isLooping.value &&
    !isUserPaused.value &&
    !isHovered.value &&
    !isDragging.value &&
    !isFocusedWithin.value &&
    isInViewport.value &&
    !prefersReducedMotion.value
  );
});

const formattedSlideNumber = computed<string>(() => {
  if (count.value === 0) {
    return '00';
  }
  return String(activeIndex.value + 1).padStart(2, '0');
});

const formattedTotalSlides = computed<string>(() => {
  return String(count.value).padStart(2, '0');
});

// Non-reactive 60fps animation/gesture state
let currentX = 0;
let startX = 0;
let targetX = 0;
let animStartTime = 0;
let rafId: number | null = null;
let autoplayTimer: number | null = null;
let wheelSettleTimer: number | null = null;

let dragStartPointerX = 0;
let dragStartPointerY = 0;
let dragStartTrackX = 0;
let dragMoved = false;
let dragDirectionLocked: 'horizontal' | 'vertical' | null = null;
let suppressClick = false;
let suppressClickTimer: number | null = null;
let lastPointerX = 0;
let lastPointerTime = 0;
let velocityX = 0;

let resizeObserver: ResizeObserver | null = null;
let intersectionObserver: IntersectionObserver | null = null;
let motionQuery: MediaQueryList | null = null;

function getCenterStart(step = stepWidth.value): number {
  return -copiesEachSide.value * count.value * step;
}

function getLoopWidth(step = stepWidth.value): number {
  return count.value * step;
}

function wrapOffset(x: number): number {
  if (!isLooping.value) {
    return 0;
  }
  const loopWidth = getLoopWidth();
  if (loopWidth <= 0) {
    return 0;
  }
  const centerStart = getCenterStart();
  const offsetFromStart = centerStart - x;
  const wrapped = ((offsetFromStart % loopWidth) + loopWidth) % loopWidth;
  return centerStart - wrapped;
}

function getLogicalIndex(x: number): number {
  const n = count.value;
  const step = stepWidth.value;
  if (n <= 1 || step <= 0) {
    return 0;
  }
  const rawStep = Math.round((getCenterStart() - x) / step);
  return ((rawStep % n) + n) % n;
}

function applyTransform(): void {
  const track = trackRef.value;
  if (!track) {
    return;
  }
  const displayX = isLooping.value ? wrapOffset(currentX) : 0;
  track.style.transform = `translate3d(${displayX}px, 0, 0)`;
  activeIndex.value = getLogicalIndex(displayX);
}

function cancelAnimation(): void {
  if (rafId !== null) {
    cancelAnimationFrame(rafId);
    rafId = null;
  }
}

function animateTo(nextTargetX: number, immediate = false): void {
  if (!isLooping.value) {
    return;
  }
  targetX = nextTargetX;

  if (immediate || prefersReducedMotion.value) {
    cancelAnimation();
    currentX = wrapOffset(targetX);
    targetX = currentX;
    startX = currentX;
    applyTransform();
    return;
  }

  startX = currentX;
  animStartTime = performance.now();

  if (rafId !== null) {
    return;
  }

  const tick = (now: number) => {
    const elapsed = now - animStartTime;
    const progress = Math.min(1, elapsed / SLIDE_DURATION_MS);
    const eased = 1 - Math.pow(1 - progress, 3);

    currentX = startX + (targetX - startX) * eased;
    applyTransform();

    if (progress < 1) {
      rafId = requestAnimationFrame(tick);
    } else {
      rafId = null;
      currentX = wrapOffset(targetX);
      targetX = currentX;
      startX = currentX;
      applyTransform();
    }
  };

  rafId = requestAnimationFrame(tick);
}

function resetToIndex(index: number): void {
  const clamped = count.value > 0 ? ((index % count.value) + count.value) % count.value : 0;
  const x = isLooping.value ? getCenterStart() - clamped * stepWidth.value : 0;
  cancelAnimation();
  currentX = x;
  targetX = x;
  startX = x;
  applyTransform();
}

function stepSlides(deltaSlides: number, fromUser = true): void {
  if (!isLooping.value || stepWidth.value <= 0) {
    return;
  }
  if (rafId === null) {
    currentX = wrapOffset(currentX);
    targetX = currentX;
  }
  const nextTarget = targetX - deltaSlides * stepWidth.value;
  animateTo(nextTarget);
  if (fromUser) {
    restartAutoplay();
  }
}

function goToService(targetLogical: number, fromUser = true): void {
  if (!isLooping.value || stepWidth.value <= 0) {
    return;
  }
  const n = count.value;
  if (rafId === null) {
    currentX = wrapOffset(currentX);
    targetX = currentX;
  }
  const curLogical = getLogicalIndex(targetX);
  let delta = targetLogical - curLogical;
  if (delta > n / 2) {
    delta -= n;
  }
  if (delta < -n / 2) {
    delta += n;
  }
  if (delta === 0) {
    snapToNearest();
    if (fromUser) {
      restartAutoplay();
    }
    return;
  }
  animateTo(targetX - delta * stepWidth.value);
  if (fromUser) {
    restartAutoplay();
  }
}

function snapToNearest(projectedOffset = currentX): void {
  if (!isLooping.value || stepWidth.value <= 0) {
    return;
  }
  const step = stepWidth.value;
  const centerStart = getCenterStart();
  const nearestStep = Math.round((centerStart - projectedOffset) / step);
  animateTo(centerStart - nearestStep * step);
}

function ensureCardVisibleOnFocus(serviceIdx: number): void {
  if (!isLooping.value || !viewportRef.value || stepWidth.value <= 0) {
    return;
  }
  const n = count.value;
  const visibleSlots = Math.max(1, Math.floor(viewportRef.value.clientWidth / stepWidth.value));
  const cur = activeIndex.value;
  const offsetFromActive = ((serviceIdx - cur) % n + n) % n;
  if (offsetFromActive >= visibleSlots) {
    goToService(serviceIdx, false);
  }
}

function measureStepWidth(): void {
  const track = trackRef.value;
  if (!track || count.value <= 1) {
    return;
  }
  const cards = track.children;
  if (cards.length < 2) {
    return;
  }
  const first = cards[0] as HTMLElement;
  const second = cards[1] as HTMLElement;
  const measured = second.offsetLeft - first.offsetLeft;
  if (measured > 0 && Math.abs(measured - stepWidth.value) > 0.5) {
    const oldStep = stepWidth.value;
    const oldCenterStart = getCenterStart(oldStep);
    const progressSteps = oldStep > 0 ? (oldCenterStart - currentX) / oldStep : 0;
    const targetSteps = oldStep > 0 ? (oldCenterStart - targetX) / oldStep : 0;

    stepWidth.value = measured;
    const newCenterStart = getCenterStart(measured);
    currentX = newCenterStart - progressSteps * measured;
    targetX = newCenterStart - targetSteps * measured;
    startX = currentX;
    applyTransform();
  }
}

// Autoplay scheduling
function stopAutoplay(): void {
  if (autoplayTimer !== null) {
    window.clearInterval(autoplayTimer);
    autoplayTimer = null;
  }
}

function startAutoplay(): void {
  stopAutoplay();
  if (!isAutoPlayActive.value) {
    return;
  }
  autoplayTimer = window.setInterval(() => {
    if (!document.hidden && isAutoPlayActive.value) {
      stepSlides(1, false);
    }
  }, AUTOPLAY_INTERVAL_MS);
}

function restartAutoplay(): void {
  stopAutoplay();
  startAutoplay();
}

function toggleAutoplay(): void {
  isUserPaused.value = !isUserPaused.value;
}

// Pointer drag / swipe handling (without blocking RouterLink clicks)
function handlePointerDown(e: PointerEvent): void {
  if (!isLooping.value || e.button !== 0) {
    return;
  }

  dragStartPointerX = e.clientX;
  dragStartPointerY = e.clientY;
  dragStartTrackX = wrapOffset(currentX);
  dragMoved = false;
  dragDirectionLocked = null;
  lastPointerX = e.clientX;
  lastPointerTime = performance.now();
  velocityX = 0;

  window.addEventListener('pointermove', handleWindowPointerMove, { passive: true });
  window.addEventListener('pointerup', handleWindowPointerUp);
  window.addEventListener('pointercancel', handleWindowPointerUp);
}

function handleWindowPointerMove(e: PointerEvent): void {
  const dx = e.clientX - dragStartPointerX;
  const dy = e.clientY - dragStartPointerY;

  if (dragDirectionLocked === null) {
    if (Math.abs(dx) < DRAG_THRESHOLD_PX && Math.abs(dy) < DRAG_THRESHOLD_PX) {
      return;
    }
    if (Math.abs(dy) > Math.abs(dx) && e.pointerType !== 'mouse') {
      dragDirectionLocked = 'vertical';
      cleanupWindowPointerListeners();
      return;
    }
    dragDirectionLocked = 'horizontal';
    dragMoved = true;
    isDragging.value = true;

    // Engage drag cleanly from current visual frame with zero threshold jump
    cancelAnimation();
    currentX = wrapOffset(currentX);
    targetX = currentX;
    startX = currentX;
    dragStartTrackX = currentX;
    dragStartPointerX = e.clientX;
    lastPointerX = e.clientX;
    lastPointerTime = performance.now();
    return;
  }

  if (dragDirectionLocked !== 'horizontal') {
    return;
  }

  const now = performance.now();
  const dt = Math.max(1, now - lastPointerTime);
  const instantVelocity = (e.clientX - lastPointerX) / dt;
  velocityX = velocityX * 0.4 + instantVelocity * 0.6;
  lastPointerX = e.clientX;
  lastPointerTime = now;

  currentX = dragStartTrackX + (e.clientX - dragStartPointerX);
  applyTransform();
}

function handleWindowPointerUp(): void {
  cleanupWindowPointerListeners();

  if (dragMoved) {
    suppressClick = true;
    if (suppressClickTimer !== null) {
      window.clearTimeout(suppressClickTimer);
    }
    suppressClickTimer = window.setTimeout(() => {
      suppressClick = false;
      suppressClickTimer = null;
    }, 100);

    const step = stepWidth.value;
    const flickMomentum = Math.max(-step * 1.4, Math.min(step * 1.4, velocityX * 190));
    snapToNearest(currentX + flickMomentum);
    restartAutoplay();
  }

  dragMoved = false;
  dragDirectionLocked = null;
  isDragging.value = false;
}

function cleanupWindowPointerListeners(): void {
  window.removeEventListener('pointermove', handleWindowPointerMove);
  window.removeEventListener('pointerup', handleWindowPointerUp);
  window.removeEventListener('pointercancel', handleWindowPointerUp);
}

function handleTrackClickCapture(e: MouseEvent): void {
  if (suppressClick) {
    e.preventDefault();
    e.stopPropagation();
    suppressClick = false;
  }
}

function handleWheel(e: WheelEvent): void {
  if (!isLooping.value) {
    return;
  }
  const absX = Math.abs(e.deltaX);
  const absY = Math.abs(e.deltaY);
  const horizontalIntent = absX > absY && absX > 2;
  const shiftWheelIntent = e.shiftKey && absY > 2;

  if (!horizontalIntent && !shiftWheelIntent) {
    return;
  }

  e.preventDefault();
  cancelAnimation();

  const delta = horizontalIntent ? e.deltaX : e.deltaY;
  currentX = wrapOffset(currentX) - delta;
  targetX = currentX;
  applyTransform();

  if (wheelSettleTimer !== null) {
    window.clearTimeout(wheelSettleTimer);
  }
  wheelSettleTimer = window.setTimeout(() => {
    wheelSettleTimer = null;
    snapToNearest();
    restartAutoplay();
  }, 140);
}

function handleKeydown(e: KeyboardEvent): void {
  if (e.key === 'ArrowLeft') {
    e.preventDefault();
    stepSlides(-1);
  } else if (e.key === 'ArrowRight') {
    e.preventDefault();
    stepSlides(1);
  }
}

function handleFocusIn(): void {
  isFocusedWithin.value = true;
}

function handleFocusOut(e: FocusEvent): void {
  const nextTarget = e.relatedTarget as Node | null;
  if (!sectionRef.value || !nextTarget || !sectionRef.value.contains(nextTarget)) {
    isFocusedWithin.value = false;
  }
}

function handleMotionChange(e: MediaQueryListEvent): void {
  prefersReducedMotion.value = e.matches;
}

useGsapScope(sectionRef, ({ gsap }) => {
  gsap.from('.carousel-header-anim', {
    y: 22,
    opacity: 0,
    duration: 0.6,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: sectionRef.value,
      start: 'top 82%',
      once: true,
    },
  });

  gsap.from('.carousel-viewport-anim', {
    y: 28,
    opacity: 0,
    duration: 0.65,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: viewportRef.value,
      start: 'top 82%',
      once: true,
    },
  });
});

watch(isAutoPlayActive, (active) => {
  if (active) {
    startAutoplay();
  } else {
    stopAutoplay();
  }
});

watch(
  () => props.services.length,
  async () => {
    await nextTick();
    measureStepWidth();
    resetToIndex(0);
    restartAutoplay();
  },
);

onMounted(async () => {
  motionQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
  prefersReducedMotion.value = motionQuery.matches;
  motionQuery.addEventListener('change', handleMotionChange);

  await nextTick();
  measureStepWidth();
  resetToIndex(0);

  if (typeof ResizeObserver !== 'undefined' && viewportRef.value) {
    resizeObserver = new ResizeObserver(() => {
      measureStepWidth();
    });
    resizeObserver.observe(viewportRef.value);
    if (trackRef.value) {
      resizeObserver.observe(trackRef.value);
    }
  }

  if (typeof IntersectionObserver !== 'undefined' && sectionRef.value) {
    intersectionObserver = new IntersectionObserver(
      (entries) => {
        isInViewport.value = entries[0]?.isIntersecting ?? true;
      },
      { threshold: 0.15 },
    );
    intersectionObserver.observe(sectionRef.value);
  }

  startAutoplay();
});

onUnmounted(() => {
  cancelAnimation();
  stopAutoplay();
  cleanupWindowPointerListeners();
  if (wheelSettleTimer !== null) {
    window.clearTimeout(wheelSettleTimer);
  }
  if (suppressClickTimer !== null) {
    window.clearTimeout(suppressClickTimer);
  }
  resizeObserver?.disconnect();
  intersectionObserver?.disconnect();
  motionQuery?.removeEventListener('change', handleMotionChange);
});
</script>

<template>
  <section
    ref="sectionRef"
    aria-roledescription="carousel"
    aria-label="Tiling services"
    class="px-3 md:px-4"
    @mouseenter="isHovered = true"
    @mouseleave="isHovered = false"
    @focusin="handleFocusIn"
    @focusout="handleFocusOut"
    @keydown="handleKeydown"
  >
    <div class="rounded-3xl bg-[#E8EFE6] py-16 md:py-24 px-6 md:px-12 overflow-hidden">
      <div class="max-w-[1320px] mx-auto">
        <!-- Header Row -->
        <div
          class="carousel-header-anim flex flex-col md:flex-row md:items-end justify-between gap-6 mb-10 md:mb-12"
        >
          <h2 class="text-3xl sm:text-4xl md:text-5xl font-medium tracking-[-0.03em] leading-[1.08]">
            <span class="block text-[#141A15]">One local crew.</span>
            <span class="block text-[#22C55E]">The whole interior &amp; patio.</span>
          </h2>

          <div class="flex items-center gap-4 md:gap-5">
            <p class="text-sm text-[#5A655C] max-w-[230px] leading-snug">
              Licensed &amp; insured tiling across Texas Hill Country
            </p>
            <div v-if="isLooping" class="flex items-center gap-2 shrink-0">
              <button
                type="button"
                aria-label="Previous service"
                class="w-11 h-11 rounded-full bg-white text-[#141A15] border border-[#E2E6E0] flex items-center justify-center hover:bg-[#1A241B] hover:text-white transition-colors shadow-sm cursor-pointer"
                @click="stepSlides(-1)"
              >
                <ArrowLeft class="w-4 h-4" />
              </button>
              <button
                type="button"
                aria-label="Next service"
                class="w-11 h-11 rounded-full bg-white text-[#141A15] border border-[#E2E6E0] flex items-center justify-center hover:bg-[#1A241B] hover:text-white transition-colors shadow-sm cursor-pointer"
                @click="stepSlides(1)"
              >
                <ArrowRight class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

        <!-- Infinite Looping Carousel Viewport -->
        <div
          ref="viewportRef"
          class="carousel-viewport-anim relative overflow-hidden select-none touch-pan-y"
          :class="
            isLooping
              ? isDragging
                ? 'cursor-grabbing'
                : 'cursor-grab'
              : ''
          "
          @pointerdown="handlePointerDown"
          @wheel="handleWheel"
          @click.capture="handleTrackClickCapture"
        >
          <div
            ref="trackRef"
            class="flex gap-5 will-change-transform pb-1"
          >
            <RouterLink
              v-for="item in renderedSlides"
              :key="item.key"
              :to="`/services/${item.service.slug}`"
              draggable="false"
              :tabindex="item.isPrimary ? 0 : -1"
              :aria-hidden="item.isPrimary ? undefined : true"
              class="group relative shrink-0 w-[252px] sm:w-[272px] xl:w-[286px] aspect-[3/4] rounded-3xl overflow-hidden bg-[#1A241B]"
              @dragstart.prevent
              @focus="item.isPrimary ? ensureCardVisibleOnFocus(item.serviceIdx) : undefined"
            >
              <img
                :src="item.service.imageUrl"
                :alt="item.isPrimary ? item.service.title : ''"
                draggable="false"
                decoding="async"
                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105 pointer-events-none"
              />
              <div
                class="absolute inset-0 bg-gradient-to-t from-[#141A15]/90 via-[#141A15]/20 to-transparent pointer-events-none"
              />

              <div class="absolute inset-x-0 bottom-0 p-6 text-white pointer-events-none">
                <div class="flex items-center justify-between gap-2">
                  <h3 class="text-2xl font-medium tracking-tight text-white">
                    {{ item.service.title }}
                  </h3>
                  <ArrowUpRight
                    class="w-5 h-5 text-[#4ADE80] opacity-0 -translate-x-1 translate-y-1 transition-all duration-300 group-hover:opacity-100 group-hover:translate-x-0 group-hover:translate-y-0"
                  />
                </div>
                <p
                  class="text-sm text-white/80 mt-1.5 max-h-0 opacity-0 overflow-hidden transition-all duration-300 group-hover:max-h-16 group-hover:opacity-100"
                >
                  {{ item.service.shortSubtitle }}
                </p>
              </div>
            </RouterLink>
          </div>
        </div>

        <!-- Bottom Carousel Pagination & Controls -->
        <div
          v-if="isLooping"
          class="mt-7 flex items-center justify-between gap-4"
        >
          <div class="flex items-center gap-2" role="group" aria-label="Service slide pagination">
            <button
              v-for="(service, idx) in services"
              :key="service.id"
              type="button"
              :aria-label="`Show ${service.title}`"
              :aria-pressed="idx === activeIndex"
              class="h-2 rounded-full transition-all duration-300 cursor-pointer"
              :class="
                idx === activeIndex
                  ? 'w-8 bg-[#1A241B]'
                  : 'w-2 bg-[#1A241B]/20 hover:bg-[#1A241B]/45'
              "
              @click="goToService(idx)"
            />
          </div>

          <div class="flex items-center gap-3">
            <span
              class="text-xs font-medium tabular-nums tracking-wider text-[#5A655C]"
              :aria-live="isAutoPlayActive ? 'off' : 'polite'"
            >
              {{ formattedSlideNumber }} / {{ formattedTotalSlides }}
            </span>

            <button
              type="button"
              :aria-label="isUserPaused ? 'Resume automatic rotation' : 'Pause automatic rotation'"
              :aria-pressed="isUserPaused"
              class="w-8 h-8 rounded-full bg-white/85 text-[#141A15] border border-[#E2E6E0] flex items-center justify-center hover:bg-[#1A241B] hover:text-white transition-colors shadow-sm cursor-pointer"
              @click="toggleAutoplay"
            >
              <Play v-if="isUserPaused" class="w-3.5 h-3.5 ml-0.5" />
              <Pause v-else class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
