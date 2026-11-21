import { onMounted, onUnmounted, type Ref } from 'vue';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

let pluginRegistered = false;

export function allowsMotion(): boolean {
  return (
    typeof window !== 'undefined' &&
    window.matchMedia('(prefers-reduced-motion: no-preference)').matches
  );
}

export function ensureGsapRegistered(): void {
  if (!pluginRegistered && typeof window !== 'undefined') {
    gsap.registerPlugin(ScrollTrigger);
    pluginRegistered = true;
  }
}

export function useGsapScope(
  scopeRef: Ref<HTMLElement | null>,
  setup: (ctx: { gsap: typeof gsap; ScrollTrigger: typeof ScrollTrigger }) => void,
): void {
  let mm: gsap.MatchMedia | null = null;

  onMounted(() => {
    ensureGsapRegistered();
    if (!scopeRef.value) {
      return;
    }

    mm = gsap.matchMedia();
    mm.add('(prefers-reduced-motion: no-preference)', () => {
      if (!scopeRef.value) {
        return;
      }
      const ctx = gsap.context(() => {
        setup({ gsap, ScrollTrigger });
      }, scopeRef.value);

      return () => {
        ctx.revert();
      };
    });
  });

  onUnmounted(() => {
    mm?.revert();
    mm = null;
  });
}
