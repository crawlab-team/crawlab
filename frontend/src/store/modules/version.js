import request from '../../api/request'

const state = {
  version: '',
  latestRelease: {
    name: '',
    body: ''
  }
}

const getters = {}

const mutations = {
  SET_VERSION: (state, value) => {
    state.version = value
  },
  SET_LATEST_RELEASE: (state, value) => {
    state.latestRelease = value
  }
}

const actions = {
  async getLatestRelease ({ commit }) {
    const res = await request.get('/releases/latest')
    if (!res.data.error) {
      commit('SET_LATEST_RELEASE', res.data.data)
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
