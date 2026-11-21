<script setup lang="ts">
import { ref, watch } from 'vue';
import gsap from 'gsap';
import { useGsapScope } from '@/composables/useGsapScope';
import { assetUrl } from '@/utils/assetUrl';

interface LayerItem {
  id: number;
  name: string;
  spec: string;
}

const layers: LayerItem[] = [
  {
    id: 5,
    name: 'Porcelain / Stone Tile & Edge Profile',
    spec: 'Kiln-fired porcelain or natural stone locked with anodized trim',
  },
  {
    id: 4,
    name: 'Ridged Thin-Set Mortar Bed',
    spec: '100% coverage polymer-modified mortar combed directionally',
  },
  {
    id: 3,
    name: 'Waterproof Uncoupling Membrane',
    spec: 'Absorbs slab movement and seals against subfloor moisture',
  },
  {
    id: 2,
    name: 'Cement Backer Board / Leveling Prep',
    spec: 'Laser-flattened substrate to < 1/8" over 10 feet',
  },
  {
    id: 1,
    name: 'Structural Subfloor Joists / Slab',
    spec: 'Checked for deflection and movement joints before setting',
  },
];

const activeLayerId = ref<number>(5);
const sectionRef = ref<HTMLElement | null>(null);
const svgRef = ref<SVGSVGElement | null>(null);

interface FeatureCard {
  title: string;
  description: string;
  imageUrl: string;
  imageAlt: string;
}

const featureCards: FeatureCard[] = [
  {
    title: 'Built to last',
    description:
      'Thick, heavy-duty porcelain and stone slabs built to take heavy foot traffic and impact.',
    imageUrl: assetUrl('/images/cutaway-slab.jpg'),
    imageAlt: 'Close-up of thick heavy-duty porcelain tile slab',
  },
  {
    title: 'Locked down tight',
    description:
      'Waterproof membranes and metal profiles seal every wall-to-floor seam against moisture and leaks.',
    imageUrl: assetUrl('/images/cutaway-membrane.jpg'),
    imageAlt: 'Waterproof membrane and aluminum profile flashing seam',
  },
  {
    title: 'Holds its line',
    description: 'A crisp, dead-straight grout line and zero lippage — true from wall to wall.',
    imageUrl: assetUrl('/images/cutaway-grout.jpg'),
    imageAlt: 'Precision laser-straight grout lines with zero lippage',
  },
];

useGsapScope(sectionRef, ({ gsap: g }) => {
  g.from('.cutaway-header', {
    y: 26,
    opacity: 0,
    duration: 0.65,
    stagger: 0.08,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: sectionRef.value,
      start: 'top 78%',
      once: true,
    },
  });

  // Exploded-to-locked 5-layer isometric SVG assembly
  const layerOffsets = [-6, -14, -22, -30, -38];
  const layerEls = sectionRef.value?.querySelectorAll('.iso-layer');
  if (layerEls && layerEls.length > 0) {
    layerEls.forEach((el, idx) => {
      g.from(el, {
        y: layerOffsets[idx] ?? -20,
        opacity: 0,
        duration: 0.65,
        delay: idx * 0.08,
        ease: 'power3.out',
        scrollTrigger: {
          trigger: '.iso-card',
          start: 'top 80%',
          once: true,
        },
      });
    });
  }

  g.from('.cutaway-feature-card', {
    x: 24,
    opacity: 0,
    duration: 0.6,
    stagger: 0.12,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.cutaway-right-panel',
      start: 'top 78%',
      once: true,
    },
  });
});

watch(activeLayerId, (activeId) => {
  if (!svgRef.value || window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    return;
  }
  const groups = svgRef.value.querySelectorAll<SVGGElement>('.iso-layer');
  groups.forEach((el) => {
    const id = Number(el.dataset.layerId);
    gsap.to(el, {
      y: id === activeId ? -5 : 0,
      duration: 0.25,
      ease: 'power2.out',
    });
  });
});
</script>

