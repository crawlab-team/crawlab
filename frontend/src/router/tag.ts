import {RouteRecordRaw} from 'vue-router';
import {TAB_NAME_OVERVIEW} from '@/constants/tab';

const endpoint = 'tags';

export default [
  {
    path: endpoint,
    component: () => import('@/views/tag/list/TagList.vue'),
  },
  {
    path: `${endpoint}/:id`,
    redirect: to => {
      return {path: to.path + '/overview'};
    },
    component: () => import('@/views/tag/detail/TagDetail.vue'),
    children: [
      {
        path: TAB_NAME_OVERVIEW,
        component: () => import('@/views/tag/detail/tabs/TagDetailTabOverview.vue'),
      },
    ]
  }
] as Array<RouteRecordRaw>;
