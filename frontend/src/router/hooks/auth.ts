import {Router} from 'vue-router';
import {getToken} from '@/utils/auth';

export const ANOMALOUS_ROUTES = [
  '/login',
];

export const initRouterAuth = (router: Router) => {
  router.beforeEach((to, from, next) => {
    if (ANOMALOUS_ROUTES.includes(to.path)) {
      return next();
    }

    const token = getToken();
    if (!token) {
      return next('/login');
    }

    return next();
  });
};
