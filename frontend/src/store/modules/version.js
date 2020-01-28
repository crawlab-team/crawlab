const state = {
  version: ''
}

const getters = {}

const mutations = {
  SET_VERSION: (state, value) => {
    state.version = value
  }
}

const actions = {}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
