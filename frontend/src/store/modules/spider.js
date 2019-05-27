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
    console.log(params)
    return request.get('/spiders', params)
      .then(response => {
        commit('SET_SPIDER_LIST', response.data.items)
      })
  },
  addSpider ({ state, dispatch }) {
    return request.put('/spiders', {
      name: state.spiderForm.name,
      col: state.spiderForm.col,
      type: 'configurable',
      site: state.spiderForm.site
    })
      .then(() => {
        dispatch('getSpiderList')
      })
  },
  editSpider ({ state, dispatch }) {
    return request.post(`/spiders/${state.spiderForm._id}`, {
      name: state.spiderForm.name,
      src: state.spiderForm.src,
      cmd: state.spiderForm.cmd,
      type: state.spiderForm.type,
      lang: state.spiderForm.lang,
      col: state.spiderForm.col,
      site: state.spiderForm.site,
      // configurable spider
      crawl_type: state.spiderForm.crawl_type,
      start_url: state.spiderForm.start_url,
      item_selector: state.spiderForm.item_selector,
      item_selector_type: state.spiderForm.item_selector_type,
      pagination_selector: state.spiderForm.pagination_selector,
      pagination_selector_type: state.spiderForm.pagination_selector_type,
      obey_robots_txt: state.spiderForm.obey_robots_txt
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
  updateSpiderEnvs ({ state }) {
    return request.post(`/spiders/${state.spiderForm._id}/update_envs`, {
      envs: JSON.stringify(state.spiderForm.envs)
    })
  },
  updateSpiderFields ({ state }) {
    return request.post(`/spiders/${state.spiderForm._id}/update_fields`, {
      fields: JSON.stringify(state.spiderForm.fields)
    })
  },
  updateSpiderDetailFields ({ state }) {
    return request.post(`/spiders/${state.spiderForm._id}/update_detail_fields`, {
      detail_fields: JSON.stringify(state.spiderForm.detail_fields)
    })
  },
  getSpiderData ({ state, commit }, id) {
    return request.get(`/spiders/${id}`)
      .then(response => {
        let data = response.data
        data.cron_enabled = !!data.cron_enabled
        commit('SET_SPIDER_FORM', data)
      })
  },
  deploySpider ({ state, dispatch }, id) {
    return request.post(`/spiders/${id}/deploy`)
      .then(response => {
        console.log(response.data)
      })
      .then(response => {
        dispatch('getSpiderData', id)
        dispatch('getSpiderList')
      })
  },
  crawlSpider ({ state, dispatch }, id) {
    return request.post(`/spiders/${id}/on_crawl`)
      .then(response => {
        console.log(response.data)
      })
  },
  getDeployList ({ state, commit }, id) {
    return request.get(`/spiders/${id}/get_deploys`)
      .then(response => {
        commit('deploy/SET_DEPLOY_LIST',
          response.data.items.map(d => {
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
  },
  deployAll () {
    return request.post('/spiders/manage/deploy_all')
      .then(response => {
        console.log(response)
      })
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
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
