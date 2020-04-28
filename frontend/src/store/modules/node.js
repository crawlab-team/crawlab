import request from '../../api/request'

const state = {
  // NodeList
  nodeList: [],
  nodeForm: {},

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
  },
  SET_NODE_SYSTEM_INFO (state, payload) {
    const { id, systemInfo } = payload
    for (let i = 0; i < state.nodeList.length; i++) {
      if (state.nodeList[i]._id === id) {
        state.nodeList[i].systemInfo = systemInfo
        break
      }
    }
  }
}

const actions = {
  getNodeList ({ state, commit }) {
    request.get('/nodes', {})
      .then(response => {
        commit('SET_NODE_LIST', response.data.data.map(d => {
          d.systemInfo = {
            os: '',
            arch: '',
            num_cpu: '',
            executables: []
          }
          return d
        }))
      })
  },
  editNode ({ state, dispatch }) {
    request.post(`/nodes/${state.nodeForm._id}`, state.nodeForm)
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
        commit('SET_NODE_FORM', response.data.data)
      })
  },
  getTaskList ({ state, commit }, id) {
    return request.get(`/nodes/${id}/tasks`)
      .then(response => {
        if (response.data.data) {
          commit('task/SET_TASK_LIST',
            response.data.data.map(d => d)
              .sort((a, b) => a.create_ts < b.create_ts ? 1 : -1),
            { root: true })
        }
      })
  },
  getNodeSystemInfo ({ state, commit }, id) {
    return request.get(`/nodes/${id}/system`)
      .then(response => {
        commit('SET_NODE_SYSTEM_INFO', { id, systemInfo: response.data.data })
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
