<script setup lang="ts">
import { computed, ref } from 'vue';
import { useElementBounding } from '@vueuse/core';
import { ChevronLeft, ChevronRight } from 'lucide-vue-next';
import { useGsapScope } from '@/composables/useGsapScope';

interface Props {
  beforeImage: string;
  afterImage: string;
  beforeAlt?: string;
  afterAlt?: string;
  variant?: 'hero' | 'compact';
  initialPosition?: number;
  autoSweepOnVisible?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  beforeAlt: 'Before Copely restoration',
  afterAlt: 'After Copely installation',
  variant: 'compact',
  initialPosition: 50,
  autoSweepOnVisible: false,
});

const containerRef = ref<HTMLElement | null>(null);
const position = ref<number>(props.initialPosition);
const isDragging = ref<boolean>(false);
let demoSweepTl: gsap.core.Timeline | null = null;

const { left, width } = useElementBounding(containerRef);

useGsapScope(containerRef, ({ gsap }) => {
  if (!props.autoSweepOnVisible || !containerRef.value) {
    return;
  }

  const proxy = { val: props.initialPosition };
  demoSweepTl = gsap.timeline({
    scrollTrigger: {
      trigger: containerRef.value,
      start: 'top 76%',
      once: true,
    },
  });

  demoSweepTl
    .to(proxy, {
      val: 32,
      duration: 0.55,
      ease: 'power2.inOut',
      onUpdate: () => {
        position.value = Math.round(proxy.val * 10) / 10;
      },
    })
    .to(proxy, {
      val: 52,
      duration: 0.65,
      ease: 'power3.out',
      onUpdate: () => {
        position.value = Math.round(proxy.val * 10) / 10;
      },
    });
});

function cancelDemoSweep(): void {
  if (demoSweepTl) {
    demoSweepTl.kill();
    demoSweepTl = null;
  }
}

function updatePositionFromClientX(clientX: number): void {
  if (width.value <= 0) {
    return;
  }
  const pct = ((clientX - left.value) / width.value) * 100;
  position.value = Math.min(96, Math.max(4, Math.round(pct * 10) / 10));
}

function onPointerDown(event: PointerEvent): void {
  cancelDemoSweep();
  const el = event.currentTarget as HTMLElement | null;
  el?.setPointerCapture?.(event.pointerId);
  isDragging.value = true;
  updatePositionFromClientX(event.clientX);
}

function onPointerMove(event: PointerEvent): void {
  if (!isDragging.value) {
    return;
  }
  updatePositionFromClientX(event.clientX);
}

function onPointerUp(event: PointerEvent): void {
  const el = event.currentTarget as HTMLElement | null;
  if (el?.hasPointerCapture?.(event.pointerId)) {
    el.releasePointerCapture(event.pointerId);
  }
  isDragging.value = false;
}

function onKeyDown(event: KeyboardEvent): void {
  cancelDemoSweep();
  const step = event.shiftKey ? 10 : 3;
  if (event.key === 'ArrowLeft') {
    event.preventDefault();
    position.value = Math.max(4, position.value - step);
  } else if (event.key === 'ArrowRight') {
    event.preventDefault();
    position.value = Math.min(96, position.value + step);
  }
}

const clipPathStyle = computed(() => ({
  clipPath: `inset(0 ${100 - position.value}% 0 0)`,
}));
</script>

<template>
  <div
    ref="containerRef"
    class="relative w-full h-full overflow-hidden select-none touch-none cursor-ew-resize bg-[#1A241B]"
    :class="variant === 'hero' ? 'rounded-3xl aspect-[16/9]' : 'rounded-2xl aspect-[4/3]'"
    role="slider"
    tabindex="0"
    aria-label="Before and after comparison slider"
    aria-valuemin="4"
    aria-valuemax="96"
    :aria-valuenow="Math.round(position)"
    @pointerdown="onPointerDown"
    @pointermove="onPointerMove"
    @pointerup="onPointerUp"
    @pointercancel="onPointerUp"
    @keydown="onKeyDown"
  >
    <!-- After layer (base full image) -->
    <img
      :src="afterImage"
      :alt="afterAlt"
      draggable="false"
      class="absolute inset-0 w-full h-full object-cover pointer-events-none"
    />

    <!-- Before layer (clipped on the left up to slider position) -->
    <div class="absolute inset-0 pointer-events-none" :style="clipPathStyle">
      <img
        :src="beforeImage"
        :alt="beforeAlt"
        draggable="false"
        class="absolute inset-0 w-full h-full object-cover"
      />
    </div>

    <!-- Top-corner or Bottom-corner Badges per spec -->
    <template v-if="variant === 'hero'">
      <span
        class="pointer-events-none absolute top-5 left-5 z-10 rounded-full bg-[#1A241B]/80 backdrop-blur-sm border border-white/15 px-3.5 py-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-white"
      >
        BEFORE
      </span>
      <span
        class="pointer-events-none absolute top-5 right-5 z-10 rounded-full bg-[#34D361] px-3.5 py-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-[#141A15]"
      >
        AFTER
      </span>
    </template>
    <template v-else>
      <span
        class="pointer-events-none absolute bottom-3.5 left-3.5 z-10 rounded-full bg-[#1A241B]/80 backdrop-blur-sm px-2.5 py-0.5 text-[11px] font-medium text-white"
      >
        Before
      </span>
      <span
        class="pointer-events-none absolute bottom-3.5 right-3.5 z-10 rounded-full bg-white/90 backdrop-blur-sm px-2.5 py-0.5 text-[11px] font-medium text-[#141A15]"
      >
        After
      </span>
    </template>

    <!-- Vertical white hairline & center drag handle -->
    <div
      class="pointer-events-none absolute inset-y-0 z-20 w-[2px] -translate-x-1/2 bg-white shadow-[0_0_12px_rgba(0,0,0,0.45)]"
      :style="{ left: `${position}%` }"
    >
      <div
        class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-white text-[#141A15] shadow-xl flex items-center justify-center ring-4 ring-black/15 transition-transform duration-150"
        :class="isDragging ? 'scale-110' : 'scale-100'"
      >
        <ChevronLeft class="w-4 h-4 -mr-0.5 stroke-[2.2]" />
        <ChevronRight class="w-4 h-4 -ml-0.5 stroke-[2.2]" />
      </div>
    </div>
  </div>
</template>
