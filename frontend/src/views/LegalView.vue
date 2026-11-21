<script setup lang="ts">
interface Props {
  docType: 'privacy' | 'terms';
}

defineProps<Props>();
</script>

<template>
  <section class="pt-32 pb-16 px-6 md:px-10 bg-[#F7F7F3]">
    <div class="max-w-[820px] mx-auto space-y-6">
      <span class="text-xs font-semibold uppercase tracking-[0.18em] text-[#22C55E]">
        COPELY TILING · LEGAL
      </span>
      <h1 class="text-4xl md:text-5xl font-medium tracking-[-0.03em] text-[#141A15]">
        {{ docType === 'privacy' ? 'Privacy Policy' : 'Terms of Service' }}
      </h1>

      <div
        v-if="docType === 'privacy'"
        class="rounded-3xl bg-white border border-[#E2E6E0] p-8 space-y-4 text-sm md:text-base text-[#364137] leading-relaxed"
      >
        <p>
          Copely collects only the contact details, address information, and project notes you
          voluntarily submit when requesting a free estimate or scheduling an on-site visit in the
          Texas Hill Country.
        </p>
        <p>
          We never sell or broker homeowner information to lead aggregators. Your project details
          are used strictly to prepare written estimates, schedule site measurements, and document
          your 25-year written installation warranty.
        </p>
      </div>

      <div
        v-else
        class="rounded-3xl bg-white border border-[#E2E6E0] p-8 space-y-4 text-sm md:text-base text-[#364137] leading-relaxed"
      >
        <p>
          On-site tile measurements and written estimates provided by Copely carry no obligation.
          Every installed tile and waterproofing scope is governed by the finalized written work
          order signed prior to project start.
        </p>
        <p>
          Our 25-year transferable installation warranty covers substrate preparation, bonded
          waterproofing membranes, mortar bond integrity, and structural grout cracking under normal
          residential use.
        </p>
      </div>
    </div>
  </section>
</template>
