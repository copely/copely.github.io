<script setup lang="ts">
interface Props {
  showWordmark?: boolean;
  wordmarkClass?: string;
  iconSize?: number;
}

withDefaults(defineProps<Props>(), {
  showWordmark: true,
  wordmarkClass: 'text-white',
  iconSize: 24,
});
</script>

<template>
  <span class="inline-flex items-center gap-2.5 select-none">
    <svg
      :width="iconSize"
      :height="iconSize"
      viewBox="0 0 28 28"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden="true"
      class="shrink-0"
    >
      <!-- Geometric double-chevron / tile-peak brand mark -->
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
      v-if="showWordmark"
      class="font-sans font-semibold tracking-tight text-xl leading-none"
      :class="wordmarkClass"
    >
      Copely
    </span>
  </span>
</template>
