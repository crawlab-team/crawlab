import axios from 'axios'
import router from '../router'
import { Message } from 'element-ui'

let baseUrl = process.env.VUE_APP_BASE_URL ? process.env.VUE_APP_BASE_URL : 'http://localhost:8000'

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
    if (response.status === 400) {
      Message.error(response.data.error)
    }
    if (response.status === 401 && router.currentRoute.path !== '/login') {
      console.log('login')
      router.push('/login')
    }
    if (response.status === 500) {
      Message.error(response.data.error)
    }
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
