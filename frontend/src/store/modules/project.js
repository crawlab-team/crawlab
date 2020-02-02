import request from '../../api/request'

const state = {
  projectForm: {},
  projectList: []
}

const getters = {}

const mutations = {
  SET_PROJECT_FORM: (state, value) => {
    state.projectForm = value
  },
  SET_PROJECT_LIST: (state, value) => {
    state.projectList = value
  }
}

const actions = {
  getProjectList ({ state, commit }) {
    request.get('/projects')
      .then(response => {
        if (response.data.data) {
          commit('SET_PROJECT_LIST', response.data.data.map(d => {
            if (!d.spiders) d.spiders = []
            return d
          }))
        }
      })
  },
  addProject ({ state }) {
    request.put('/projects', state.projectForm)
  },
  editProject ({ state }, id) {
    request.post(`/projects/${id}`, state.projectForm)
  },
  removeProject ({ state }, id) {
    request.delete(`/projects/${id}`)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
