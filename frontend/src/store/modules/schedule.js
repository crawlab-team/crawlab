import request from '../../api/request'

const state = {
  scheduleList: [],
  scheduleForm: {}
}

const getters = {}

const mutations = {
  SET_SCHEDULE_LIST (state, value) {
    state.scheduleList = value
  },
  SET_SCHEDULE_FORM (state, value) {
    state.scheduleForm = value
  }
}

const actions = {
  getScheduleList ({ state, commit }) {
    request.get('/schedules')
      .then(response => {
        commit('SET_SCHEDULE_LIST', response.data.items)
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
