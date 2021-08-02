import {
  getDefaultStoreActions,
  getDefaultStoreGetters,
  getDefaultStoreMutations,
  getDefaultStoreState
} from '@/utils/store';

const state = {
  ...getDefaultStoreState<Tag>('tag'),
} as TagStoreState;

const getters = {
  ...getDefaultStoreGetters<Tag>(),
} as TagStoreGetters;

const mutations = {
  ...getDefaultStoreMutations<Tag>(),
} as TagStoreMutations;

const actions = {
  ...getDefaultStoreActions<Tag>('/tags'),
} as TagStoreActions;

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions,
} as TagStoreModule;
