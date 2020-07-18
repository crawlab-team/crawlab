import Vue from 'vue'
import Router from 'vue-router'

import store from '../store'
import request from '../api/request'
import stats from '../utils/stats'

/* Layout */
import Layout from '../views/layout/Layout'

// in development-env not use lazy-loading, because lazy-loading too many pages will cause webpack hot update too slow. so only in production use lazy-loading;
// detail: https://panjiachen.github.io/vue-element-admin-site/#/lazy-loading

Vue.use(Router)

/**
 * hidden: true                   if `hidden:true` will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu, whatever its child routes length
 *                                if not set alwaysShow, only more than one route under the children
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noredirect           if `redirect:noredirect` will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    title: 'title'               the name show in submenu and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if false, the item will hidden in breadcrumb(default is true)
  }
 **/
export const constantRouterMap = [
  {
    path: '/login',
    component: () => import('../views/login/index'),
    hidden: true
  },
  {
    path: '/signup',
    component: () => import('../views/login/index'),
    hidden: true
  },
  { path: '/404', component: () => import('../views/404'), hidden: true },
  { path: '/', redirect: '/home' },

  // Crawlab Pages
  {
    path: '/home',
    component: Layout,
    children: [
      {
        path: '',
        component: () => import('../views/home/Home'),
        meta: {
          title: 'Home',
          icon: 'fa fa-home'
        }
      }
    ]
  },
  {
    path: '/projects',
    component: Layout,
    meta: {
      title: 'Project',
      icon: 'fa fa-gear'
    },
    children: [
      {
        path: '',
        name: 'Project',
        component: () => import('../views/project/ProjectList'),
        meta: {
          title: 'Project',
          icon: 'fa fa-code-fork'
        }
      }
    ]
  },
  {
    path: '/spiders',
    component: Layout,
    meta: {
      title: 'Spider',
      icon: 'fa fa-bug'
    },
    children: [
      {
        path: '',
        name: 'SpiderList',
        component: () => import('../views/spider/SpiderList'),
        meta: {
          title: 'Spiders',
          icon: 'fa fa-bug'
        }
      },
      {
        path: ':id',
        name: 'SpiderDetail',
        component: () => import('../views/spider/SpiderDetail'),
        meta: {
          title: 'Spider Detail',
          icon: 'fa fa-circle-o'
        },
        hidden: true
      }
    ]
  },
  {
    path: '/tasks',
    component: Layout,
    meta: {
      title: 'Task',
      icon: 'fa fa-list'
    },
    children: [
      {
        path: '',
        name: 'TaskList',
        component: () => import('../views/task/TaskList'),
        meta: {
          title: 'Tasks',
          icon: 'fa fa-list'
        }
      },
      {
        path: ':id',
        name: 'TaskDetail',
        component: () => import('../views/task/TaskDetail'),
        meta: {
          title: 'Task Detail',
          icon: 'fa fa-circle-o'
        },
        hidden: true
      }
    ]
  },
  {
    path: '/schedules',
    component: Layout,
    meta: {
      title: 'Schedules',
      icon: 'fa fa-calendar'
    },
    hidden: false,
    children: [
      {
        path: '',
        name: 'ScheduleList',
        component: () => import('../views/schedule/ScheduleList'),
        meta: {
          title: 'Schedules',
          icon: 'fa fa-calendar'
        }
      }
    ]
  },
  {
    path: '/nodes',
    component: Layout,
    meta: {
      title: 'Node',
      icon: 'fa fa-server'
    },
    children: [
      {
        path: '',
        name: 'NodeList',
        component: () => import('../views/node/NodeList'),
        meta: {
          title: 'Nodes',
          icon: 'fa fa-server'
        }
      },
      {
        path: ':id',
        name: 'NodeDetail',
        component: () => import('../views/node/NodeDetail'),
        meta: {
          title: 'Node Detail',
          icon: 'fa fa-circle-o'
        },
        hidden: true
      }
    ]
  },
  {
    path: '/repos',
    component: Layout,
    meta: {
      title: 'Spider Market',
      icon: 'fa fa-cloud'
    },
    children: [
      {
        path: '',
        name: 'RepoList',
        component: () => import('../views/repo/RepoList'),
        meta: {
          title: 'Spider Market',
          icon: 'fa fa-cloud'
        }
      }
    ]
  },
  {
    path: '/disclaimer',
    component: Layout,
    meta: {
      title: 'Disclaimer',
      icon: 'fa fa-exclamation-triangle'
    },
    children: [
      {
        path: '',
        name: 'Disclaimer',
        component: () => import('../views/doc/Disclaimer'),
        meta: {
          title: 'Disclaimer',
          icon: 'fa fa-exclamation-triangle'
        }
      }
    ]
  },
  // {
  //   path: '/challenges',
  //   component: Layout,
  //   meta: {
  //     title: 'ChallengeList',
  //     icon: 'fa fa-flash'
  //   },
  //   children: [
  //     {
  //       path: '',
  //       name: 'ChallengeList',
  //       component: () => import('../views/challenge/ChallengeList'),
  //       meta: {
  //         title: 'Challenges',
  //         icon: 'fa fa-flash'
  //       }
  //     }
  //   ]
  // },
  {
    path: '/feedback',
    component: Layout,
    meta: {
      title: 'Feedback',
      icon: 'fa fa-commenting-o'
    },
    children: [
      {
        path: '',
        name: 'Feedback',
        component: () => import('../views/feedback/Feedback'),
        meta: {
          title: 'Feedback',
          icon: 'fa fa-commenting'
        }
      }
    ]
  },
  {
    path: '/users',
    component: Layout,
    meta: {
      title: 'User',
      icon: 'fa fa-users'
    },
    children: [
      {
        path: '',
        name: 'UserList',
        component: () => import('../views/user/UserList'),
        meta: {
          title: 'Users',
          icon: 'fa fa-users'
        }
      }
    ]
  },
  {
    path: '/setting',
    component: Layout,
    meta: {
      title: 'Setting',
      icon: 'fa fa-gear'
    },
    children: [
      {
        path: '',
        name: 'Setting',
        component: () => import('../views/setting/Setting'),
        meta: {
          title: 'Setting',
          icon: 'fa fa-gear'
        }
      }
    ]
  },

  { path: '*', redirect: '/404', hidden: true }
]

const router = new Router({
  // mode: 'history', //后端支持可开
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRouterMap
})

router.beforeEach((to, from, next) => {
  if (to.meta && to.meta.title) {
    window.document.title = `Crawlab - ${to.meta.title}`
  } else {
    window.document.title = 'Crawlab'
  }

  if (['/login', '/signup'].includes(to.path)) {
    next()
  } else {
    if (window.localStorage.getItem('token')) {
      next()
    } else {
      next('/login')
    }
  }
})

router.afterEach(async(to, from, next) => {
  if (to.path) {
    await store.dispatch('setting/getSetting')
    const res = await request.get('/version')
    const version = res.data.data
    store.commit('version/SET_VERSION', version)
    sessionStorage.setItem('v', version)
    stats.sendPv(to.path)
  }
})

export default router
