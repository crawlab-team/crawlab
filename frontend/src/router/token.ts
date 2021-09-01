import {RouteRecordRaw} from 'vue-router';

const endpoint = 'tokens';

export default [
  {
    name: 'TokenList',
    path: endpoint,
    component: () => import('@/views/token/list/TokenList.vue'),
  },
] as Array<RouteRecordRaw>;
