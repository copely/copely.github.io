<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { ArrowUpRight } from 'lucide-vue-next';
import { useCopelyContent } from '@/composables/useCopelyContent';
import { useGsapScope } from '@/composables/useGsapScope';
import type { BlogPost } from '@/types/copely';

const { content, loadContent } = useCopelyContent();
const pageRef = ref<HTMLElement | null>(null);

useGsapScope(pageRef, ({ gsap }) => {
  gsap.from('.blog-header-anim', { y: 22, opacity: 0, duration: 0.6, ease: 'power3.out' });
  gsap.from('.blog-featured-anim', { y: 26, opacity: 0, duration: 0.65, ease: 'power3.out', delay: 0.1 });
  gsap.from('.blog-card-anim', {
    y: 26,
    opacity: 0,
    duration: 0.55,
    stagger: 0.08,
    ease: 'power3.out',
    scrollTrigger: { trigger: '.blog-grid-wrap', start: 'top 84%', once: true },
  });
});

onMounted(() => {
  void loadContent();
});

const featuredPost = computed<BlogPost | undefined>(() => {
  return content.value.blogPosts.find((p) => p.featured) ?? content.value.blogPosts[0];
});

const gridPosts = computed<BlogPost[]>(() => {
  const featId = featuredPost.value?.id;
  return content.value.blogPosts.filter((p) => p.id !== featId);
});
</script>

<template>
  <section ref="pageRef" class="pt-32 pb-16 px-6 md:px-10 bg-[#F7F7F3]">
    <div class="max-w-[1320px] mx-auto space-y-14">
      <!-- Header Split -->
      <div class="blog-header-anim grid grid-cols-1 lg:grid-cols-12 gap-8 items-end pb-10 border-b border-[#E2E6E0]">
        <div class="lg:col-span-7">
          <h1 class="text-4xl sm:text-5xl md:text-6xl font-medium tracking-[-0.03em] leading-[1.06] text-[#141A15]">
            Tiling advice <span class="font-serif italic font-normal">you can</span><br />
            actually use.
          </h1>
        </div>
        <div class="lg:col-span-5">
          <p class="text-base md:text-lg text-[#5A655C] leading-relaxed">
            Costs, waterproofing, stone care, and how tile systems actually work — written by the
            crew that lays them.
          </p>
        </div>
      </div>

      <!-- Featured Article Banner (Large 2-Column Horizontal Card) -->
      <RouterLink
        v-if="featuredPost"
        :to="`/blog/${featuredPost.slug}`"
        class="blog-featured-anim group rounded-3xl bg-white border border-[#E2E6E0] overflow-hidden grid grid-cols-1 lg:grid-cols-12 shadow-sm hover:shadow-xl transition-shadow duration-300"
      >
        <div class="lg:col-span-7 aspect-[16/10] lg:aspect-auto overflow-hidden bg-[#1A241B]">
          <img
            :src="featuredPost.coverImage"
            :alt="featuredPost.title"
            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
        </div>

        <div class="lg:col-span-5 p-8 md:p-10 flex flex-col justify-between gap-6">
          <div class="space-y-4">
            <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-[0.14em] text-[#22C55E]">
              <span>{{ featuredPost.category }}</span>
              <span class="text-[#5A655C]">·</span>
              <span class="text-[#5A655C]">{{ featuredPost.publishedAt }}</span>
              <span class="text-[#5A655C]">·</span>
              <span class="text-[#5A655C]">{{ featuredPost.readTime }}</span>
            </div>

            <h2 class="text-2xl sm:text-3xl font-medium tracking-tight text-[#141A15] group-hover:text-[#22C55E] transition-colors">
              {{ featuredPost.title }}
            </h2>

            <p class="text-sm md:text-base text-[#5A655C] leading-relaxed">
              {{ featuredPost.excerpt }}
            </p>
          </div>

          <div class="flex items-center justify-between pt-4 border-t border-[#E2E6E0]">
            <div class="flex items-center gap-3">
              <img
                :src="featuredPost.authorAvatar"
                :alt="featuredPost.authorName"
                class="w-10 h-10 rounded-full object-cover"
              />
              <div>
                <p class="text-sm font-semibold text-[#141A15]">{{ featuredPost.authorName }}</p>
                <p class="text-xs text-[#5A655C]">{{ featuredPost.authorRole }}</p>
              </div>
            </div>
            <span class="inline-flex items-center gap-1 text-sm font-medium text-[#141A15]">
              <span>Read</span>
              <ArrowUpRight class="w-4 h-4" />
            </span>
          </div>
        </div>
      </RouterLink>

      <!-- 3-Column Article Grid -->
      <div class="blog-grid-wrap grid grid-cols-1 md:grid-cols-3 gap-7">
        <RouterLink
          v-for="post in gridPosts"
          :key="post.id"
          :to="`/blog/${post.slug}`"
          class="blog-card-anim group rounded-3xl bg-white border border-[#E2E6E0] overflow-hidden flex flex-col justify-between shadow-sm hover:shadow-xl transition-shadow duration-300"
        >
          <div>
            <div class="aspect-[4/3] overflow-hidden bg-[#1A241B]">
              <img
                :src="post.coverImage"
                :alt="post.title"
                loading="lazy"
                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
              />
            </div>

            <div class="p-6 space-y-3">
              <div class="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-[0.14em] text-[#5A655C]">
                <span class="text-[#22C55E]">{{ post.category }}</span>
                <span>·</span>
                <span>{{ post.publishedAt }}</span>
                <span>·</span>
                <span>{{ post.readTime }}</span>
              </div>

              <h3 class="text-xl font-medium text-[#141A15] tracking-tight group-hover:text-[#22C55E] transition-colors">
                {{ post.title }}
              </h3>

              <p class="text-sm text-[#5A655C] leading-relaxed">
                {{ post.excerpt }}
              </p>
            </div>
          </div>

          <div class="px-6 pb-6 pt-3 flex items-center gap-3">
            <img
              :src="post.authorAvatar"
              :alt="post.authorName"
              class="w-8 h-8 rounded-full object-cover"
            />
            <div class="text-xs">
              <p class="font-semibold text-[#141A15]">{{ post.authorName }}</p>
              <p class="text-[#5A655C]">{{ post.authorRole }}</p>
            </div>
          </div>
        </RouterLink>
      </div>
    </div>
  </section>
</template>
