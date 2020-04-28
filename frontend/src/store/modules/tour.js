const state = {
  tourFinishSteps: {
    'spider-list': 3,
    'spider-list-add': 8,
    'spider-detail': 9,
    'spider-detail-config': 12,
    'task-list': 4,
    'task-detail': 7,
    'node-detail': 4,
    'schedule-list': 1,
    'schedule-list-add': 8,
    'setting': 2
  },
  tourSteps: {}
}

const getters = {}

const mutations = {
  SET_TOUR_STEP: (state, payload) => {
    const { tourName, step } = payload
    state.tourSteps[tourName] = step
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
