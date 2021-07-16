import {
  getDefaultStoreActions,
  getDefaultStoreGetters,
  getDefaultStoreMutations,
  getDefaultStoreState
} from '@/utils/store';

const state = {
  ...getDefaultStoreState<Token>('token'),
} as TokenStoreState;

const getters = {
  ...getDefaultStoreGetters<Token>(),
} as TokenStoreGetters;

const mutations = {
  ...getDefaultStoreMutations<Token>(),
} as TokenStoreMutations;

const actions = {
  ...getDefaultStoreActions<Token>('/tokens'),
} as TokenStoreActions;

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions,
} as TokenStoreModule;
