import {RouteRecordRaw} from 'vue-router';
import {TAB_NAME_OVERVIEW} from '@/constants/tab';

export default [
  {
    path: 'users',
    component: () => import('@/views/user/list/UserList.vue'),
  },
  {
    path: 'users/:id',
    redirect: to => {
      return {path: to.path + '/' + TAB_NAME_OVERVIEW};
    },
    component: () => import('@/views/user/detail/UserDetail.vue'),
    children: [
      {
        path: TAB_NAME_OVERVIEW,
        component: () => import('@/views/user/detail/tabs/UserDetailTabOverview.vue'),
      },
    ]
  },
] as Array<RouteRecordRaw>;
