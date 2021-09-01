import {RouteRecordRaw} from 'vue-router';

const endpoint = '/login';

export default [
  {
    name: 'Login',
    path: endpoint,
    component: () => import('@/views/login/Login.vue'),
  },
] as Array<RouteRecordRaw>;
