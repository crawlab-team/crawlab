import request from '../../api/request'

const state = {
  // NodeList
  nodeList: [],
  nodeForm: { _id: {} },

  // spider to deploy/run
  activeSpider: {}
}

const getters = {}

const mutations = {
  SET_NODE_FORM (state, value) {
    state.nodeForm = value
  },
  SET_NODE_LIST (state, value) {
    state.nodeList = value
  },
  SET_ACTIVE_SPIDER (state, value) {
    state.activeSpider = value
  }
}

const actions = {
  getNodeList ({ state, commit }) {
    request.get('/nodes', {})
      .then(response => {
        commit('SET_NODE_LIST', response.data.items)
      })
  },
  addNode ({ state, dispatch }) {
    request.put('/nodes', {
      name: state.nodeForm.name,
      ip: state.nodeForm.ip,
      port: state.nodeForm.port,
      description: state.nodeForm.description
    })
      .then(() => {
        dispatch('getNodeList')
      })
  },
  editNode ({ state, dispatch }) {
    request.post(`/nodes/${state.nodeForm._id}`, {
      name: state.nodeForm.name,
      ip: state.nodeForm.ip,
      port: state.nodeForm.port,
      description: state.nodeForm.description
    })
      .then(() => {
        dispatch('getNodeList')
      })
  },
  deleteNode ({ state, dispatch }, id) {
    request.delete(`/nodes/${id}`)
      .then(() => {
        dispatch('getNodeList')
      })
  },
  getNodeData ({ state, commit }, id) {
    request.get(`/nodes/${id}`)
      .then(response => {
        commit('SET_NODE_FORM', response.data)
      })
  },
  getDeployList ({ state, commit }, id) {
    return request.get(`/nodes/${id}/get_deploys`)
      .then(response => {
        commit('deploy/SET_DEPLOY_LIST',
          response.data.items.map(d => d)
            .sort((a, b) => a.finish_ts < b.finish_ts ? 1 : -1),
          { root: true })
      })
  },
  getTaskList ({ state, commit }, id) {
    return request.get(`/nodes/${id}/get_tasks`)
      .then(response => {
        commit('task/SET_TASK_LIST',
          response.data.items.map(d => d)
            .sort((a, b) => a.create_ts < b.create_ts ? 1 : -1),
          { root: true })
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
