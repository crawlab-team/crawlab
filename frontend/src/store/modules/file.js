import request from '../../api/request'

const state = {
  currentPath: '',
  fileList: [],
  fileContent: ''
}

const getters = {}

const mutations = {
  SET_CURRENT_PATH (state, value) {
    state.currentPath = value
  },
  SET_FILE_LIST (state, value) {
    state.fileList = value
  },
  SET_FILE_CONTENT (state, value) {
    state.fileContent = value
  }
}

const actions = {
  getFileList ({ commit, rootState }, payload) {
    const { path } = payload
    const spiderId = rootState.spider.spiderForm._id
    commit('SET_CURRENT_PATH', path)
    request.get(`/spiders/${spiderId}/dir`, { path })
      .then(response => {
        commit(
          'SET_FILE_LIST',
          response.data.data
            .sort((a, b) => a.name > b.name ? -1 : 1)
            .sort((a, b) => a.is_dir > b.is_dir ? -1 : 1)
        )
      })
  },
  getFileContent ({ commit, rootState }, payload) {
    const { path } = payload
    const spiderId = rootState.spider.spiderForm._id
    request.get(`/spiders/${spiderId}/file`, { path })
      .then(response => {
        commit('SET_FILE_CONTENT', response.data.data)
      })
  },

}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
