import axios from 'axios'
import router from '../router'

let baseUrl = process.env.VUE_APP_BASE_URL ? process.env.VUE_APP_BASE_URL : 'http://localhost:8000'

const request = async (method, path, params, data, others = {}) => {
  try {
    const url = baseUrl + path
    const headers = {
      'Authorization': window.localStorage.getItem('token')
    }
    const response = await axios({
      method,
      url,
      params,
      data,
      headers,
      ...others
    })
    // console.log(response)
    return response
  } catch (e) {
    if (e.response.status === 401) {
      if (e.response.data.code === 11006) {
        if (router.currentRoute.path !== '/change_password') {
          await router.push('/change_password')
        }
      } else if (router.currentRoute.path !== '/login') {
        await router.push('/login')
      }
    }
    await Promise.reject(e)
  }

  // return new Promise((resolve, reject) => {
  //   const url = baseUrl + path
  //   const headers = {
  //     'Authorization': window.localStorage.getItem('token')
  //   }
  //   axios({
  //     method,
  //     url,
  //     params,
  //     data,
  //     headers,
  //     ...others
  //   })
  //     .then(resolve)
  //     .catch(error => {
  //       console.log(error)
  //       if (error.response.status === 401) {
  //         router.push('/login')
  //       }
  //       reject(error)
  //     })
  // })
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
  return request('DELETE', path)
}

export default {
  baseUrl,
  request,
  get,
  post,
  put,
  delete: del
}
