import {Store} from 'vuex';

export const getDefaultService = <T>(ns: string, store: Store<RootStoreState>): Services<T> => {
  const {dispatch} = store;

  return {
    getById: (id: string) => dispatch(`${ns}/getById`, id),
    create: (form: T) => dispatch(`${ns}/create`, form),
    updateById: (id: string, form: T) => dispatch(`${ns}/updateById`, {id, form}),
    deleteById: (id: string) => dispatch(`${ns}/deleteById`, id),
    getList: (params?: ListRequestParams) => {
      if (params) {
        return dispatch(`${ns}/getListWithParams`, params);
      } else {
        return dispatch(`${ns}/getList`);
      }
    },
    getAll: () => dispatch(`${ns}/getAllList`),
    createList: (data: T[]) => dispatch(`${ns}/createList`, data),
    updateList: (ids: string[], data: T, fields: string[]) => dispatch(`${ns}/updateList`, {ids, data, fields}),
    deleteList: (ids: string[]) => dispatch(`${ns}/deleteList`, ids),
  };
};
