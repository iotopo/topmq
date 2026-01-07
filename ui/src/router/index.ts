import type { RouteRecordRaw } from 'vue-router'
import { ElMessage } from 'element-plus'
import nProgress from 'nprogress'
import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/pages/layout.vue'
import { getActiveStore } from '@/store'

export const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/login.vue'),
    meta: {
      title: '账号密码登录',
      hidden: true,
    },
  },
  {
    path: '/:pathMatch(.*)*',
    name: '404',
    component: () => import('@/pages/404.vue'),
    meta: {
      title: '页面不存在',
      hidden: true,
    },
  },
  {
    path: '/user-center',
    name: 'UserCenter',
    component: () => import('@/pages/user-center/index.vue'),
    redirect: '/user-center/basic',
    meta: {
      auth: true,
      hidden: true,
    },
    children: [
      {
        path: 'basic',
        name: 'AccountBasic',
        component: () => import('@/pages/user-center/basic.vue'),
        meta: {
          // title: '账户资料',
          auth: true,
          activeMenu: 'basic',
          hidden: true,
        },
      },
      {
        path: 'security',
        name: 'AccountSecurity',
        component: () => import('@/pages/user-center/security.vue'),
        meta: {
          // title: '修改密码',
          auth: true,
          activeMenu: 'security',
          hidden: true,
        },
      },
    ],
  },
  {
    path: '/',
    redirect: '/access/account',
    meta: {
      auth: true,
      hidden: true,
    },
  },
  {
    path: '/access',
    name: 'access',
    component: Layout,
    redirect: '/access/acl',
    meta: {
      title: 'menu.access_control', // 访问控制
      icon: 'ep:lock',
      auth: true,
    },
    children: [
      {
        path: 'account',
        name: 'Account',
        component: () => import('@/pages/access/account.vue'),
        meta: {
          title: 'menu.account', // 客户端认证
          auth: true,
          cache: true,
        },
      },
      {
        path: 'acl',
        name: 'Acl',
        component: () => import('@/pages/access/acl.vue'),
        meta: {
          title: 'menu.acl', // 客户端授权
          auth: true,
          cache: true,
        },
      },
      {
        path: 'blacklist',
        name: 'Blacklist',
        component: () => import('@/pages/access/blacklist.vue'),
        meta: {
          title: 'menu.blacklist', // 黑名单
          auth: true,
          cache: true,
        },
      },
    ],
  },
  {
    path: '/monitoring',
    name: 'Monitoring',
    component: Layout,
    redirect: '/monitoring/clients',
    meta: {
      title: 'menu.monitoring',
      icon: 'ep:monitor',
      auth: true,
    },
    children: [
      {
        path: 'metrics',
        name: 'Metrics',
        component: () => import('@/pages/monitoring/metrics.vue'),
        meta: {
          title: 'menu.metrics',
          cache: true,
          auth: true,
          // visibleFunc(): boolean {
          //   const store = getActiveStore()
          //   return !!store.metrics
          // },
        },
      },
      {
        path: 'clients',
        name: 'Clients',
        component: () => import('@/pages/monitoring/clients.vue'),
        meta: {
          title: 'menu.clients',
          cache: true,
          auth: true,
        },
      },
      {
        path: 'clients/:clientID',
        name: 'ClientDetail',
        component: () => import('@/pages/monitoring/client-detail.vue'),
        props: true,
        meta: {
          title: 'menu.clients',
          cache: false,
          hidden: true,
          auth: true,
        },
      },
      {
        path: 'subscriptions',
        name: 'Subscriptions',
        component: () => import('@/pages/monitoring/subscriptions.vue'),
        meta: {
          title: 'menu.subscriptions',
          cache: true,
          auth: true,
        },
      },
      {
        path: 'retained-messages',
        name: 'RetainedMessages',
        component: () => import('@/pages/monitoring/retained-messages.vue'),
        meta: {
          title: 'menu.retained_messages',
          cache: true,
          auth: true,
        },
      },
    ],
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, from, next) => {
  nProgress.start()
  const store = getActiveStore()
  //   console.log('beforeEach', to)
  if (to.meta.auth && !store.isLogin) {
    try {
      await store.authenticate()
      next({ ...to, replace: true })
      return
    } catch (error: any) {
      const response = error.response
      if (response?.status === 401) {
        if (!to.path.startsWith('/login')) {
          // 重定向到登录页
          next({ path: '/login', query: { redirect: to.fullPath } })
          return
        }
      } else if (response?.status === 500) {
        const data = await response.clone().json()
        if (data.error) {
          ElMessage.error(data.error)
        } else if (data.msg) {
          ElMessage.error(data.msg)
        } else {
          console.warn(data)
        }
      }
      if (to.path !== '/login') {
        next({ path: '/login', query: { redirect: to.fullPath } })
        return
      }
    }
  }

  next()
})
router.afterEach(() => {
  nProgress.done()
})
