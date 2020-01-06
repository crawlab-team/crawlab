import request from '../../api/request'
const state = {
  scheduleList: [],
  scheduleForm: {
    node_ids: []
  }
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
        commit('SET_SCHEDULE_LIST', response.data.data.map(d => {
          const arr = d.cron.split(' ')
          arr.splice(0, 1)
          d.cron = arr.join(' ')
          return d
        }))
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
  enableSchedule ({ state, dispatch }, id) {
    return request.post(`/schedules/${id}/enable`)
  },
  disableSchedule ({ state, dispatch }, id) {
    return request.post(`/schedules/${id}/disable`)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
