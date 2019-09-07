import Vue from 'vue'
import Router from 'vue-router'
import store from '../store'
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
  { path: '/login', component: () => import('../views/login/index'), hidden: true },
  { path: '/change_password', component: () => import('../views/change_password/index'), hidden: true },
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
    path: '/sites',
    component: Layout,
    hidden: true,
    meta: {
      title: 'Site',
      icon: 'fa fa-sitemap'
    },
    children: [
      {
        path: '',
        name: 'SiteList',
        component: () => import('../views/site/SiteList'),
        meta: {
          title: 'Sites',
          icon: 'fa fa-sitemap'
        }
      }
    ]
  },
  {
    path: '/users',
    component: Layout,
    meta: {
      title: 'User',
      icon: 'fa fa-user'
    },
    children: [
      {
        path: '',
        name: 'UserList',
        component: () => import('../views/user/UserList'),
        meta: {
          title: 'Users',
          icon: 'fa fa-user'
        }
      }
    ]
  }
]

const router = new Router({
  // mode: 'history', //后端支持可开
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRouterMap
})
router.onReady(async () => {
  const settings = await store.dispatch('settings/getSettings')
  let routes = []
  if (settings.can_register) {
    routes = [
      ...routes,
      ...[
        { path: '/signup', component: () => import('../views/login/index'), hidden: true }
      ]
    ]
  }

  routes.push({ path: '*', redirect: '/404', hidden: true })
  router.addRoutes(routes)
  console.log('ready')
})

router.beforeEach(async (to, from, next) => {
  if (to.meta && to.meta.title) {
    window.document.title = `Crawlab - ${to.meta.title}`
  } else {
    window.document.title = 'Crawlab'
  }
  // try {
  // switch (to.path) {
  //   case '/signup':
  //     const settings = await store.dispatch('settings/getSettings')
  //     console.log(settings)
  //     if (settings.can_register) {
  //       router.addRoutes([
  //         { path: '/signup', component: () => import('../views/login/index'), hidden: true }
  //
  //       ])
  //       router.options.routes = store.getters.routers
  //       // next({ ...to, replace: true })// hack方法 确保addRoutes已完成 ,set the replace: true so the navigation will not leave a history record
  //     }
  //     break
  // }
  // } catch (e) {
  //
  // }

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

router.afterEach((to, from, next) => {
  if (to.path) {
    stats.sendPv(to.path)
  }
})

export default router
