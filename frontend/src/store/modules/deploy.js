import request from '../../api/request'

const state = {
  deployList: []
}

const getters = {}

const mutations = {
  SET_DEPLOY_LIST(state, value) {
    state.deployList = value
  }
}

const actions = {
  getDeployList({ state, commit }) {
    request.get('/deploys')
      .then(response => {
        commit('SET_DEPLOY_LIST', response.data.items.map(d => {
          return d
        }).sort((a, b) => a.finish_ts < b.finish_ts ? 1 : -1))
      })
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
