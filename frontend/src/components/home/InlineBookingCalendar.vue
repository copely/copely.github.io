<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue';
import gsap from 'gsap';
import {
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  Clock,
  Globe,
  MapPin,
} from 'lucide-vue-next';
import { SERVICE_OPTIONS, TIMELINE_OPTIONS } from '@/data/copelyData';
import { useEstimateDraft } from '@/composables/useEstimateDraft';
import { allowsMotion, useGsapScope } from '@/composables/useGsapScope';
import { createBooking, fetchBookingAvailability } from '@/services/api';
import type {
  BookingResponse,
  DayAvailability,
  ProjectTimeline,
  TileServiceCategory,
} from '@/types/copely';

interface Props {
  towns: string[];
}

const props = defineProps<Props>();

const { estimateDraft, clearHighlight } = useEstimateDraft();

const currentMonth = ref<string>('2026-11');
const availabilityDays = ref<DayAvailability[]>([]);
const isLoadingSlots = ref<boolean>(false);
const selectedDate = ref<string>('2026-11-18');
const selectedSlot = ref<string>('');
const selectedTown = ref<string>(props.towns[0] ?? 'Dripping Springs');
const selectedService = ref<TileServiceCategory>(estimateDraft.service);
const selectedTimeline = ref<ProjectTimeline>(estimateDraft.timeline);
const use24Hour = ref<boolean>(false);

const customerName = ref<string>('');
const customerEmail = ref<string>('');
const customerPhone = ref<string>('');
const isSubmitting = ref<boolean>(false);
const bookingError = ref<string>('');
const confirmedBooking = ref<BookingResponse | null>(null);
const scheduleSectionRef = ref<HTMLElement | null>(null);
const slotsListRef = ref<HTMLElement | null>(null);
const calendarGridRef = ref<HTMLElement | null>(null);
const dayHeaderRef = ref<HTMLElement | null>(null);
const confirmDrawerRef = ref<HTMLElement | null>(null);
const confirmedCardRef = ref<HTMLElement | null>(null);
const confirmCircleRef = ref<SVGCircleElement | null>(null);
const confirmCheckRef = ref<SVGPathElement | null>(null);
const heroBannerRef = ref<HTMLElement | null>(null);
const calServiceWrapRef = ref<HTMLElement | null>(null);
const timeToggleThumbRef = ref<HTMLElement | null>(null);
const wayCountRef = ref<HTMLElement | null>(null);

let bookingCtx: gsap.Context | null = null;
let hasCardEnteredViewport = false;
let hasPlayedInitialGridCascade = false;
let pendingMonthDirection = 0;
let heroHandoffTimer: number | null = null;

function pulseSelectedDayRing(el: HTMLElement | null): void {
  if (!el || !allowsMotion()) {
    return;
  }
  gsap.fromTo(
    el,
    { boxShadow: '0 0 0 0 rgba(52, 211, 97, 0.58)' },
    {
      boxShadow: '0 0 0 11px rgba(52, 211, 97, 0)',
      duration: 0.78,
      ease: 'power2.out',
      clearProps: 'boxShadow',
    },
  );
}

async function runCalendarGridCascade(direction = 0): Promise<void> {
  if (!allowsMotion()) {
    return;
  }
  await nextTick();
  const grid = calendarGridRef.value;
  if (!grid) {
    return;
  }
  const slots = grid.querySelectorAll<HTMLElement>('.cal-grid-slot');
  if (!slots.length) {
    return;
  }
  const rows = Math.max(4, Math.ceil(slots.length / 7));
  const activeBtn = grid.querySelector<HTMLElement>('[data-selected-day="true"]');

  bookingCtx?.add(() => {
    gsap.killTweensOf(slots);
    gsap.fromTo(
      slots,
      {
        scale: 0.72,
        y: 10,
        x: direction * 18,
        opacity: 0,
      },
      {
        scale: 1,
        y: 0,
        x: 0,
        opacity: 1,
        duration: 0.34,
        ease: 'back.out(1.45)',
        stagger: {
          amount: 0.22,
          grid: [rows, 7],
          from: direction < 0 ? 'end' : 0,
        },
        clearProps: 'transform',
        onComplete: () => {
          pulseSelectedDayRing(activeBtn);
        },
      },
    );
  });
}

async function animateSlotPills(): Promise<void> {
  if (!allowsMotion()) {
    return;
  }
  await nextTick();
  const pills = slotsListRef.value?.querySelectorAll<HTMLElement>('.slot-pill-anim');
  if (pills && pills.length > 0) {
    bookingCtx?.add(() => {
      gsap.killTweensOf(pills);
      gsap.fromTo(
        pills,
        { x: 12, opacity: 0 },
        {
          x: 0,
          opacity: 1,
          duration: 0.28,
          stagger: { amount: 0.16 },
          ease: 'back.out(1.3)',
          clearProps: 'transform,opacity',
        },
      );
    });
  }
}

