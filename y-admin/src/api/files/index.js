import request from '@/utils/request'

export function queryAsk(params) {
  return request({
    url: '/yungo/queryask',
    method: 'get',
    params
  })
}
export function dealAsk(params) {
  return request({
    url: '/yungo/deal',
    method: 'post',
    params
  })
}
export function getFtypes(params) {
  return request({
    url: '/yungo/ftypes',
    method: 'get',
    params
  })
}
export function fileDownload(params) {
  return request({
    url: '/yungo/download',
    method: 'get',
    params
  })
}