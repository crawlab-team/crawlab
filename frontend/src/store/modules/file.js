import request from '../../api/request'

const state = {
  currentPath: '',
  fileList: []
}

const getters = {}

const mutations = {
  SET_CURRENT_PATH (state, value) {
    state.currentPath = value
  },
  SET_FILE_LIST (state, value) {
    state.fileList = value
  }
}

const actions = {
  getFileList ({ commit }, path) {
    commit('SET_CURRENT_PATH', path)
    request.get('/files', { path })
      .then(response => {
        let list = []
        list = list.concat(response.data.folders.map(d => {
          return { path: d, type: 2 }
        }))
        list = list.concat(response.data.files.map(d => {
          return { path: d, type: 1 }
        }))
        commit('SET_FILE_LIST', list)
      })
  },
  getDefaultPath ({ commit }) {
    return new Promise((resolve, reject) => {
      request.get('/files/getDefaultPath')
        .then(response => {
          commit('SET_CURRENT_PATH', response.data.defaultPath)
          resolve(response.data.defaultPath)
        })
        .catch(reject)
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
