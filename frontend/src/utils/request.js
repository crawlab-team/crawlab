import axios from 'axios'
import { MessageBox, Message } from 'element-ui'
import store from '@/store'
import { getToken } from '@/utils/auth'
import i18n from '@/i18n'
import router from '@/router'

const codeMessage = {
  200: '服务器成功返回请求的数据。',
  201: '新建或修改数据成功。',
  202: '一个请求已经进入后台排队（异步任务）。',
  204: '删除数据成功。',
  400: '发出的请求有错误，服务器没有进行新建或修改数据的操作。',
  401: '用户没有权限（令牌、用户名、密码错误）。',
  403: '用户得到授权，但是访问是被禁止的。',
  404: '发出的请求针对的是不存在的记录，服务器没有进行操作。',
  406: '请求的格式不可得。',
  410: '请求的资源被永久删除，且不会再得到的。',
  422: '当创建一个对象时，发生一个验证错误。',
  500: '服务器发生错误，请检查服务器。',
  502: '网关错误。',
  503: '服务不可用，服务器暂时过载或维护。',
  504: '网关超时。'
}

/**
 * 异常处理程序
 */
const errorHandler = (error) => {
  const { response } = error
  const routePath = router.currentRoute.path
  if (response && response.status) {
    const errorText = codeMessage[response.status] || response.statusText
    const { status } = response
    Message({
      message: `请求错误 ${status}: ${response.request.responseURL},${errorText}`,
      type: 'error',
      duration: 5 * 1000
    })
    switch (status) {
      case 401:
        if (routePath !== '/login' && routePath !== '/') {
          MessageBox.confirm(
            i18n.t('auth.login_expired_message'),
            i18n.t('auth.login_expired_title'), {
              confirmButtonText: i18n.t('auth.login_expired_confirm'),
              cancelButtonText: i18n.t('auth.login_expired_cancel'),
              type: 'warning'
            }).then(() => {
            store.dispatch('user/resetToken').then(() => {
              location.reload()
            })
          })
        }
        break
      default:
    }
  } else if (!response) {
    Message({
      message: `您的网络发生异常，无法连接服务器`,
      type: 'error',
      duration: 5 * 1000
    })
  }
  return response
}

// 根据 VUE_APP_BASE_URL 生成 baseUrl
let baseUrl = process.env.VUE_APP_BASE_URL
  ? process.env.VUE_APP_BASE_URL
  : 'http://localhost:8000'
if (!baseUrl.match(/^https?/i)) {
  baseUrl = `${window.location.protocol}//${window.location.host}${process.env.VUE_APP_BASE_URL}`
}

// 如果 Docker 中设置了 CRAWLAB_API_ADDRESS 这个环境变量，则会将 baseUrl 覆盖
const CRAWLAB_API_ADDRESS = '###CRAWLAB_API_ADDRESS###'
if (!CRAWLAB_API_ADDRESS.match('CRAWLAB_API_ADDRESS')) {
  baseUrl = CRAWLAB_API_ADDRESS
}
// create an axios instance
const service = axios.create({
  baseURL: baseUrl, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 5000 // request timeout
})
// request interceptor
service.interceptors.request.use(
  config => {
    // do something before request is sent

    if (store.getters.token) {
      // let each request carry token
      // ['X-Token'] is a custom headers key
      // please modify it according to the actual situation
      config.headers['Authorization'] = getToken()
    }
    return config
  },
  error => {
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
   */
  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    return response
  },
  errorHandler
)
export default service
