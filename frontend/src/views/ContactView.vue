<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { CheckCircle2, ChevronDown, Mail, Phone } from 'lucide-vue-next';
import CopelyLogo from '@/components/ui/CopelyLogo.vue';
import { SERVICE_OPTIONS, TIMELINE_OPTIONS } from '@/data/copelyData';
import { useCopelyContent } from '@/composables/useCopelyContent';
import { useEstimateDraft } from '@/composables/useEstimateDraft';
import { useGsapScope } from '@/composables/useGsapScope';
import { submitEstimate } from '@/services/api';
import type { EstimateResponse, ProjectTimeline, TileServiceCategory } from '@/types/copely';

const { content, loadContent } = useCopelyContent();
const { estimateDraft } = useEstimateDraft();
const pageRef = ref<HTMLElement | null>(null);

useGsapScope(pageRef, ({ gsap }) => {
  gsap.from('.contact-left-anim', { y: 24, opacity: 0, duration: 0.65, ease: 'power3.out' });
  gsap.from('.contact-card-anim', {
    y: 30,
    scale: 0.98,
    opacity: 0,
    duration: 0.7,
    delay: 0.12,
    ease: 'power3.out',
  });
});

const form = reactive<{
  firstName: string;
  lastName: string;
  phone: string;
  email: string;
  address: string;
  town: string;
  service: TileServiceCategory;
  timeline: ProjectTimeline;
  notes: string;
}>({
  firstName: '',
  lastName: '',
  phone: '',
  email: '',
  address: '',
  town: '',
  service: estimateDraft.service,
  timeline: estimateDraft.timeline,
  notes: '',
});

const isSubmitting = ref<boolean>(false);
const errorMessage = ref<string>('');
const confirmation = ref<EstimateResponse | null>(null);

onMounted(() => {
  void loadContent();
});

