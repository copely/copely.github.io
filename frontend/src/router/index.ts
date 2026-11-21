import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { allowsMotion } from '@/composables/useGsapScope';
import { useSiteExplorer } from '@/composables/useSiteExplorer';
import HomeView from '@/views/HomeView.vue';
import ServicesView from '@/views/ServicesView.vue';
import ServiceDetailView from '@/views/ServiceDetailView.vue';
import AboutView from '@/views/AboutView.vue';
import BlogView from '@/views/BlogView.vue';
import BlogDetailView from '@/views/BlogDetailView.vue';
import FaqView from '@/views/FaqView.vue';
import ContactView from '@/views/ContactView.vue';
import LegalView from '@/views/LegalView.vue';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
    meta: { headerTone: 'dark-hero' },
  },
  {
    path: '/services',
    name: 'services',
    component: ServicesView,
    meta: { headerTone: 'light' },
  },
  {
    path: '/services/:slug',
    name: 'service-detail',
    component: ServiceDetailView,
    meta: { headerTone: 'light' },
  },
  {
    path: '/about',
    name: 'about',
    component: AboutView,
    meta: { headerTone: 'light' },
  },
  {
    path: '/blog',
    name: 'blog',
    component: BlogView,
    meta: { headerTone: 'light' },
  },
  {
    path: '/blog/:slug',
    name: 'blog-detail',
    component: BlogDetailView,
    meta: { headerTone: 'light' },
  },
  {
    path: '/faq',
    name: 'faq',
    component: FaqView,
    meta: { headerTone: 'light' },
  },
  {
    path: '/contact',
    name: 'contact',
    component: ContactView,
    meta: { headerTone: 'dark-hero' },
  },
  {
    path: '/privacy',
    name: 'privacy',
    component: LegalView,
    props: { docType: 'privacy' },
    meta: { headerTone: 'light' },
  },
  {
    path: '/terms',
    name: 'terms',
    component: LegalView,
    props: { docType: 'terms' },
    meta: { headerTone: 'light' },
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (from.matched.length > 0 && to.path !== from.path && allowsMotion()) {
      const { setPendingRouteScroll } = useSiteExplorer();
      if (savedPosition) {
        setPendingRouteScroll({ top: savedPosition.top, left: savedPosition.left });
      } else if (to.hash) {
        setPendingRouteScroll({ hash: to.hash });
      } else {
        setPendingRouteScroll({ top: 0, left: 0 });
      }
      return false;
    }
    if (savedPosition) {
      return savedPosition;
    }
    if (to.hash) {
      return {
        el: to.hash,
        behavior: allowsMotion() ? 'smooth' : 'auto',
        top: 96,
      };
    }
    return { top: 0, behavior: 'auto' };
  },
});

export default router;
