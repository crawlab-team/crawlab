import {
  getDefaultStoreActions,
  getDefaultStoreGetters,
  getDefaultStoreMutations,
  getDefaultStoreState
} from '@/utils/store';
import useRequest from '@/services/request';

const {
  post,
} = useRequest();

const state = {
  ...getDefaultStoreState<User>('user'),
} as UserStoreState;

const getters = {
  ...getDefaultStoreGetters<User>(),
} as UserStoreGetters;

const mutations = {
  ...getDefaultStoreMutations<User>(),
} as UserStoreMutations;

const actions = {
  ...getDefaultStoreActions<User>('/users'),
  changePassword: async (ctx: StoreActionContext, {id, password}: { id: string; password: string }) => {
    return await post(`/users/${id}/change-password`, {password});
  },
} as UserStoreActions;

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions,
} as UserStoreModule;