function handleDayHover(day: DayAvailability, event: MouseEvent, entering: boolean): void {
  if (!allowsMotion() || !day.available || selectedDate.value === day.date) {
    return;
  }
  const btn = event.currentTarget as HTMLElement | null;
  if (!btn) {
    return;
  }
  gsap.killTweensOf(btn, 'scale,y');
  gsap.to(btn, {
    scale: entering ? 1.07 : 1,
    y: entering ? -2 : 0,
    duration: entering ? 0.22 : 0.2,
    ease: entering ? 'back.out(2)' : 'power2.out',
    ...(entering ? {} : { clearProps: 'transform' }),
  });
}

async function setTimeFormat(is24: boolean): Promise<void> {
  if (use24Hour.value === is24) {
    return;
  }
  use24Hour.value = is24;
  if (!allowsMotion()) {
    if (timeToggleThumbRef.value) {
      timeToggleThumbRef.value.style.transform = is24 ? 'translateX(100%)' : 'translateX(0%)';
    }
    return;
  }
  bookingCtx?.add(() => {
    if (timeToggleThumbRef.value) {
      gsap.to(timeToggleThumbRef.value, {
        xPercent: is24 ? 100 : 0,
        duration: 0.3,
        ease: 'back.out(1.4)',
      });
    }
  });
  await nextTick();
  const labels = slotsListRef.value?.querySelectorAll<HTMLElement>('.slot-time-label');
  if (labels && labels.length > 0) {
    bookingCtx?.add(() => {
      gsap.fromTo(
        labels,
        { rotationX: 75, opacity: 0.2, transformPerspective: 400 },
        {
          rotationX: 0,
          opacity: 1,
          duration: 0.24,
          stagger: 0.018,
          ease: 'power2.out',
          clearProps: 'transform,opacity',
        },
      );
    });
  }
}

useGsapScope(scheduleSectionRef, ({ gsap: g }) => {
  if (timeToggleThumbRef.value) {
    g.set(timeToggleThumbRef.value, { xPercent: use24Hour.value ? 100 : 0 });
  }

  g.fromTo(
    '.cal-live-halo',
    { scale: 1, opacity: 0.65 },
    {
      scale: 2.3,
      opacity: 0,
      duration: 1.7,
      repeat: -1,
      ease: 'power2.out',
    },
  );

  const headerTl = g.timeline({
    scrollTrigger: {
      trigger: scheduleSectionRef.value,
      start: 'top 80%',
      once: true,
    },
  });

  headerTl
    .fromTo(
      '.cal-eyebrow-pill',
      { y: 12, opacity: 0, scale: 0.94 },
      { y: 0, opacity: 1, scale: 1, duration: 0.45, ease: 'back.out(1.6)', clearProps: 'transform,opacity' },
    )
    .fromTo(
      '.cal-head-line',
      { clipPath: 'inset(100% 0% 0% 0%)', y: 18 },
      {
        clipPath: 'inset(0% 0% 0% 0%)',
        y: 0,
        duration: 0.68,
        stagger: 0.08,
        ease: 'expo.out',
        clearProps: 'clipPath,transform',
      },
      '-=0.25',
    )
    .fromTo(
      '.cal-subtext',
      { y: 12, opacity: 0 },
      { y: 0, opacity: 1, duration: 0.48, ease: 'power2.out', clearProps: 'transform,opacity' },
      '-=0.38',
    );

  const cardTl = g.timeline({
    scrollTrigger: {
      trigger: '.booking-widget-card',
      start: 'top 84%',
      once: true,
      onEnter: () => {
        hasCardEnteredViewport = true;
      },
    },
  });

  cardTl
    .fromTo(
      '.booking-widget-card',
      {
        y: 38,
        rotationX: 5,
        scale: 0.975,
        transformPerspective: 1100,
        clipPath: 'inset(0% 0% 10% 0% round 24px)',
        opacity: 0,
      },
      {
        y: 0,
        rotationX: 0,
        scale: 1,
        clipPath: 'inset(0% 0% 0% 0% round 24px)',
        opacity: 1,
        duration: 0.76,
        ease: 'power4.out',
        clearProps: 'transform,clipPath,opacity',
      },
    )
    .fromTo(
      '.booking-col-panel',
      { x: -16, opacity: 0 },
      {
        x: 0,
        opacity: 1,
        duration: 0.52,
        stagger: 0.09,
        ease: 'power3.out',
        clearProps: 'transform,opacity',
      },
      '-=0.45',
    )
    .fromTo(
      '.copely-diamond-mark',
      { rotation: -90, scale: 0.5 },
      { rotation: 0, scale: 1, duration: 0.52, ease: 'back.out(2)', clearProps: 'transform' },
      '-=0.4',
    )
    .call(
      () => {
        if (!isLoadingSlots.value && availabilityDays.value.length > 0 && !hasPlayedInitialGridCascade) {
          hasPlayedInitialGridCascade = true;
          void runCalendarGridCascade(0);
          void animateSlotPills();
        }
      },
      [],
      '-=0.22',
    );

  const trustTl = g.timeline({
    scrollTrigger: {
      trigger: '.financing-trust-bar',
      start: 'top 95%',
      once: true,
    },
  });

  const counterProxy = { ways: 0 };
  trustTl
    .fromTo(
      '.trust-item-anim',
      { y: 10, opacity: 0 },
      {
        y: 0,
        opacity: 1,
        duration: 0.42,
        stagger: 0.06,
        ease: 'power2.out',
        clearProps: 'transform,opacity',
      },
    )
    .fromTo(
      '.trust-dot-sep',
      { scale: 0, opacity: 0 },
      { scale: 1, opacity: 1, duration: 0.3, stagger: 0.05, ease: 'back.out(2)', clearProps: 'transform' },
      '-=0.3',
    )
    .to(
      counterProxy,
      {
        ways: 3,
        duration: 0.65,
        ease: 'power2.out',
        onUpdate: () => {
          if (wayCountRef.value) {
            wayCountRef.value.textContent = String(Math.round(counterProxy.ways));
          }
        },
      },
      '<',
    );

  bookingCtx = g.context(() => {}, scheduleSectionRef.value ?? undefined);
});

