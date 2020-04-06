import axios from 'axios'
import router from '../router'
import { Message } from 'element-ui'

// 根据 VUE_APP_BASE_URL 生成 baseUrl
let baseUrl = process.env.VUE_APP_BASE_URL ? process.env.VUE_APP_BASE_URL : 'http://localhost:8000'
if (!baseUrl.match(/^https?/i)) {
  baseUrl = `${window.location.protocol}//${window.location.host}${process.env.VUE_APP_BASE_URL}`
}

// 如果 Docker 中设置了 CRAWLAB_API_ADDRESS 这个环境变量，则会将 baseUrl 覆盖
const CRAWLAB_API_ADDRESS = '###CRAWLAB_API_ADDRESS###'
if (!CRAWLAB_API_ADDRESS.match('CRAWLAB_API_ADDRESS')) {
  baseUrl = CRAWLAB_API_ADDRESS
}

const request = (method, path, params, data, others = {}) => {
  const url = baseUrl + path
  const headers = {
    'Authorization': window.localStorage.getItem('token')
  }
  return axios({
    method,
    url,
    params,
    data,
    headers,
    ...others
  }).then((response) => {
    if (response.status === 200) {
      return Promise.resolve(response)
    }
    return Promise.reject(response)
  }).catch((e) => {
    let response = e.response
    if (!response) {
      return e
    }
    if (response.status === 400) {
      Message.error(response.data.error)
    }
    if (response.status === 401 && router.currentRoute.path !== '/login') {
      router.push('/login')
    }
    if (response.status === 500) {
      Message.error(response.data.error)
    }
    return e
  })
}

const get = (path, params) => {
  return request('GET', path, params)
}

const post = (path, data) => {
  return request('POST', path, {}, data)
}

const put = (path, data) => {
  return request('PUT', path, {}, data)
}

const del = (path, data) => {
  return request('DELETE', path, {}, data)
}

export default {
  baseUrl,
  request,
  get,
  post,
  put,
  delete: del
}
