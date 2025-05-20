// router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import Home from '../pages/Home.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/guide',
    name: 'Guide',
    component: () => import('../pages/Guide.vue'),
  },
  {
    path: '/about/history',
    name: 'History',
    component: () => import('../pages/about/History.vue'),
  },
  {
    path: '/about/mission-vision',
    name: 'Mission',
    component: () => import('../pages/about/MissionVision.vue'),
  },
  {
    path: '/support/help',
    name: 'Help',
    component: () => import('../pages/support/HelpFAQ.vue'),
  },
  {
    path: '/support/contact',
    name: 'Contact',
    component: () => import('../pages/support/ContactForm.vue'),
  },
  {
    path: '/support/feedback',
    name: 'Feedback',
    component: () => import('../pages/support/Feedback.vue'),
  },
  {
    path: '/legal/privacy',
    name: 'Privacy',
    component: () => import('../pages/legal/PrivacyPolicy.vue'),
  },
  {
    path: '/legal/terms',
    name: 'Terms',
    component: () => import('../pages/legal/TermsConditions.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../pages/NotFound.vue'),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior () {
    // Siempre desplazarse al inicio cuando se cambia de página
    return { top: 0 }
  },
});

export default router;
