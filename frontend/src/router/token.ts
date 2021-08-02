import {RouteRecordRaw} from 'vue-router';

const endpoint = 'tokens';

export default [
  {
    path: endpoint,
    component: () => import('@/views/token/list/TokenList.vue'),
  },
] as Array<RouteRecordRaw>;
