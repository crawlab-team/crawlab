import {RouteRecordRaw} from 'vue-router';

const endpoint = '';

export default [
  {
    name: 'Home',
    path: endpoint,
    component: () => import('@/views/home/Home.vue'),
  },
] as Array<RouteRecordRaw>;
