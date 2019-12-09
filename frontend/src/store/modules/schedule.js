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
        commit('SET_SCHEDULE_LIST', response.data.data)
      })
  },
  addSchedule ({ state }) {
    request.put('/schedules', state.scheduleForm)
  },
  editSchedule ({ state }, id) {
    request.post(`/schedules/${id}`, state.scheduleForm)
  },
  removeSchedule ({ state }, id) {
    request.delete(`/schedules/${id}`)
  },
  stopSchedule ({ state, dispatch }, id) {
    return request.post(`/schedules/${id}/stop`)
  },
  runSchedule ({ state, dispatch }, id) {
    return request.post(`/schedules/${id}/run`)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
