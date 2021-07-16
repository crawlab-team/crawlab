import {onBeforeMount, Ref} from 'vue';
import {Store} from 'vuex';
import {setupAutoUpdate} from '@/utils/auto';

export const getDefaultUseListOptions = <T = any>(navActions: Ref<ListActionGroup[]>, tableColumns: Ref<TableColumns<T>>): UseListOptions<T> => {
  return {
    navActions,
    tableColumns,
  };
};

export const setupGetAllList = (store: Store<RootStoreState>, allListNamespaces: ListStoreNamespace[]) => {
  onBeforeMount(async () => {
    await Promise.all(allListNamespaces?.map(ns => store.dispatch(`${ns}/getAllList`)) || []);
  });
};

export const setupListComponent = (ns: ListStoreNamespace, store: Store<RootStoreState>, allListNamespaces?: ListStoreNamespace[]) => {
  if (!allListNamespaces) allListNamespaces = [];

  // get all list
  setupGetAllList(store, allListNamespaces);

  // auto update
  setupAutoUpdate(async () => {
    await store.dispatch(`${ns}/getList`);
  });
};
