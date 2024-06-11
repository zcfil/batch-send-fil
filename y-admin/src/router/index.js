import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/register',
    component: () => import('@/views/login/register'),
    hidden: true
  },
  {
    path: '/download',
    component: () => import('@/views/login/googleAuthenticator/app'),
    hidden: true
  },
  {
    path: '/verify',
    component: () => import('@/views/login/googleAuthenticator/verify'),
    hidden: true
  },
  {
    path: '/loginVerify',
    component: () => import('@/views/login/googleAuthenticator/loginVerify'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },

  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    meta: { title: '首页', icon: 'dashboard' },
    children: [{
      path: 'Dashboard',
      name: '首页',
      component: () => import('@/views/dashboard/index'),
      meta: { title: '首页', icon: 'dashboard' }
    }]
  },
  {
    path: '/finance',
    component: Layout,
    redirect: '/finance/index/index',
    name: 'All',
    meta: { title: '财务管理', icon: 'el-icon-money' },
    children: [
      {
        path: 'index',
        name: 'batch',
        component: () => import('@/layout/components/EmptyView'),
        meta: { title: '提币审核', icon: 'el-icon-bank-card' },
        redirect: '/finance/index/index',
        children: [
          {
            path: 'index',
            name: 'batch',
            component: () => import('@/views/finance/index'),
            meta: { title: '', icon: 'el-icon-bank-card'},
          },
          {
            path: 'audita',
            name: 'MentionAudit',
            hidden: true,
            component: () => import('@/views/finance/audit'),
            meta: { title: '批次详情', icon: 'el-icon-bank-card'}    
          },
        ]
      },
            {
        // path: 'record',
        // name: 'MentionRecord',
        // component: () => import('@/views/finance/record'),
        // meta: { title: '提币记录', icon: 'el-icon-notebook-1' }
        path: 'record',
        name: 'Record',
        component: () => import('@/layout/components/EmptyView'),
        meta: { title: '提币记录', icon: 'el-icon-notebook-1' },
        redirect: '/finance/index/index',
        children: [
          {
            path: 'index',
            name: 'Record',
            component: () => import('@/views/finance/batchlist'),
            meta: { title: '', icon: 'el-icon-notebook-1'},
          },
          {
            path: 'record-1',
            name: 'MentionRecord',
            hidden: true,
            component: () => import('@/views/finance/record'),
            meta: { title: '批次详情', icon: 'el-icon-bank-card' }    
          },
        ]
      },
      {
        path: 'manual',
        name: 'Manual',
        component: () => import('@/layout/components/EmptyView'),
        meta: { title: '手动提币', icon: 'el-icon-tickets' },
        redirect: '/finance/index/index',
        children: [
          {
            path: 'index',
            name: 'Manual',
            component: () => import('@/views/finance/manual'),
            meta: { title: '', icon: 'el-icon-tickets'},
          },
        ]
      },
      {
        path: 'manualRecord',
        name: 'ManualRecord',
        component: () => import('@/layout/components/EmptyView'),
        meta: { title: '手动提币记录', icon: 'el-icon-document' },
        redirect: '/finance/index/index',
        children: [
          {
            path: 'index',
            name: 'ManualRecord',
            component: () => import('@/views/finance/manualRecord'),
            meta: { title: '', icon: 'el-icon-document'},
          },
        ]
      }
    ]
  },
  {
    path: '/wallet',
    component: Layout,
    redirect: '/wallet/index/index',
    name: 'All',
    meta: { title: '钱包管理', icon: 'el-icon-s-cooperation' },
    children: [
      {
        path: 'index',
        name: 'wallet',
        component: () => import('@/layout/components/EmptyView'),
        meta: { title: '创建钱包', icon: 'el-icon-s-cooperation' },
        redirect: '/wallet/index/index',
        children: [
          {
            path: 'index',
            name: 'create',
            component: () => import('@/views/wallet/index'),
            meta: { title: '', icon: 'el-icon-s-platform'},
          },
          {
            path: 'wallet-1',
            name: 'walletCreate',
            hidden: true,
            component: () => import('@/views/wallet/addwallet'),
            meta: { title: '确认助记词', icon: 'el-icon-s-claim' }    
          },
        ]
      },
      {
        path: 'backups',
        name: 'backups',
        component: () => import('@/layout/components/EmptyView'),
        meta: { title: '备份钱包', icon: 'el-icon-s-claim' },
        redirect: '/finance/index/index',
        children: [
          {
            path: 'index',
            name: 'backups',
            component: () => import('@/views/wallet/backups'),
            meta: { title: '', icon: 'el-icon-s-claim'},
          },
        ]
      }
    ]
  },
  {
    path: '/profile',
    component: Layout,
    meta: { title: 'Profile', icon: 'el-icon-user-solid' },
    children: [
      {
        path: 'index',
        name: 'Profile',
        component: () => import('@/views/profile/index'),
        meta: { title: 'Profile', icon: 'el-icon-user-solid' }
      }
    ]
  },
  

  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
