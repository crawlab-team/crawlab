import {RouteRecordRaw} from 'vue-router';
import {
  TAB_NAME_DATA,
  TAB_NAME_FILES,
  TAB_NAME_OVERVIEW,
  TAB_NAME_SCHEDULES,
  TAB_NAME_SETTINGS,
  TAB_NAME_TASKS
} from '@/constants/tab';

export default [
  {
    name: 'SpiderList',
    path: 'spiders',
    component: () => import('@/views/spider/list/SpiderList.vue'),
  },
  {
    name: 'SpiderDetail',
    path: 'spiders/:id',
    redirect: to => {
      return {path: to.path + '/' + TAB_NAME_OVERVIEW};
    },
    component: () => import('@/views/spider/detail/SpiderDetail.vue'),
    children: [
      {
        path: TAB_NAME_OVERVIEW,
        component: () => import('@/views/spider/detail/tabs/SpiderDetailTabOverview.vue'),
      },
      {
        path: TAB_NAME_FILES,
        component: () => import('@/views/spider/detail/tabs/SpiderDetailTabFiles.vue'),
      },
      {
        path: TAB_NAME_TASKS,
        component: () => import('@/views/spider/detail/tabs/SpiderDetailTabTasks.vue'),
      },
      {
        path: TAB_NAME_SCHEDULES,
        component: () => import('@/views/spider/detail/tabs/SpiderDetailTabSchedules.vue'),
      },
      {
        path: TAB_NAME_DATA,
        component: () => import('@/views/spider/detail/tabs/SpiderDetailTabData.vue'),
      },
      {
        path: TAB_NAME_SETTINGS,
        component: () => import('@/views/spider/detail/tabs/SpiderDetailTabSettings.vue'),
      },
    ]
  },
] as Array<RouteRecordRaw>;
