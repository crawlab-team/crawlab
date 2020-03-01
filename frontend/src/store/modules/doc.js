import request from '../../api/request'

const state = {
  docData: []
}

const getters = {}

const mutations = {
  SET_DOC_DATA (state, value) {
    state.docData = value
  }
}

const actions = {
  async getDocData ({ commit }) {
    const res = await request.get('/docs')

    const data = JSON.parse(res.data.data.string)

    // init cache
    const cache = {}

    // iterate paths
    for (let path in data) {
      if (data.hasOwnProperty(path)) {
        const d = data[path]
        if (path.match(/\/$/)) {
          cache[path] = d
          cache[path].children = []
        } else if (path.match(/\.html$/)) {
          const parentPath = path.split('/')[0] + '/'
          cache[parentPath].children.push(d)
        }
      }
    }

    commit('SET_DOC_DATA', Object.values(cache).map(d => {
      d.level = 1
      d.label = d.title
      d.url = process.env.VUE_APP_DOC_URL + '/' + d.url
      if (d.children) {
        d.children = d.children.map(c => {
          c.level = 2
          c.label = c.title
          c.url = process.env.VUE_APP_DOC_URL + '/' + c.url
          return c
        })
      }
      return d
    }))
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
