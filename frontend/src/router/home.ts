import {RouteRecordRaw} from 'vue-router';

const endpoint = '';

export default [
  {
    path: endpoint,
    component: () => import('@/views/home/Home.vue'),
  },
] as Array<RouteRecordRaw>;
