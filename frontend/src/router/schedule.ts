import {RouteRecordRaw} from 'vue-router';
import {TAB_NAME_OVERVIEW, TAB_NAME_TASKS} from '@/constants/tab';

export default [
  {
    name: 'ScheduleList',
    path: 'schedules',
    component: () => import('@/views/schedule/list/ScheduleList.vue'),
  },
  {
    name: 'ScheduleDetail',
    path: 'schedules/:id',
    redirect: to => {
      return {path: to.path + '/' + TAB_NAME_OVERVIEW};
    },
    component: () => import('@/views/schedule/detail/ScheduleDetail.vue'),
    children: [
      {
        path: TAB_NAME_OVERVIEW,
        component: () => import('@/views/schedule/detail/tabs/ScheduleDetailTabOverview.vue'),
      },
      {
        path: TAB_NAME_TASKS,
        component: () => import('@/views/schedule/detail/tabs/ScheduleDetailTabTasks.vue'),
      },
    ]
  },
] as Array<RouteRecordRaw>;
