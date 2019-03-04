import dayjs from 'dayjs'
import request from '../../api/request'

const state = {
  // list of spiders
  spiderList: [],

  // active spider data
  spiderForm: { _id: {} },

  // node to deploy/run
  activeNode: {},

  // upload form for importing spiders
  importForm: {
    url: '',
    type: 'github'
  }
}

const getters = {}

const mutations = {
  SET_SPIDER_FORM (state, value) {
    state.spiderForm = value
  },
  SET_SPIDER_LIST (state, value) {
    state.spiderList = value
  },
  SET_ACTIVE_NODE (state, value) {
    state.activeNode = value
  },
  SET_IMPORT_FORM (state, value) {
    state.importForm = value
  }
}

const actions = {
  getSpiderList ({ state, commit }) {
    return request.get('/spiders', {})
      .then(response => {
        commit('SET_SPIDER_LIST', response.data.items)
      })
  },
  addSpider ({ state, dispatch }) {
    return request.put('/spiders', {
      name: state.spiderForm.name,
      src: state.spiderForm.src,
      cmd: state.spiderForm.cmd,
      type: state.spiderForm.type,
      lang: state.spiderForm.lang
    })
      .then(() => {
        dispatch('getSpiderList')
      })
  },
  editSpider ({ state, dispatch }) {
    return request.post(`/spiders/${state.spiderForm._id.$oid}`, {
      name: state.spiderForm.name,
      src: state.spiderForm.src,
      cmd: state.spiderForm.cmd,
      type: state.spiderForm.type,
      lang: state.spiderForm.lang
    })
      .then(() => {
        dispatch('getSpiderList')
      })
  },
  deleteSpider ({ state, dispatch }, id) {
    return request.delete(`/spiders/${id}`)
      .then(() => {
        dispatch('getSpiderList')
      })
  },
  getSpiderData ({ state, commit }, id) {
    return request.get(`/spiders/${id}`)
      .then(response => {
        commit('SET_SPIDER_FORM', response.data)
      })
  },
  deploySpider ({ state, dispatch }, { id, nodeId }) {
    return request.post(`/spiders/${id}/deploy`, {
      node_id: nodeId
    })
      .then(response => {
        console.log(response.data)
      })
  },
  crawlSpider ({ state, dispatch }, payload) {
    const { id, nodeId } = payload
    return request.post(`/spiders/${id}/crawl`, {
      node_id: nodeId
    })
      .then(response => {
        console.log(response.data)
      })
  },
  getDeployList ({ state, commit }, id) {
    return request.get(`/spiders/${id}/get_deploys`)
      .then(response => {
        commit('deploy/SET_DEPLOY_LIST',
          response.data.items.map(d => {
            if (d.finish_ts) d.finish_ts = dayjs(d.finish_ts.$date).format('YYYY-MM-DD HH:mm:ss')
            return d
          }).sort((a, b) => a.finish_ts < b.finish_ts ? 1 : -1),
          { root: true })
      })
  },
  getTaskList ({ state, commit }, id) {
    return request.get(`/spiders/${id}/get_tasks`)
      .then(response => {
        commit('task/SET_TASK_LIST',
          response.data.items.map(d => {
            if (d.create_ts) d.create_ts = dayjs(d.create_ts.$date).format('YYYY-MM-DD HH:mm:ss')
            if (d.finish_ts) d.finish_ts = dayjs(d.finish_ts.$date).format('YYYY-MM-DD HH:mm:ss')
            return d
          }).sort((a, b) => a.create_ts < b.create_ts ? 1 : -1),
          { root: true })
      })
  },
  importGithub ({ state }) {
    const url = state.importForm.url
    return request.post('/spiders/import/github', { url })
      .then(response => {
        console.log(response)
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
