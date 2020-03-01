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

    // set default enable_tutorial
    const enableTutorial = res.data.data.enable_tutorial
    if (localStorage.getItem('enableTutorial') === undefined) {
      localStorage.setItem('enableTutorial', enableTutorial ? '1' : '0')
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
