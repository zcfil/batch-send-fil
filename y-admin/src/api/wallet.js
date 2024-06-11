import request from '@/utils/request'


export const exportMnemonic = (data) => request({
  url: '/yungo/exportMnemonic',
  params: data,
  method: 'get'
})

export const exportPrivateKey = (data) => request({
    url: '/yungo/exportPrivateKey',
    params: data,
    method: 'get'
  })
  
export const getMnemonic = (data) => request({
  url: '/yungo/getMnemonic',
  params: data,
  method: 'get'
})

export const newWallet = (data) => request({
  url: '/yungo/newWallet',
  params: data,
  method: 'post'
})
export const importWallet = (data) => request({
  url: '/yungo/importWallet',
  params: data,
  method: 'post'
})
export const importMnemonic = (data) => request({
  url: '/yungo/importMnemonic',
  params: data,
  method: 'post'
})
export const walletList = (data) => request({
  url: '/yungo/walletList',
  params: data,
  method: 'get'
})

export const delWallet = (data) => request({
  url: '/yungo/delWallet',
  params: data,
  method: 'post'
})

export const setWallet = (data) => request({
  url: '/yungo/setWallet',
  params: data,
  method: 'post'
})