watch(
  () => [estimateDraft.service, estimateDraft.timeline, estimateDraft.highlightedFromHero] as const,
  async ([svc, tl, highlighted]) => {
    selectedService.value = svc;
    selectedTimeline.value = tl;
    if (!highlighted) {
      return;
    }
    confirmedBooking.value = null;
    if (!allowsMotion()) {
      return;
    }
    await nextTick();
    if (heroBannerRef.value) {
      gsap.fromTo(
        heroBannerRef.value,
        { scale: 0.85, opacity: 0 },
        { scale: 1, opacity: 1, duration: 0.45, ease: 'back.out(2)', clearProps: 'transform,opacity' },
      );
    }
    if (heroHandoffTimer !== null) {
      window.clearTimeout(heroHandoffTimer);
    }
    heroHandoffTimer = window.setTimeout(() => {
      heroHandoffTimer = null;
      const card = scheduleSectionRef.value?.querySelector<HTMLElement>('.booking-widget-card');
      if (card) {
        gsap.fromTo(
          card,
          { boxShadow: '0 0 0 3px rgba(52, 211, 97, 0.58), 0 20px 25px -5px rgba(0, 0, 0, 0.1)' },
          {
            boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)',
            duration: 1.05,
            ease: 'power2.out',
            clearProps: 'boxShadow',
          },
        );
      }
      if (calServiceWrapRef.value) {
        gsap.fromTo(
          calServiceWrapRef.value,
          { scale: 1, boxShadow: '0 0 0 2px rgba(52, 211, 97, 0.65)' },
          {
            scale: 1.02,
            duration: 0.22,
            yoyo: true,
            repeat: 1,
            ease: 'power2.out',
            onComplete: () => {
              if (calServiceWrapRef.value) {
                gsap.set(calServiceWrapRef.value, { clearProps: 'transform,boxShadow' });
              }
            },
          },
        );
      }
    }, 450);
  },
);

watch(selectedSlot, async (newSlot, oldSlot) => {
  if (!allowsMotion() || !newSlot) {
    return;
  }
  await nextTick();
  bookingCtx?.add(() => {
    const drawer = confirmDrawerRef.value;
    if (!drawer) {
      return;
    }
    if (!oldSlot) {
      const targetHeight = drawer.scrollHeight;
      gsap.fromTo(
        drawer,
        { height: 0, opacity: 0 },
        {
          height: targetHeight,
          opacity: 1,
          duration: 0.38,
          ease: 'expo.out',
          onComplete: () => {
            gsap.set(drawer, { height: 'auto' });
          },
        },
      );
      const fields = drawer.querySelectorAll('.confirm-field');
      if (fields.length > 0) {
        gsap.fromTo(
          fields,
          { y: 10, opacity: 0 },
          {
            y: 0,
            opacity: 1,
            duration: 0.3,
            stagger: 0.045,
            delay: 0.05,
            ease: 'power2.out',
            clearProps: 'transform,opacity',
          },
        );
      }
    } else {
      const cta = drawer.querySelector('.book-cta-btn');
      if (cta) {
        gsap.fromTo(
          cta,
          { scale: 0.96 },
          { scale: 1, duration: 0.28, ease: 'back.out(2)', clearProps: 'transform' },
        );
      }
    }
  });
});

