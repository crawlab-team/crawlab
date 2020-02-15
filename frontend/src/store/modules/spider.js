import Vue from 'vue'
import request from '../../api/request'

const state = {
  // list of spiders
  spiderList: [],

  spiderTotal: 0,

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
  previewCrawlData: [],

  // template list
  templateList: [],

  // spider file tree
  fileTree: {}
}

const getters = {}

const mutations = {
  SET_SPIDER_TOTAL (state, value) {
    state.spiderTotal = value
  },
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
  },
  SET_SPIDER_FORM_CONFIG_SETTINGS (state, payload) {
    const settings = {}
    payload.forEach(row => {
      settings[row.name] = row.value
    })
    Vue.set(state.spiderForm.config, 'settings', settings)
  },
  SET_TEMPLATE_LIST (state, value) {
    state.templateList = value
  },
  SET_FILE_TREE (state, value) {
    state.fileTree = value
  }
}

const actions = {
  getSpiderList ({ state, commit }, params = {}) {
    return request.get('/spiders', params)
      .then(response => {
        commit('SET_SPIDER_LIST', response.data.data.list)
        commit('SET_SPIDER_TOTAL', response.data.data.total)
      })
  },
  editSpider ({ state, dispatch }) {
    return request.post(`/spiders/${state.spiderForm._id}`, state.spiderForm)
  },
  deleteSpider ({ state, dispatch }, id) {
    return request.delete(`/spiders/${id}`)
  },
  getSpiderData ({ state, commit }, id) {
    return request.get(`/spiders/${id}`)
      .then(response => {
        let data = response.data.data
        commit('SET_SPIDER_FORM', data)
      })
  },
  async getSpiderScrapySpiders ({ state, commit }, id) {
    const res = await request.get(`/spiders/${id}/scrapy/spiders`)
    state.spiderForm.spider_names = res.data.data
    commit('SET_SPIDER_FORM', state.spiderForm)
  },
  crawlSpider ({ state, dispatch }, payload) {
    const { spiderId, runType, nodeIds, param } = payload
    return request.put(`/tasks`, {
      spider_id: spiderId,
      run_type: runType,
      node_ids: nodeIds,
      param: param
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
  getDir ({ state, commit }, path) {
    const id = state.spiderForm._id
    return request.get(`/spiders/${id}/dir`)
      .then(response => {
        commit('')
      })
  },
  importGithub ({ state }) {
    const url = state.importForm.url
    return request.post('/spiders/import/github', { url })
  },
  getSpiderStats ({ state, commit }) {
    return request.get(`/spiders/${state.spiderForm._id}/stats`)
      .then(response => {
        commit('SET_OVERVIEW_STATS', response.data.data.overview)
        // commit('SET_STATUS_STATS', response.data.task_count_by_status)
        commit('SET_DAILY_STATS', response.data.data.daily)
        // commit('SET_NODE_STATS', response.data.task_count_by_node)
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
  },
  postConfigSpiderConfig ({ state }) {
    return request.post(`/config_spiders/${state.spiderForm._id}/config`, state.spiderForm.config)
  },
  saveConfigSpiderSpiderfile ({ state, rootState }) {
    const content = rootState.file.fileContent
    return request.post(`/config_spiders/${state.spiderForm._id}/spiderfile`, { content })
  },
  addConfigSpider ({ state }) {
    return request.put(`/config_spiders`, state.spiderForm)
  },
  addSpider ({ state }) {
    return request.put(`/spiders`, state.spiderForm)
  },
  async getTemplateList ({ state, commit }) {
    const res = await request.get(`/config_spiders_templates`)
    commit('SET_TEMPLATE_LIST', res.data.data)
  },
  async getScheduleList ({ state, commit }, payload) {
    const { id } = payload
    const res = await request.get(`/spiders/${id}/schedules`)
    commit('schedule/SET_SCHEDULE_LIST', res.data.data, { root: true })
  },
  async getFileTree ({ state, commit }, payload) {
    const id = payload ? payload.id : state.spiderForm._id
    const res = await request.get(`/spiders/${id}/file/tree`)
    commit('SET_FILE_TREE', res.data.data)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
