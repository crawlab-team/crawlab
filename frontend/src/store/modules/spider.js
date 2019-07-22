import request from '../../api/request'

const state = {
  // list of spiders
  spiderList: [],

  // active spider data
  spiderForm: {},

  // node to deploy/run
  activeNode: {},

  // upload form for importing spiders
  importForm: {
    url: '',
    type: 'github'
  },

  // spider overview stats
  overviewStats: {},

  // spider status stats
  statusStats: [],

  // spider daily stats
  dailyStats: [],

  // spider node stats
  nodeStats: [],

  // filters
  filterSite: '',

  // preview crawl data
  previewCrawlData: []
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
  },
  SET_OVERVIEW_STATS (state, value) {
    state.overviewStats = value
  },
  SET_STATUS_STATS (state, value) {
    state.statusStats = value
  },
  SET_DAILY_STATS (state, value) {
    state.dailyStats = value
  },
  SET_NODE_STATS (state, value) {
    state.nodeStats = value
  },
  SET_FILTER_SITE (state, value) {
    state.filterSite = value
  },
  SET_PREVIEW_CRAWL_DATA (state, value) {
    state.previewCrawlData = value
  }
}

const actions = {
  getSpiderList ({ state, commit }) {
    let params = {}
    if (state.filterSite) {
      params.site = state.filterSite
    }
    return request.get('/spiders', params)
      .then(response => {
        commit('SET_SPIDER_LIST', response.data.data)
      })
  },
  editSpider ({ state, dispatch }) {
    return request.post(`/spiders/${state.spiderForm._id}`, state.spiderForm)
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
        let data = response.data.data
        commit('SET_SPIDER_FORM', data)
      })
  },
  crawlSpider ({ state, dispatch }, id) {
    return request.put(`/tasks`, {
      spider_id: id
      // TODO: node_id
    })
  },
  getTaskList ({ state, commit }, id) {
    return request.get(`/spiders/${id}/tasks`)
      .then(response => {
        commit('task/SET_TASK_LIST',
          response.data.data ? response.data.data.map(d => {
            return d
          }).sort((a, b) => a.create_ts < b.create_ts ? 1 : -1) : [],
          { root: true })
      })
  },
  importGithub ({ state }) {
    const url = state.importForm.url
    return request.post('/spiders/import/github', { url })
  },
  getSpiderStats ({ state, commit }) {
    return request.get('/stats/get_spider_stats?spider_id=' + state.spiderForm._id)
      .then(response => {
        commit('SET_OVERVIEW_STATS', response.data.overview)
        commit('SET_STATUS_STATS', response.data.task_count_by_status)
        commit('SET_DAILY_STATS', response.data.daily_stats)
        commit('SET_NODE_STATS', response.data.task_count_by_node)
      })
  },
  getPreviewCrawlData ({ state, commit }) {
    return request.post(`/spiders/${state.spiderForm._id}/preview_crawl`)
      .then(response => {
        commit('SET_PREVIEW_CRAWL_DATA', response.data.items)
      })
  },
  extractFields ({ state, commit }) {
    return request.post(`/spiders/${state.spiderForm._id}/extract_fields`)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