async function handleSubmit(): Promise<void> {
  errorMessage.value = '';
  if (!form.firstName.trim() || !form.lastName.trim()) {
    errorMessage.value = 'Please enter your first and last name.';
    return;
  }
  if (!form.phone.trim() && !form.email.trim()) {
    errorMessage.value = 'Please include a phone number or email address so we can reply.';
    return;
  }

  isSubmitting.value = true;
  try {
    confirmation.value = await submitEstimate({
      firstName: form.firstName.trim(),
      lastName: form.lastName.trim(),
      phone: form.phone.trim(),
      email: form.email.trim(),
      address: form.address.trim(),
      town: form.town,
      service: form.service,
      timeline: form.timeline,
      notes: form.notes.trim() || undefined,
    });
  } catch (err) {
    errorMessage.value =
      err instanceof Error ? err.message : 'Could not send estimate request. Please call us.';
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<template>
  <section ref="pageRef" class="min-h-screen bg-[#1A241B] text-white pt-32 pb-20 px-6 md:px-10">
    <div class="max-w-[1280px] mx-auto grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16 items-start">
      <!-- Left Column -->
      <div class="contact-left-anim lg:col-span-5 space-y-6 pt-2">
        <span class="block text-xs font-semibold uppercase tracking-[0.18em] text-[#4ADE80]">
          CONTACT US
        </span>

        <h1 class="text-4xl sm:text-5xl md:text-6xl font-medium tracking-[-0.03em] leading-[1.06]">
          Free estimate.<br />
          <span class="inline-flex items-center gap-3 mt-1">
            <CopelyLogo :show-wordmark="false" :icon-size="34" />
            <span class="font-serif italic font-normal text-white/95">We’ll follow up.</span>
          </span>
        </h1>

        <p class="text-base md:text-lg text-white/75 leading-relaxed max-w-md">
          Fill out the form and we’ll get back to you within 24 hours. Active leak? Water damage?
          Call
          <a href="tel:+15125550142" class="text-[#4ADE80] font-semibold underline">
            (512) 555-0142
          </a>
          and we’ll move fast.
        </p>

        <!-- Contact Badges with Green Circular Icons -->
        <div class="space-y-4 pt-4">
          <a
            href="tel:+15125550142"
            class="flex items-center gap-4 rounded-2xl bg-[#253227]/70 border border-white/10 p-4 hover:border-[#4ADE80]/40 transition"
          >
            <span
              class="w-11 h-11 rounded-full bg-[#34D361] text-[#141A15] flex items-center justify-center shrink-0"
            >
              <Phone class="w-5 h-5" />
            </span>
            <div>
              <span class="block text-xs text-white/60">Call or text</span>
              <span class="block text-base font-semibold text-white">(512) 555-0142</span>
            </div>
          </a>

          <a
            href="mailto:hello@copelytiling.com"
            class="flex items-center gap-4 rounded-2xl bg-[#253227]/70 border border-white/10 p-4 hover:border-[#4ADE80]/40 transition"
          >
            <span
              class="w-11 h-11 rounded-full bg-[#34D361] text-[#141A15] flex items-center justify-center shrink-0"
            >
              <Mail class="w-5 h-5" />
            </span>
            <div>
              <span class="block text-xs text-white/60">Email</span>
              <span class="block text-base font-semibold text-white">hello@copelytiling.com</span>
            </div>
          </a>
        </div>
      </div>

      <!-- Right Column: Tell us about your project Form -->
      <div class="lg:col-span-7">
        <div class="contact-card-anim rounded-3xl bg-[#212D22] border border-white/10 p-7 md:p-10 shadow-2xl">
          <!-- Confirmation State -->
          <div v-if="confirmation" class="py-10 text-center space-y-4">
            <CheckCircle2 class="w-12 h-12 text-[#34D361] mx-auto" />
            <span class="block text-xs font-semibold uppercase tracking-widest text-[#4ADE80]">
              Request Received · {{ confirmation.estimateId }}
            </span>
            <h2 class="text-3xl font-medium text-white">We’ve got your project details.</h2>
            <p class="text-white/75 max-w-md mx-auto">
              Expect a follow-up from our Hill Country estimator
              <strong class="text-white">{{ confirmation.responseWindow.toLowerCase() }}</strong>.
            </p>
            <button
              type="button"
              class="mt-4 rounded-full bg-white text-[#141A15] px-6 py-3 text-sm font-semibold hover:bg-[#E8EFE6] transition cursor-pointer"
              @click="confirmation = null"
            >
              Send another request
            </button>
          </div>

          <!-- Form -->
          <form v-else class="space-y-5" @submit.prevent="handleSubmit">
            <h2 class="text-2xl font-medium text-white tracking-tight">
              Tell us about your project
            </h2>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label for="contact-first" class="block text-xs text-white/65 mb-1.5">
                  First name
                </label>
                <input
                  id="contact-first"
                  v-model="form.firstName"
                  type="text"
                  placeholder="Sam"
                  required
                  class="w-full rounded-xl bg-[#253227] border border-white/15 px-4 py-3 text-sm text-white placeholder:text-white/35 focus:outline-none focus:border-[#4ADE80]"
                />
              </div>

              <div>
                <label for="contact-last" class="block text-xs text-white/65 mb-1.5">
                  Last name
                </label>
                <input
                  id="contact-last"
                  v-model="form.lastName"
                  type="text"
                  placeholder="Calloway"
                  required
                  class="w-full rounded-xl bg-[#253227] border border-white/15 px-4 py-3 text-sm text-white placeholder:text-white/35 focus:outline-none focus:border-[#4ADE80]"
                />
              </div>

              <div>
                <label for="contact-phone" class="block text-xs text-white/65 mb-1.5">Phone</label>
                <input
                  id="contact-phone"
                  v-model="form.phone"
                  type="tel"
                  placeholder="(512) 555-0100"
                  class="w-full rounded-xl bg-[#253227] border border-white/15 px-4 py-3 text-sm text-white placeholder:text-white/35 focus:outline-none focus:border-[#4ADE80]"
                />
              </div>

              <div>
                <label for="contact-email" class="block text-xs text-white/65 mb-1.5">Email</label>
                <input
                  id="contact-email"
                  v-model="form.email"
                  type="email"
                  placeholder="you@example.com"
                  class="w-full rounded-xl bg-[#253227] border border-white/15 px-4 py-3 text-sm text-white placeholder:text-white/35 focus:outline-none focus:border-[#4ADE80]"
                />
              </div>

              <div>
                <label for="contact-address" class="block text-xs text-white/65 mb-1.5">
                  Home address
                </label>
                <input
                  id="contact-address"
                  v-model="form.address"
                  type="text"
                  placeholder="123 Ridgeline Dr"
                  class="w-full rounded-xl bg-[#253227] border border-white/15 px-4 py-3 text-sm text-white placeholder:text-white/35 focus:outline-none focus:border-[#4ADE80]"
                />
              </div>

              <div>
                <label for="contact-town" class="block text-xs text-white/65 mb-1.5">Town</label>
                <div class="relative">
                  <select
                    id="contact-town"
                    v-model="form.town"
                    class="w-full appearance-none rounded-xl bg-[#253227] border border-white/15 px-4 py-3 pr-10 text-sm text-white focus:outline-none focus:border-[#4ADE80] cursor-pointer"
                  >
                    <option value="">Choose your town</option>
                    <option
                      v-for="town in content.towns"
                      :key="town"
                      :value="town"
                      class="bg-[#1A241B] text-white"
                    >
                      {{ town }}
                    </option>
                  </select>
                  <ChevronDown
                    class="w-4 h-4 text-white/50 absolute right-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
                  />
                </div>
              </div>

              <div>
                <label for="contact-service" class="block text-xs text-white/65 mb-1.5">
                  Service
                </label>
                <div class="relative">
                  <select
                    id="contact-service"
                    v-model="form.service"
                    class="w-full appearance-none rounded-xl bg-[#253227] border border-white/15 px-4 py-3 pr-10 text-sm text-white focus:outline-none focus:border-[#4ADE80] cursor-pointer"
                  >
                    <option
                      v-for="opt in SERVICE_OPTIONS"
                      :key="opt.value"
                      :value="opt.value"
                      class="bg-[#1A241B] text-white"
                    >
                      {{ opt.label }}
                    </option>
                  </select>
                  <ChevronDown
                    class="w-4 h-4 text-white/50 absolute right-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
                  />
                </div>
              </div>

              <div>
                <label for="contact-timeline" class="block text-xs text-white/65 mb-1.5">
                  Timeline
                </label>
                <div class="relative">
                  <select
                    id="contact-timeline"
                    v-model="form.timeline"
                    class="w-full appearance-none rounded-xl bg-[#253227] border border-white/15 px-4 py-3 pr-10 text-sm text-white focus:outline-none focus:border-[#4ADE80] cursor-pointer"
                  >
                    <option
                      v-for="opt in TIMELINE_OPTIONS"
                      :key="opt.value"
                      :value="opt.value"
                      class="bg-[#1A241B] text-white"
                    >
                      {{ opt.label }}
                    </option>
                  </select>
                  <ChevronDown
                    class="w-4 h-4 text-white/50 absolute right-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
                  />
                </div>
              </div>
            </div>

            <div>
              <label for="contact-notes" class="block text-xs text-white/65 mb-1.5">
                Project notes (optional)
              </label>
              <textarea
                id="contact-notes"
                v-model="form.notes"
                rows="3"
                placeholder="Approximate sq ft, room type, or active shower leak details…"
                class="w-full rounded-xl bg-[#253227] border border-white/15 px-4 py-3 text-sm text-white placeholder:text-white/35 focus:outline-none focus:border-[#4ADE80]"
              />
            </div>

            <p v-if="errorMessage" class="text-sm text-rose-400">{{ errorMessage }}</p>

            <button
              type="submit"
              :disabled="isSubmitting"
              class="w-full rounded-full bg-[#34D361] text-[#141A15] py-3.5 px-6 text-sm font-semibold hover:bg-[#4ADE80] transition-colors disabled:opacity-50 cursor-pointer"
            >
              {{ isSubmitting ? 'Sending request…' : 'Get My Free Estimate' }}
            </button>
          </form>
        </div>
      </div>
    </div>
  </section>
</template>
