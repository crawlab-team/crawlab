import axios, {AxiosRequestConfig} from 'axios';

// TODO: request interception

// TODO: response interception

const useRequest = () => {
  // implementation
  const baseUrl = process.env.VUE_APP_API_BASE_URL || 'http://localhost:8000';

  const request = async <R = any>(opts: AxiosRequestConfig): Promise<R> => {
    // base url
    const baseURL = baseUrl;

    // headers
    const headers = {} as any;

    // add token to headers
    const token = localStorage.getItem('token');
    if (token) {
      headers['Authorization'] = token;
    }

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
    return await getList(url, {}, opts);
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
  };
};

export default useRequest;
