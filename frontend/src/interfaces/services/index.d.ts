interface Services<T = any> {
  getById: (id: string) => Promise<ResponseWithData<T>>;
  create: (form: T) => Promise<ResponseWithData<T>>;
  updateById: (id: string, form: T) => Promise<ResponseWithData<T>>;
  deleteById: (id: string) => Promise<Response>;
  getList: (params?: ListRequestParams) => Promise<ResponseWithListData<T>>;
  getAll: () => Promise<ResponseWithListData<T>>;
  createList: (data: T[]) => Promise<ResponseWithListData<T>>;
  updateList: (ids: string[], data: T, fields: string[]) => Promise<Response>;
  deleteList: (ids: string[]) => Promise<Response>;
}
