<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ArrowLeft } from 'lucide-vue-next';
import { useCopelyContent } from '@/composables/useCopelyContent';
import { useGsapScope } from '@/composables/useGsapScope';
import { fetchPostBySlug } from '@/services/api';
import type { BlogPost } from '@/types/copely';

const route = useRoute();
const router = useRouter();
const { content, loadContent } = useCopelyContent();
const articleRef = ref<HTMLElement | null>(null);

useGsapScope(articleRef, ({ gsap }) => {
  gsap.from('.blog-detail-reveal', { y: 22, opacity: 0, duration: 0.6, stagger: 0.08, ease: 'power3.out' });
});

const fetchedPost = ref<BlogPost | null>(null);
const slug = computed<string>(() => String(route.params.slug ?? ''));

const post = computed<BlogPost | undefined>(() => {
  if (fetchedPost.value && fetchedPost.value.slug === slug.value) {
    return fetchedPost.value;
  }
  return content.value.blogPosts.find((p) => p.slug === slug.value) ?? content.value.blogPosts[0];
});

async function loadPost(): Promise<void> {
  await loadContent();
  try {
    fetchedPost.value = await fetchPostBySlug(slug.value);
  } catch {
    fetchedPost.value = null;
  }
}

function scrollToBooking(): void {
  void router.push({ path: '/', hash: '#schedule' });
}

watch(slug, () => {
  void loadPost();
});

onMounted(() => {
  void loadPost();
});
</script>

<template>
  <article v-if="post" ref="articleRef" class="pt-32 pb-16 px-6 md:px-10 bg-[#F7F7F3]">
    <div class="max-w-[860px] mx-auto space-y-10">
      <RouterLink
        to="/blog"
        class="inline-flex items-center gap-2 text-sm font-medium text-[#5A655C] hover:text-[#141A15] transition"
      >
        <ArrowLeft class="w-4 h-4" />
        <span>Back to all guides</span>
      </RouterLink>

      <header class="blog-detail-reveal space-y-5">
        <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-[0.14em] text-[#5A655C]">
          <span class="text-[#22C55E]">{{ post.category }}</span>
          <span>·</span>
          <span>{{ post.publishedAt }}</span>
          <span>·</span>
          <span>{{ post.readTime }}</span>
        </div>

        <h1 class="text-3xl sm:text-4xl md:text-5xl font-medium tracking-[-0.03em] leading-[1.1] text-[#141A15]">
          {{ post.title }}
        </h1>

        <p class="text-lg text-[#5A655C] leading-relaxed">
          {{ post.excerpt }}
        </p>

        <div class="flex items-center gap-3 pt-2">
          <img
            :src="post.authorAvatar"
            :alt="post.authorName"
            class="w-11 h-11 rounded-full object-cover"
          />
          <div>
            <p class="text-sm font-semibold text-[#141A15]">{{ post.authorName }}</p>
            <p class="text-xs text-[#5A655C]">{{ post.authorRole }}</p>
          </div>
        </div>
      </header>

      <div class="blog-detail-reveal rounded-3xl overflow-hidden bg-[#1A241B] aspect-[16/9] shadow-lg">
        <img :src="post.coverImage" :alt="post.title" class="w-full h-full object-cover" />
      </div>

      <div class="space-y-8 pt-2">
        <section v-for="sec in post.sections" :key="sec.heading" class="space-y-3">
          <h2 class="text-2xl font-medium tracking-tight text-[#141A15]">{{ sec.heading }}</h2>
          <p
            v-for="(para, pIdx) in sec.paragraphs"
            :key="pIdx"
            class="text-base md:text-lg text-[#364137] leading-relaxed"
          >
            {{ para }}
          </p>
        </section>
      </div>

      <div class="rounded-3xl bg-[#1A241B] text-white p-8 flex flex-col sm:flex-row sm:items-center justify-between gap-6">
        <div>
          <h3 class="text-xl font-medium">Want an honest look at your tile project?</h3>
          <p class="text-sm text-white/75 mt-1">
            Book a free 30-minute laser measure &amp; moisture inspection in Texas Hill Country.
          </p>
        </div>
        <button
          type="button"
          class="shrink-0 rounded-full bg-[#34D361] text-[#141A15] px-6 py-3 text-sm font-semibold hover:bg-[#4ADE80] transition cursor-pointer"
          @click="scrollToBooking"
        >
          Get a free estimate
        </button>
      </div>
    </div>
  </article>
</template>
