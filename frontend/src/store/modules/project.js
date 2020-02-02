import request from '../../api/request'

const state = {
  projectForm: {},
  projectList: [],
  projectTags: []
}

const getters = {}

const mutations = {
  SET_PROJECT_FORM: (state, value) => {
    state.projectForm = value
  },
  SET_PROJECT_LIST: (state, value) => {
    state.projectList = value
  },
  SET_PROJECT_TAGS: (state, value) => {
    state.projectTags = value
  }
}

const actions = {
  getProjectList ({ state, commit }, payload) {
    request.get('/projects', payload)
      .then(response => {
        if (response.data.data) {
          commit('SET_PROJECT_LIST', response.data.data.map(d => {
            if (!d.spiders) d.spiders = []
            return d
          }))
        }
      })
  },
  getProjectTags ({ state, commit }) {
    request.get('/projects/tags')
      .then(response => {
        if (response.data.data) {
          commit('SET_PROJECT_TAGS', response.data.data.map(d => d.tag))
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
