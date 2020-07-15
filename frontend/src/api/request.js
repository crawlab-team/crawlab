import service from '@/utils/request'

const get = (path, params) => {
  return service.get(path, {
    params
  })
}

const post = (path, data) => {
  return service.post(path, data)
}

const put = (path, data) => {
  return service.put(path, data)
}

const del = (path, data) => {
  return service.delete(path, {
    data
  })
}
const request = service.request

export default {
  baseUrl: service.defaults.baseURL,
  request,
  get,
  post,
  put,
  delete: del
}
