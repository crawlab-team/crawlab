import useRequest from '@/services/request';

const {
  get,
  put,
  post,
  del,
  getList,
  getAll,
  putList,
  postList,
  delList,
} = useRequest();

export const useService = <T = any>(endpoint: string): Services<T> => {
  return {
    getById: async (id: string) => {
      return await get<T>(`${endpoint}/${id}`);
    },
    create: async (form: T) => {
      return await put<T>(`${endpoint}`, form);
    },
    updateById: async (id: string, form: T) => {
      return await post<T>(`${endpoint}/${id}`, form);
    },
    deleteById: async (id: string) => {
      return await del(`${endpoint}/${id}`);
    },
    getList: async (params?: ListRequestParams) => {
      return await getList<T>(`${endpoint}`, params);
    },
    getAll: async () => {
      return await getAll<T>(`${endpoint}`);
    },
    createList: async (data: T[]) => {
      return await putList<T>(`${endpoint}/batch`, data);
    },
    updateList: async (ids: string[], data: T, fields: string[]) => {
      return await postList<T>(`${endpoint}`, {ids, data: JSON.stringify(data), fields});
    },
    deleteList: async (ids: string[]) => {
      return await delList(`${endpoint}`, {ids});
    },
  };
};
