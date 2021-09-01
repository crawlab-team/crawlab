import {RouteRecordRaw} from 'vue-router';
import {TAB_NAME_OVERVIEW} from '@/constants/tab';

const endpoint = 'plugins';

export default [
  {
    name: 'PluginList',
    path: endpoint,
    component: () => import('@/views/plugin/list/PluginList.vue'),
  },
  {
    name: 'PluginDetail',
    path: `${endpoint}/:id`,
    redirect: to => {
      return {path: to.path + '/' + TAB_NAME_OVERVIEW};
    },
    component: () => import('@/views/plugin/detail/PluginDetail.vue'),
    children: [
      {
        path: TAB_NAME_OVERVIEW,
        component: () => import('@/views/plugin/detail/tabs/PluginDetailTabOverview.vue'),
      },
    ]
  },
] as Array<RouteRecordRaw>;
