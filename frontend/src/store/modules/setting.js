import request from '../../api/request'

const state = {
  setting: {}
}

const getters = {}

const mutations = {
  SET_SETTING (state, value) {
    state.setting = value
  }
}

const actions = {
  async getSetting ({ commit }) {
    const res = await request.get('/setting')
    commit('SET_SETTING', res.data.data)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
