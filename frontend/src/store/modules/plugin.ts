import {
  getDefaultStoreActions,
  getDefaultStoreGetters,
  getDefaultStoreMutations,
  getDefaultStoreState
} from '@/utils/store';

type Plugin = CPlugin;

const state = {
  ...getDefaultStoreState<Plugin>('plugin'),
} as PluginStoreState;

const getters = {
  ...getDefaultStoreGetters<Plugin>(),
} as PluginStoreGetters;

const mutations = {
  ...getDefaultStoreMutations<Plugin>(),
} as PluginStoreMutations;

const actions = {
  ...getDefaultStoreActions<Plugin>('/plugins'),
} as PluginStoreActions;

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions,
} as PluginStoreModule;
