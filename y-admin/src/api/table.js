import request from '@/utils/request'

// export function getList(params) {
//   return request({
//     url: '/vue-admin-template/table/list',
//     method: 'get',
//     params
//   })
// }

export function getList(params) {
  return request({
    url: '/yungo/list/imports',
    method: 'get',
    params
  })
}

export const netminers = params => request({
  url: '/yungo/netminers',
  method: 'get',
  params
})

export const netpeers = params => request({
  url: '/yungo/netpeers',
  method: 'get',
  params
})

export const minerinfo = params => request({
  url: '/yungo/minerinfo',
  method: 'get',
  params
})

export const queryAsk = params => request({
  url: '/yungo/queryask',
  method: 'get',
  params
})

export const getpleagefil = params => request({
  url: '/pleagefil',
  method: 'get',
  params
})

export const getstatisticalfil = params => request({
  url: '/statisticalfil',
  method: 'get',
  params
})