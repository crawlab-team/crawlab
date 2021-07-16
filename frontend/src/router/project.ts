import {RouteRecordRaw} from 'vue-router';
import {TAB_NAME_OVERVIEW, TAB_NAME_SPIDERS} from '@/constants/tab';

const endpoint = 'projects';

export default [
  {
    path: endpoint,
    component: () => import('@/views/project/list/ProjectList.vue'),
  },
  {
    path: `${endpoint}/:id`,
    redirect: to => {
      return {path: to.path + '/overview'};
    },
    component: () => import('@/views/project/detail/ProjectDetail.vue'),
    children: [
      {
        path: TAB_NAME_OVERVIEW,
        component: () => import('@/views/project/detail/tabs/ProjectDetailTabOverview.vue'),
      },
      {
        path: TAB_NAME_SPIDERS,
        component: () => import('@/views/project/detail/tabs/ProjectDetailTabSpiders.vue'),
      },
    ]
  },
] as Array<RouteRecordRaw>;
