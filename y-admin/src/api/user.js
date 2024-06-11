import request from '@/utils/request'

// export function login(data) {
//   return request({
//     url: '/vue-admin-template/user/login',
//     method: 'post',
//     data
//   })
// }
export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

// export function getInfo(token) {
//   return request({
//     url: '/vue-admin-template/user/info',
//     method: 'get',
//     params: { token }
//   })
// }  
export function getInfo(token) {
  return request({
    url: '/vue-admin-template/user/info',
    method: 'get',
    params: { token }
  })
}

export function logout() {
  return request({
    url: '/vue-admin-template/user/logout',
    method: 'post'
  })
}

export const register = data => request({
  url: '/admin/register',
  method: 'post',
  data
})

export const memberinfo = () => request({
  url: '/yungo/memberinfo',
  method: 'get'
})

export const resetpwd = data => request({
  url: '/yungo/resetpwd',
  method: 'post',
  data
})
export const getVerificationCode = () => request({
  url: '/yungo/getVerificationCode',
  method: 'get',
})

export const bindVerificationCode = (data) => request({
  url: '/yungo/bindVerificationCode',
  method: 'post',
  params:data
})

export const verifyCode = data => request({
  url: '/yungo/verifyCode',
  method: 'post',
  params: data
})