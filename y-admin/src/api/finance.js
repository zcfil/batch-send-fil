import request from '@/utils/request'



// export const register = data => request({
//   url: '/register',
//   method: 'post',
//   data
// })

export const ApplyList = (data) => request({
  url: '/yungo/getApplyList',
  params: data,
  method: 'get'
})

export const ManualList = (data) => request({
  url: '/yungo/getManualList',
  params: data,
  method: 'get'
})
export const updateManual = (data) => request({
  url: '/yungo/updateManual',
  params: data,
  method: 'post'
})
export const sendManual = (data) => request({
  url: '/yungo/sendManual',
  params: data,
  method: 'post'
})
export const manualAdd = (data) => request({
  url: '/yungo/manualAdd',
  params: data,
  method: 'post'
})

export const Sends = (data) => request({
  url: '/yungo/sends',
  params: data,
  method: 'post'
})

export const Refuse = (data) => request({
  url: '/yungo/refuse',
  params: data,
  method: 'post'
})

export const walletbalance = (data) => request({
  url: '/yungo/walletBalance',
  params: data,
  method: 'get'
})

export const batchlist = (data) => request({
  url: '/yungo/getBatchList',
  params: data,
  method: 'get'
})

export const batchSends = (data) => request({
  url: '/yungo/batchSends',
  params: data,
  method: 'post'
})

export const batchRefuse = (data) => request({
  url: '/yungo/batchRefuse',
  params: data,
  method: 'post'
})

export const GetConfig = (data) => request({
  url: '/yungo/getConfig',
  params: data,
  method: 'get'
})
export const SetConfig = (data) => request({
  url: '/yungo/setConfig',
  params: data,
  method: 'post'
})
