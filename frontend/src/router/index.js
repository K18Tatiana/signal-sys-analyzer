import { createRouter, createWebHistory } from 'vue-router';
import Home from '../pages/Home.vue';
import authService from '../services/auth.service';

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
  // Rutas de autenticación
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/user/auth/Login.vue'),
    meta: {
      redirectIfAuth: true, // Redirige si ya está autenticado
    },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../pages/user/auth/Register.vue'),
    meta: {
      redirectIfAuth: true, // Redirige si ya está autenticado
    },
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('../pages/user/auth/ForgotPassword.vue'),
    meta: {
      redirectIfAuth: true,
    },
  },
  // Rutas protegidas (requieren autenticación)
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../pages/user/Dashboard.vue'),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../pages/user/profile/Profile.vue'),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: '/user/documents',
    name: 'UserDocuments',
    component: () => import('../pages/user/Documents.vue'),
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: '/user/analysis',
    name: 'UserAnalysis',
    component: () => import('../pages/user/Analysis.vue'),
    meta: {
      requiresAuth: true,
    },
  },
  // Página 404
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

// Navegación global antes de cada ruta
router.beforeEach((to, from, next) => {
  // Verificar si la ruta requiere autenticación
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!authService.isAuthenticated()) {
      // Si no está autenticado, redirigir a Login con la URL anterior para regresar después
      next({
        path: '/login',
        query: { redirect: to.fullPath },
      });
    } else {
      // Si está autenticado, continuar normalmente
      next();
    }
  }
  // Verificar si debe redirigir porque ya está autenticado
  else if (to.matched.some(record => record.meta.redirectIfAuth)) {
    if (authService.isAuthenticated()) {
      // Si ya está autenticado, redirigir al Dashboard
      next({ path: '/dashboard' });
    } else {
      // Si no está autenticado, continuar normalmente
      next();
    }
  } else {
    // Si no tiene requisitos especiales, continuar normalmente
    next();
  }
});

export default router;
