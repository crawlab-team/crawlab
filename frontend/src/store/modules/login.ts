import {Module} from 'vuex';

export default {
  namespaced: true,
  state: {
    isSignUp: false,
  },
  mutations: {},
  actions: {}
} as Module<LoginStoreState, RootStoreState>;