watch(confirmedBooking, async (booking, prevBooking) => {
  await nextTick();
  if (!booking) {
    if (prevBooking) {
      if (timeToggleThumbRef.value) {
        if (allowsMotion()) {
          gsap.set(timeToggleThumbRef.value, { xPercent: use24Hour.value ? 100 : 0 });
        } else {
          timeToggleThumbRef.value.style.transform = use24Hour.value
            ? 'translateX(100%)'
            : 'translateX(0%)';
        }
      }
      void animateSlotPills();
    }
    return;
  }
  if (!allowsMotion()) {
    return;
  }
  bookingCtx?.add(() => {
    if (confirmedCardRef.value) {
      gsap.fromTo(
        confirmedCardRef.value,
        { scale: 0.94, opacity: 0, y: 10 },
        { scale: 1, opacity: 1, y: 0, duration: 0.45, ease: 'back.out(1.5)', clearProps: 'transform,opacity' },
      );
    }
    const circle = confirmCircleRef.value;
    const check = confirmCheckRef.value;
    if (circle && check) {
      const circleLen = circle.getTotalLength();
      const checkLen = check.getTotalLength();
      gsap.set(circle, { strokeDasharray: circleLen, strokeDashoffset: circleLen });
      gsap.set(check, { strokeDasharray: checkLen, strokeDashoffset: checkLen });
      gsap
        .timeline()
        .to(circle, { strokeDashoffset: 0, duration: 0.45, ease: 'power2.out' })
        .to(check, { strokeDashoffset: 0, duration: 0.32, ease: 'power3.out' }, '-=0.14');
    }
    const items = confirmedCardRef.value?.querySelectorAll('.confirm-summary-anim');
    if (items && items.length > 0) {
      gsap.fromTo(
        items,
        { y: 8, opacity: 0 },
        { y: 0, opacity: 1, duration: 0.32, stagger: 0.06, delay: 0.12, ease: 'power2.out', clearProps: 'transform,opacity' },
      );
    }
  });
});

const serviceLabel = computed<string>(() => {
  return (
    SERVICE_OPTIONS.find((item) => item.value === selectedService.value)?.label ??
    'Porcelain & Tile'
  );
});

const timelineLabel = computed<string>(() => {
  return (
    TIMELINE_OPTIONS.find((item) => item.value === selectedTimeline.value)?.label ?? 'This week'
  );
});

const monthTitle = computed<string>(() => {
  const [yrStr, moStr] = currentMonth.value.split('-');
  const yr = Number(yrStr);
  const mo = Number(moStr) - 1;
  const d = new Date(Date.UTC(yr, mo, 1));
  return d.toLocaleDateString('en-US', { month: 'long', year: 'numeric', timeZone: 'UTC' });
});

const leadingBlankDays = computed<number>(() => {
  if (availabilityDays.value.length === 0) {
    // 2026-11-01 is Sunday (0)
    return 0;
  }
  return availabilityDays.value[0]?.dayOfWeek ?? 0;
});

const activeDayObj = computed<DayAvailability | undefined>(() => {
  return availabilityDays.value.find((day) => day.date === selectedDate.value);
});

const availableSlotsForDay = computed(() => {
  return (activeDayObj.value?.slots ?? []).filter((slot) => slot.available);
});

const selectedDayHeader = computed<string>(() => {
  if (!selectedDate.value) {
    return '';
  }
  const parts = selectedDate.value.split('-').map(Number);
  const d = new Date(Date.UTC(parts[0] ?? 2026, (parts[1] ?? 11) - 1, parts[2] ?? 18));
  return d.toLocaleDateString('en-US', {
    weekday: 'short',
    day: 'numeric',
    timeZone: 'UTC',
  });
});

function formatSlotTime(slot24: string): string {
  if (use24Hour.value) {
    return slot24;
  }
  const [hStr, mStr] = slot24.split(':');
  const hour = Number(hStr);
  const suffix = hour >= 12 ? 'pm' : 'am';
  const hour12 = hour % 12 === 0 ? 12 : hour % 12;
  return `${hour12}:${mStr}${suffix}`;
}

function buildFallbackMonthDays(monthStr: string): DayAvailability[] {
  const [yrStr, moStr] = monthStr.split('-');
  const yr = Number(yrStr);
  const mo = Number(moStr);
  const count = new Date(Date.UTC(yr, mo, 0)).getUTCDate();
  const baseSlots = [
    '08:00',
    '08:30',
    '09:00',
    '09:30',
    '10:00',
    '10:30',
    '11:00',
    '13:00',
    '13:30',
    '14:00',
    '15:00',
    '16:00',
  ];
  const days: DayAvailability[] = [];
  for (let dayNum = 1; dayNum <= count; dayNum++) {
    const d = new Date(Date.UTC(yr, mo - 1, dayNum));
    const dow = d.getUTCDay();
    const iso = `${monthStr}-${String(dayNum).padStart(2, '0')}`;
    const open = dow !== 0 && iso >= '2026-11-18' && iso !== '2026-11-26';
    days.push({
      date: iso,
      dayOfWeek: dow,
      available: open,
      slots: baseSlots.map((time) => ({ time, available: open })),
    });
  }
  return days;
}