<template>
  <section ref="sectionRef" class="py-20 md:py-28 px-6 md:px-10 bg-[#F7F7F3]">
    <div class="max-w-[1320px] mx-auto grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-14 items-start">
      <!-- Left Column: Copy + Isometric 3D 5-Layer Cutaway -->
      <div class="lg:col-span-6 space-y-6">
        <h2
          class="cutaway-header text-4xl md:text-5xl font-medium tracking-[-0.03em] leading-[1.08] text-[#141A15]"
        >
          More than<br />
          surfaces.<br />
          <span class="font-serif italic font-normal text-[#1A241B]">It’s a system.</span>
        </h2>

        <p class="cutaway-header text-base text-[#5A655C] max-w-lg leading-relaxed">
          We build every floor and wet room as one engineered structure — subfloor and backer board
          up through uncoupling membrane, waterproofing, and kiln-fired tile — so it holds strong
          for decades, not just a few years.
        </p>

        <!-- Isometric 3D Cutaway Illustration Card -->
        <div class="iso-card rounded-3xl bg-white border border-[#E2E6E0] p-5 md:p-6 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <span class="text-xs font-semibold uppercase tracking-[0.14em] text-[#5A655C]">
              5-Layer Copely Assembly
            </span>
            <span class="text-xs font-medium text-[#22C55E]">Hover layer to inspect</span>
          </div>

          <svg
            ref="svgRef"
            viewBox="0 0 560 350"
            class="w-full h-auto"
            role="img"
            aria-label="Isometric 5-layer cutaway illustration of a Copely tile installation"
          >
            <defs>
              <linearGradient id="tileTopGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#354037" />
                <stop offset="100%" stop-color="#1E2720" />
              </linearGradient>
              <linearGradient id="membraneGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" stop-color="#4ADE80" />
                <stop offset="100%" stop-color="#22C55E" />
              </linearGradient>
            </defs>

            <!-- Layer 1: Subfloor Joists -->
            <g
              data-layer-id="1"
              class="iso-layer cursor-pointer transition-opacity duration-200"
              :opacity="activeLayerId === 1 ? 1 : 0.78"
              @mouseenter="activeLayerId = 1"
            >
              <polygon
                points="80,262 280,192 480,262 280,332"
                fill="#D6C8B4"
                stroke="#8B7A62"
                stroke-width="1.5"
              />
              <polygon points="80,262 280,332 280,344 80,274" fill="#B5A48B" />
              <polygon points="480,262 280,332 280,344 480,274" fill="#9D8B73" />
              <line x1="150" y1="238" x2="350" y2="308" stroke="#9D8B73" stroke-width="4" />
              <line x1="215" y1="215" x2="415" y2="285" stroke="#9D8B73" stroke-width="4" />
            </g>

            <!-- Layer 2: Cement Backer Board -->
            <g
              data-layer-id="2"
              class="iso-layer cursor-pointer transition-opacity duration-200"
              :opacity="activeLayerId === 2 ? 1 : 0.85"
              @mouseenter="activeLayerId = 2"
            >
              <polygon
                points="80,224 280,154 480,224 280,294"
                fill="#C7CBC6"
                stroke="#7E857D"
                stroke-width="1.5"
              />
              <polygon points="80,224 280,294 280,305 80,235" fill="#A5ABA4" />
              <polygon points="480,224 280,294 280,305 480,235" fill="#8C938B" />
            </g>

            <!-- Layer 3: Waterproof Uncoupling Membrane -->
            <g
              data-layer-id="3"
              class="iso-layer cursor-pointer transition-opacity duration-200"
              :opacity="activeLayerId === 3 ? 1 : 0.88"
              @mouseenter="activeLayerId = 3"
            >
              <polygon
                points="80,186 280,116 480,186 280,256"
                fill="url(#membraneGrad)"
                stroke="#15803D"
                stroke-width="1.5"
              />
              <polygon points="80,186 280,256 280,264 80,194" fill="#16A34A" />
              <polygon points="480,186 280,256 280,264 480,194" fill="#15803D" />
              <!-- Uncoupling waffle grid lines -->
              <line x1="145" y1="163" x2="345" y2="233" stroke="#DCFCE7" stroke-opacity="0.45" stroke-width="1.5" />
              <line x1="212" y1="140" x2="412" y2="210" stroke="#DCFCE7" stroke-opacity="0.45" stroke-width="1.5" />
              <line x1="145" y1="209" x2="345" y2="139" stroke="#DCFCE7" stroke-opacity="0.45" stroke-width="1.5" />
              <line x1="212" y1="232" x2="412" y2="162" stroke="#DCFCE7" stroke-opacity="0.45" stroke-width="1.5" />
            </g>

            <!-- Layer 4: Ridged Mortar Bed -->
            <g
              data-layer-id="4"
              class="iso-layer cursor-pointer transition-opacity duration-200"
              :opacity="activeLayerId === 4 ? 1 : 0.88"
              @mouseenter="activeLayerId = 4"
            >
              <polygon
                points="80,148 280,78 480,148 280,218"
                fill="#7B827B"
                stroke="#4F564F"
                stroke-width="1.5"
              />
              <polygon points="80,148 280,218 280,226 80,156" fill="#616861" />
              <polygon points="480,148 280,218 280,226 480,156" fill="#4E544E" />
              <!-- Trowel ridges -->
              <line x1="125" y1="133" x2="325" y2="203" stroke="#A6ADA6" stroke-width="2" stroke-dasharray="4 3" />
              <line x1="165" y1="119" x2="365" y2="189" stroke="#A6ADA6" stroke-width="2" stroke-dasharray="4 3" />
              <line x1="205" y1="105" x2="405" y2="175" stroke="#A6ADA6" stroke-width="2" stroke-dasharray="4 3" />
              <line x1="245" y1="91" x2="445" y2="161" stroke="#A6ADA6" stroke-width="2" stroke-dasharray="4 3" />
            </g>

            <!-- Layer 5: Porcelain / Stone Tile & Aluminum Edge Profile -->
            <g
              data-layer-id="5"
              class="iso-layer cursor-pointer transition-opacity duration-200"
              :opacity="activeLayerId === 5 ? 1 : 0.9"
              @mouseenter="activeLayerId = 5"
            >
              <polygon
                points="80,104 280,34 480,104 280,174"
                fill="url(#tileTopGrad)"
                stroke="#141A15"
                stroke-width="1.5"
              />
              <polygon points="80,104 280,174 280,185 80,115" fill="#2A332C" />
              <polygon points="480,104 280,174 280,185 480,115" fill="#171D18" />
              <!-- Grout cross joint on top tile -->
              <line x1="180" y1="69" x2="380" y2="139" stroke="#4ADE80" stroke-opacity="0.55" stroke-width="2" />
              <line x1="180" y1="139" x2="380" y2="69" stroke="#4ADE80" stroke-opacity="0.55" stroke-width="2" />
              <!-- Aluminum edge profile strip -->
              <polygon points="280,174 480,104 484,106 284,176" fill="#E5E7EB" />
            </g>

            <!-- Numbered Callout Badges -->
            <g v-for="(item, idx) in layers" :key="item.id">
              <circle
                :cx="44"
                :cy="104 + idx * 40"
                r="13"
                :fill="activeLayerId === item.id ? '#34D361' : '#1A241B'"
                class="cursor-pointer transition-colors"
                @mouseenter="activeLayerId = item.id"
              />
              <text
                x="44"
                :y="108 + idx * 40"
                text-anchor="middle"
                font-size="11"
                font-weight="600"
                :fill="activeLayerId === item.id ? '#141A15' : '#FFFFFF'"
                class="pointer-events-none"
              >
                {{ item.id }}
              </text>
            </g>
          </svg>

          <!-- Interactive Layer Legend Pill -->
          <div class="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-2">
            <button
              v-for="layer in layers"
              :key="layer.id"
              type="button"
              class="text-left rounded-xl px-3 py-2 border transition-all text-xs cursor-pointer"
              :class="
                activeLayerId === layer.id
                  ? 'bg-[#E8EFE6] border-[#34D361] text-[#141A15]'
                  : 'bg-[#F7F7F3] border-transparent text-[#5A655C] hover:border-[#E2E6E0]'
              "
              @mouseenter="activeLayerId = layer.id"
              @click="activeLayerId = layer.id"
            >
              <span class="font-semibold text-[#141A15]">{{ layer.id }}. {{ layer.name }}</span>
              <span class="block text-[11px] text-[#5A655C] mt-0.5 leading-snug">
                {{ layer.spec }}
              </span>
            </button>
          </div>
        </div>
      </div>

      <!-- Right Column: Dark Container (#1A241B) with 3 Stacked Feature Cards -->
      <div
        class="cutaway-right-panel lg:col-span-6 rounded-3xl bg-[#1A241B] p-6 md:p-8 text-white space-y-5 shadow-xl"
      >
        <div
          v-for="card in featureCards"
          :key="card.title"
          class="cutaway-feature-card group rounded-2xl bg-[#253227]/75 border border-white/10 p-4 md:p-5 flex flex-col sm:flex-row items-start sm:items-center gap-5 transition-colors hover:border-[#4ADE80]/40"
        >
          <div class="w-full sm:w-28 h-28 rounded-xl overflow-hidden shrink-0 bg-[#141A15]">
            <img
              :src="card.imageUrl"
              :alt="card.imageAlt"
              loading="lazy"
              class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            />
          </div>
          <div class="space-y-1.5">
            <h3 class="text-xl font-medium text-white tracking-tight">{{ card.title }}</h3>
            <p class="text-sm text-white/75 leading-relaxed">{{ card.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
