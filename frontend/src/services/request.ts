import axios, {AxiosRequestConfig, AxiosResponse} from 'axios';
import {ElMessageBox} from 'element-plus';
import router from '@/router';
import {getRequestBaseUrl} from '@/utils/request';

// TODO: request interception

// response interception
let msgBoxVisible = false;
axios.interceptors.response.use(res => {
  return res;
}, err => {
  const status = err?.response?.status;
  if (status === 401) {
    if (msgBoxVisible) return;
    msgBoxVisible = true;
    ElMessageBox.confirm('You seem to have been logged-out, try to login again?', 'Unauthorized', {type: 'warning'})
      .then(_ => router.push('/login'))
      .finally(() => {
        msgBoxVisible = false;
      });
  } else {
    console.error(err);
  }
});

const useRequest = () => {
  const baseUrl = getRequestBaseUrl();

  const getHeaders = (): any => {
    // headers
    const headers = {} as any;

    // add token to headers
    const token = localStorage.getItem('token');
    if (token) {
      headers['Authorization'] = token;
    }

    return headers;
  }

  const request = async <R = any>(opts: AxiosRequestConfig): Promise<R> => {
    // base url
    const baseURL = baseUrl;

    // headers
    const headers = getHeaders();

    // axios response
    const res = await axios.request({
      ...opts,
      baseURL,
      headers,
    });

    // response data
    return res.data;
  };

  const get = async <T = any, R = ResponseWithData<T>, PM = any>(url: string, params?: PM, opts?: AxiosRequestConfig): Promise<R> => {
    opts = {
      ...opts,
      method: 'GET',
      url,
      params,
    };
    return await request<R>(opts);
  };

  const post = async <T = any, R = ResponseWithData<T>, PM = any>(url: string, data?: T, params?: PM, opts?: AxiosRequestConfig): Promise<R> => {
    opts = {
      ...opts,
      method: 'POST',
      url,
      data,
      params,
    };
    return await request<R>(opts);
  };

  const put = async <T = any, R = ResponseWithData<T>, PM = any>(url: string, data?: T, params?: PM, opts?: AxiosRequestConfig): Promise<R> => {
    opts = {
      ...opts,
      method: 'PUT',
      url,
      data,
      params,
    };
    return await request<R>(opts);
  };

  const del = async <T = any, R = ResponseWithData<T>, PM = any>(url: string, data?: T, params?: PM, opts?: AxiosRequestConfig): Promise<R> => {
    opts = {
      ...opts,
      method: 'DELETE',
      url,
      data,
      params,
    };
    return await request<R>(opts);
  };

  const getList = async <T = any>(url: string, params?: ListRequestParams, opts?: AxiosRequestConfig) => {
    // normalize conditions
    if (params && Array.isArray(params.conditions)) {
      params.conditions = JSON.stringify(params.conditions);
    }

    // get request
    const res = await get<T, ResponseWithListData<T>, ListRequestParams>(url, params, opts);

    // normalize array data
    if (!res.data) {
      res.data = [];
    }

    return res;
  };

  const getAll = async <T = any>(url: string, opts?: AxiosRequestConfig) => {
    return await getList(url, {all: true}, opts);
  };

  const postList = async <T = any, R = Response, PM = any>(url: string, data?: BatchRequestPayloadWithJsonStringData, params?: PM, opts?: AxiosRequestConfig): Promise<R> => {
    return await post<BatchRequestPayloadWithJsonStringData, R, PM>(url, data, params, opts);
  };

  const putList = async <T = any, R = ResponseWithListData, PM = any>(url: string, data?: T[], params?: PM, opts?: AxiosRequestConfig): Promise<R> => {
    return await put<T[], R, PM>(url, data, params, opts);
  };

  const delList = async <T = any, R = Response, PM = any>(url: string, data?: BatchRequestPayload, params?: PM, opts?: AxiosRequestConfig): Promise<R> => {
    return await del<BatchRequestPayload, R, PM>(url, data, params, opts);
  };

  const requestRaw = async <R = any>(opts: AxiosRequestConfig): Promise<AxiosResponse> => {
    // base url
    const baseURL = baseUrl;

    // headers
    const headers = getHeaders();

    // axios response
    return await axios.request({
      ...opts,
      baseURL,
      headers,
    });
  };

  const getRaw = async <T = any, PM = any>(url: string, params?: PM, opts?: AxiosRequestConfig): Promise<AxiosResponse> => {
    opts = {
      ...opts,
      method: 'GET',
      url,
      params,
    };
    return await requestRaw(opts);
  };

  return {
    // public variables and methods
    baseUrl,
    request,
    get,
    post,
    put,
    del,
    getList,
    getAll,
    postList,
    putList,
    delList,
    requestRaw,
    getRaw,
  };
};

export default useRequest;