async function loadAvailability(month: string): Promise<void> {
  isLoadingSlots.value = true;
  bookingError.value = '';
  try {
    const resp = await fetchBookingAvailability(month);
    availabilityDays.value = resp.days;
  } catch {
    availabilityDays.value = buildFallbackMonthDays(month);
  } finally {
    isLoadingSlots.value = false;
  }

  const currentStillValid = availabilityDays.value.some(
    (d) => d.date === selectedDate.value && d.available,
  );
  if (!currentStillValid) {
    const firstOpen = availabilityDays.value.find((d) => d.available);
    selectedDate.value = firstOpen?.date ?? '';
    selectedSlot.value = '';
  }

  if (hasCardEnteredViewport) {
    if (!hasPlayedInitialGridCascade) {
      hasPlayedInitialGridCascade = true;
      void runCalendarGridCascade(0);
      void animateSlotPills();
    } else if (pendingMonthDirection !== 0) {
      const dir = pendingMonthDirection;
      pendingMonthDirection = 0;
      void runCalendarGridCascade(dir);
      void animateSlotPills();
    }
  } else {
    pendingMonthDirection = 0;
  }
}

function shiftMonth(delta: number): void {
  const [yrStr, moStr] = currentMonth.value.split('-');
  const d = new Date(Date.UTC(Number(yrStr), Number(moStr) - 1 + delta, 1));
  const nextMonth = `${d.getUTCFullYear()}-${String(d.getUTCMonth() + 1).padStart(2, '0')}`;
  if (nextMonth < '2026-11' || nextMonth > '2027-03') {
    return;
  }
  pendingMonthDirection = delta;
  currentMonth.value = nextMonth;
  selectedSlot.value = '';
  void loadAvailability(nextMonth);
}

async function selectDay(day: DayAvailability): Promise<void> {
  if (!day.available) {
    return;
  }
  const dayIdx = availabilityDays.value.findIndex((d) => d.date === day.date);
  const clickedGridIndex = leadingBlankDays.value + Math.max(0, dayIdx);

  selectedDate.value = day.date;
  selectedSlot.value = '';
  confirmedBooking.value = null;
  bookingError.value = '';

  if (allowsMotion()) {
    await nextTick();
    bookingCtx?.add(() => {
      const grid = calendarGridRef.value;
      const slots = grid?.querySelectorAll<HTMLElement>('.cal-grid-slot');
      if (slots && slots.length > 0) {
        const rows = Math.max(4, Math.ceil(slots.length / 7));
        gsap.killTweensOf(slots);
        gsap.fromTo(
          slots,
          { scale: 0.91 },
          {
            scale: 1,
            duration: 0.28,
            ease: 'back.out(2.1)',
            stagger: {
              amount: 0.18,
              from: clickedGridIndex,
              grid: [rows, 7],
            },
            clearProps: 'transform',
          },
        );
      }

      const activeBtn = grid?.querySelector<HTMLElement>('[data-selected-day="true"]');
      if (activeBtn) {
        gsap.fromTo(
          activeBtn,
          { scale: 0.88 },
          { scale: 1, duration: 0.44, ease: 'elastic.out(1, 0.58)', clearProps: 'transform' },
        );
        pulseSelectedDayRing(activeBtn);
      }

      if (dayHeaderRef.value) {
        gsap.fromTo(
          dayHeaderRef.value,
          { y: -8, opacity: 0.35 },
          { y: 0, opacity: 1, duration: 0.32, ease: 'power2.out', clearProps: 'transform,opacity' },
        );
      }
    });
  }

  void animateSlotPills();
}

async function confirmSlotBooking(): Promise<void> {
  if (!selectedDate.value || !selectedSlot.value || !selectedTown.value) {
    return;
  }
  isSubmitting.value = true;
  bookingError.value = '';
  try {
    const result = await createBooking({
      date: selectedDate.value,
      timeSlot: selectedSlot.value,
      town: selectedTown.value,
      service: selectedService.value,
      timeline: selectedTimeline.value,
      customerName: customerName.value.trim() || undefined,
      customerEmail: customerEmail.value.trim() || undefined,
      customerPhone: customerPhone.value.trim() || undefined,
    });
    confirmedBooking.value = result;
    clearHighlight();
    await loadAvailability(currentMonth.value);
  } catch (err) {
    bookingError.value =
      err instanceof Error ? err.message : 'Unable to reserve that slot. Please try another.';
  } finally {
    isSubmitting.value = false;
  }
}

onMounted(() => {
  void loadAvailability(currentMonth.value);
});

onUnmounted(() => {
  if (heroHandoffTimer !== null) {
    window.clearTimeout(heroHandoffTimer);
    heroHandoffTimer = null;
  }
  bookingCtx?.revert();
  bookingCtx = null;
});
</script>

<template>
  <section id="schedule" ref="scheduleSectionRef" class="py-16 md:py-24 px-6 md:px-10 bg-[#F7F7F3] scroll-mt-24 overflow-x-clip">
    <div class="max-w-[1160px] mx-auto">
      <!-- Centered Header -->
      <div class="text-center space-y-3 mb-12">
        <div
          class="cal-eyebrow-pill inline-flex items-center gap-2 rounded-full bg-white border border-[#E2E6E0] px-3.5 py-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-[#1A241B] shadow-sm"
        >
          <span class="relative flex h-2 w-2 items-center justify-center">
            <span class="cal-live-halo absolute inline-flex h-full w-full rounded-full bg-[#34D361]" />
            <span class="relative inline-flex h-2 w-2 rounded-full bg-[#34D361]" />
          </span>
          <span>Live Hill Country schedule</span>
        </div>

        <h2 class="text-4xl md:text-5xl font-medium tracking-[-0.03em] text-[#141A15]">
          <span class="inline-block overflow-hidden align-bottom">
            <span class="cal-head-line inline-block">Pick a day.&nbsp;</span>
          </span>
          <span class="inline-block overflow-hidden align-bottom">
            <span class="cal-head-line inline-block font-serif italic font-normal">We’ll be there.</span>
          </span>
        </h2>
        <p class="cal-subtext text-base text-[#5A655C] max-w-xl mx-auto">
          An honest look at your floors or walls and a written estimate — whether or not you hire
          us. Active leak? Call
          <a href="tel:+15125550142" class="font-semibold text-[#141A15] underline decoration-[#34D361]">
            (512) 555-0142
          </a>
        </p>

        <!-- Pre-filled Selection Banner when routed from Hero card -->
        <div
          v-if="estimateDraft.highlightedFromHero"
          ref="heroBannerRef"
          class="inline-flex flex-wrap items-center justify-center gap-2 rounded-full bg-[#DCE8D9] border border-[#34D361]/40 px-4 py-1.5 text-xs font-medium text-[#141A15] mt-2"
        >
          <span class="w-2 h-2 rounded-full bg-[#22C55E]" />
          <span>Pre-filled from your estimate:</span>
          <strong>{{ serviceLabel }}</strong>
          <span>·</span>
          <strong>{{ timelineLabel }}</strong>
        </div>
      </div>

      <!-- 3-Column Cal.com-Style Interactive Widget Card -->
      <div
        class="booking-widget-card rounded-3xl bg-white border border-[#E2E6E0] shadow-xl overflow-hidden grid grid-cols-1 lg:grid-cols-12"
      >
        <!-- Column 1: Session Details -->
        <div class="booking-col-panel lg:col-span-4 p-6 md:p-8 border-b lg:border-b-0 lg:border-r border-[#E2E6E0] flex flex-col justify-between gap-8">
          <div class="space-y-5">
            <!-- Copely circular dark badge -->
            <div
              class="inline-flex items-center gap-2 rounded-full bg-[#1A241B] text-white px-3.5 py-1.5 text-xs font-semibold"
            >
              <span class="copely-diamond-mark inline-block text-[#4ADE80]" aria-hidden="true">◆</span>
              <span>Copely</span>
            </div>

            <div>
              <h3 class="text-2xl font-medium tracking-tight text-[#141A15]">
                Free tile estimate
              </h3>
              <p class="text-sm text-[#5A655C] mt-1">
                Laser room measurement, subfloor moisture check &amp; written line-item quote.
              </p>
            </div>

            <ul class="space-y-3 text-sm text-[#141A15]">
              <li class="flex items-center gap-2.5">
                <Clock class="w-4 h-4 text-[#22C55E]" />
                <span>30 min</span>
              </li>
              <li class="flex items-center gap-2.5">
                <MapPin class="w-4 h-4 text-[#22C55E]" />
                <span>On-site visit</span>
              </li>
              <li class="flex items-center gap-2.5">
                <Globe class="w-4 h-4 text-[#22C55E]" />
                <span>Texas Hill Country</span>
              </li>
            </ul>

            <!-- Pre-filled Service selector -->
            <div class="pt-2">
              <label
                for="cal-service"
                class="block text-xs font-medium text-[#5A655C] mb-1.5"
              >
                Project focus
              </label>
              <div ref="calServiceWrapRef" class="relative rounded-xl">
                <select
                  id="cal-service"
                  v-model="selectedService"
                  class="w-full appearance-none rounded-xl bg-[#F7F7F3] border border-[#E2E6E0] px-3.5 py-2.5 pr-9 text-sm font-medium text-[#141A15] focus:outline-none focus:border-[#1A241B]"
                >
                  <option v-for="opt in SERVICE_OPTIONS" :key="opt.value" :value="opt.value">
                    {{ opt.label }}
                  </option>
                </select>
                <ChevronDown
                  class="w-4 h-4 text-[#5A655C] absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none"
                />
              </div>
            </div>
          </div>

          <!-- Location Selector Dropdown at bottom -->
          <div>
            <label for="cal-town" class="block text-xs font-medium text-[#5A655C] mb-1.5">
              What’s your location?
            </label>
            <div class="relative">
              <select
                id="cal-town"
                v-model="selectedTown"
                class="w-full appearance-none rounded-xl bg-[#F7F7F3] border border-[#E2E6E0] px-3.5 py-2.5 pr-9 text-sm font-medium text-[#141A15] focus:outline-none focus:border-[#1A241B] cursor-pointer"
              >
                <option value="" disabled>Choose your town</option>
                <option v-for="town in towns" :key="town" :value="town">
                  {{ town }}
                </option>
              </select>
              <ChevronDown
                class="w-4 h-4 text-[#5A655C] absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none"
              />
            </div>
          </div>
        </div>

        <!-- Column 2: Interactive Month Calendar -->
        <div class="booking-col-panel lg:col-span-5 p-6 md:p-8 border-b lg:border-b-0 lg:border-r border-[#E2E6E0]">
          <div class="flex items-center justify-between mb-6">
            <h4 class="text-base font-semibold text-[#141A15]">{{ monthTitle }}</h4>
            <div class="flex items-center gap-1.5">
              <button
                type="button"
                aria-label="Previous month"
                :disabled="currentMonth <= '2026-11'"
                class="w-8 h-8 rounded-full border border-[#E2E6E0] flex items-center justify-center text-[#141A15] hover:bg-[#E8EFE6] disabled:opacity-35 disabled:cursor-not-allowed cursor-pointer transition-colors"
                @click="shiftMonth(-1)"
              >
                <ChevronLeft class="w-4 h-4" />
              </button>
              <button
                type="button"
                aria-label="Next month"
                :disabled="currentMonth >= '2027-03'"
                class="w-8 h-8 rounded-full border border-[#E2E6E0] flex items-center justify-center text-[#141A15] hover:bg-[#E8EFE6] disabled:opacity-35 disabled:cursor-not-allowed cursor-pointer transition-colors"
                @click="shiftMonth(1)"
              >
                <ChevronRight class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Weekday Labels -->
          <div class="grid grid-cols-7 gap-1.5 text-center mb-2">
            <span
              v-for="dayLabel in ['SUN', 'MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT']"
              :key="dayLabel"
              class="text-[11px] font-semibold tracking-wider text-[#5A655C] py-1"
            >
              {{ dayLabel }}
            </span>
          </div>

          <!-- Day Grid -->
          <div ref="calendarGridRef" class="grid grid-cols-7 gap-1.5">
            <div
              v-for="blank in leadingBlankDays"
              :key="`blank-${blank}`"
              class="cal-grid-slot pointer-events-none"
              aria-hidden="true"
            />
            <button
              v-for="day in availabilityDays"
              :key="day.date"
              type="button"
              :disabled="!day.available"
              :data-selected-day="selectedDate === day.date ? 'true' : undefined"
              class="cal-day-cell cal-grid-slot aspect-square rounded-xl text-sm font-medium transition-colors flex items-center justify-center"
              :class="
                selectedDate === day.date
                  ? 'bg-[#1A241B] text-white shadow-md font-semibold cursor-pointer'
                  : day.available
                    ? 'bg-[#DCE8D9] text-[#141A15] hover:bg-[#34D361]/40 cursor-pointer'
                    : 'text-[#5A655C]/35 cursor-not-allowed'
              "
              @mouseenter="handleDayHover(day, $event, true)"
              @mouseleave="handleDayHover(day, $event, false)"
              @click="selectDay(day)"
            >
              {{ Number(day.date.split('-')[2]) }}
            </button>
          </div>
        </div>

        <!-- Column 3: Time Slot Picker & Confirmation -->
        <div class="booking-col-panel lg:col-span-3 p-6 md:p-8 flex flex-col">
          <!-- Confirmed View -->
          <div
            v-if="confirmedBooking"
            ref="confirmedCardRef"
            class="flex-1 flex flex-col justify-center items-start space-y-4 bg-[#E8EFE6] rounded-2xl p-5"
          >
            <svg
              class="w-9 h-9 text-[#22C55E]"
              viewBox="0 0 36 36"
              fill="none"
              aria-hidden="true"
            >
              <circle
                ref="confirmCircleRef"
                cx="18"
                cy="18"
                r="15"
                stroke="currentColor"
                stroke-width="2.4"
              />
              <path
                ref="confirmCheckRef"
                d="M11.5 18.5L16 23L25 13.5"
                stroke="currentColor"
                stroke-width="2.6"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            <div class="space-y-1">
              <span class="confirm-summary-anim inline-block text-xs font-semibold uppercase tracking-wider text-[#22C55E]">
                Confirmed · {{ confirmedBooking.bookingId }}
              </span>
              <h4 class="confirm-summary-anim text-lg font-semibold text-[#141A15]">Visit Scheduled</h4>
              <p class="confirm-summary-anim text-sm text-[#364137]">
                {{ confirmedBooking.date }} at {{ formatSlotTime(confirmedBooking.timeSlot) }} in
                {{ confirmedBooking.town }}
              </p>
            </div>
            <button
              type="button"
              class="confirm-summary-anim text-xs font-semibold underline text-[#141A15] cursor-pointer"
              @click="confirmedBooking = null"
            >
              Book another time
            </button>
          </div>

          <!-- Active Slot Picker -->
          <template v-else>
            <div class="flex items-center justify-between gap-2 mb-4">
              <span ref="dayHeaderRef" class="text-sm font-semibold text-[#141A15] inline-block">
                {{ selectedDayHeader || 'Select date' }}
              </span>
              <div class="relative grid grid-cols-2 rounded-full bg-[#F7F7F3] border border-[#E2E6E0] p-0.5 text-xs">
                <span
                  ref="timeToggleThumbRef"
                  aria-hidden="true"
                  class="pointer-events-none absolute top-0.5 bottom-0.5 left-0.5 w-[calc(50%-2px)] rounded-full bg-[#1A241B]"
                />
                <button
                  type="button"
                  class="relative z-10 px-2.5 py-1 rounded-full font-medium transition-colors cursor-pointer"
                  :class="!use24Hour ? 'text-white' : 'text-[#5A655C]'"
                  @click="setTimeFormat(false)"
                >
                  12h
                </button>
                <button
                  type="button"
                  class="relative z-10 px-2.5 py-1 rounded-full font-medium transition-colors cursor-pointer"
                  :class="use24Hour ? 'text-white' : 'text-[#5A655C]'"
                  @click="setTimeFormat(true)"
                >
                  24h
                </button>
              </div>
            </div>

            <div
              v-if="isLoadingSlots"
              class="flex-1 flex items-center justify-center text-sm text-[#5A655C] py-8"
            >
              Loading slots…
            </div>

            <div
              v-else-if="availableSlotsForDay.length === 0"
              class="flex-1 flex items-center justify-center text-sm text-[#5A655C] text-center py-8"
            >
              No open slots on this day. Pick another green date.
            </div>

            <div v-else ref="slotsListRef" class="space-y-2 max-h-[260px] overflow-y-auto overflow-x-hidden pr-1">
              <button
                v-for="slot in availableSlotsForDay"
                :key="slot.time"
                type="button"
                class="slot-pill-anim w-full rounded-xl border px-3.5 py-2.5 text-sm font-medium flex items-center gap-2 transition-colors cursor-pointer"
                :class="
                  selectedSlot === slot.time
                    ? 'bg-[#1A241B] text-white border-[#1A241B]'
                    : 'bg-white text-[#141A15] border-[#E2E6E0] hover:border-[#34D361] hover:bg-[#E8EFE6]/45'
                "
                @click="selectedSlot = slot.time"
              >
                <span
                  class="w-1.5 h-1.5 rounded-full"
                  :class="selectedSlot === slot.time ? 'bg-[#4ADE80]' : 'bg-[#22C55E]'"
                />
                <span class="slot-time-label inline-block">{{ formatSlotTime(slot.time) }}</span>
              </button>
            </div>

            <!-- Inline Booking Confirm Drawer -->
            <div
              v-if="selectedSlot"
              ref="confirmDrawerRef"
              class="mt-4 pt-4 border-t border-[#E2E6E0] space-y-2.5 overflow-hidden"
            >
              <input
                v-model="customerName"
                type="text"
                placeholder="Your name (optional)"
                class="confirm-field w-full rounded-xl bg-[#F7F7F3] border border-[#E2E6E0] px-3 py-2 text-xs text-[#141A15] focus:outline-none focus:border-[#1A241B]"
              />
              <input
                v-model="customerEmail"
                type="email"
                placeholder="Email for calendar invite"
                class="confirm-field w-full rounded-xl bg-[#F7F7F3] border border-[#E2E6E0] px-3 py-2 text-xs text-[#141A15] focus:outline-none focus:border-[#1A241B]"
              />
              <input
                v-model="customerPhone"
                type="tel"
                placeholder="Phone (optional)"
                class="confirm-field w-full rounded-xl bg-[#F7F7F3] border border-[#E2E6E0] px-3 py-2 text-xs text-[#141A15] focus:outline-none focus:border-[#1A241B]"
              />

              <p v-if="bookingError" class="text-xs text-rose-600">{{ bookingError }}</p>

              <button
                type="button"
                :disabled="isSubmitting"
                class="book-cta-btn confirm-field w-full rounded-full bg-[#34D361] text-[#141A15] py-2.5 px-4 text-xs font-semibold hover:bg-[#4ADE80] transition-colors disabled:opacity-50 cursor-pointer"
                @click="confirmSlotBooking"
              >
                {{ isSubmitting ? 'Reserving…' : `Book ${formatSlotTime(selectedSlot)}` }}
              </button>
            </div>
          </template>
        </div>
      </div>

      <!-- Financing Trust Bar Below Widget -->
      <div
        class="financing-trust-bar mt-8 flex flex-wrap items-center justify-center gap-x-5 gap-y-2 text-sm text-[#5A655C]"
      >
        <span class="trust-item-anim"><strong class="font-semibold text-[#141A15]">$0</strong> down payment</span>
        <span class="trust-dot-sep inline-block" aria-hidden="true">·</span>
        <span class="trust-item-anim"><strong class="font-semibold text-[#141A15]">0%</strong> interest</span>
        <span class="trust-dot-sep inline-block" aria-hidden="true">·</span>
        <span class="trust-item-anim">insurance <strong class="font-semibold text-[#141A15]">billed direct</strong></span>
        <span class="trust-dot-sep inline-block" aria-hidden="true">·</span>
        <span class="trust-item-anim"><strong ref="wayCountRef" class="font-semibold text-[#141A15]">3</strong> ways to pay</span>
      </div>
    </div>
  </section>
</template>
