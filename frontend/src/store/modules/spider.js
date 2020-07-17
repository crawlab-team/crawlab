import Vue from 'vue'
import request from '../../api/request'

const state = {
  // list of spiders
  spiderList: [],

  // total number of spiders
  spiderTotal: 0,

  // list of all spiders
  allSpiderList: [],

  // active spider data
  spiderForm: {},

  // spider scrapy settings
  spiderScrapySettings: [],

  // spider scrapy items
  spiderScrapyItems: [],

  // spider scrapy pipelines
  spiderScrapyPipelines: [],

  // scrapy errors
  spiderScrapyErrors: {},

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
  fileTree: {},

  // config list ts
  configListTs: undefined
}

const getters = {}

const mutations = {
  SET_SPIDER_TOTAL(state, value) {
    state.spiderTotal = value
  },
  SET_SPIDER_FORM(state, value) {
    state.spiderForm = value
  },
  SET_SPIDER_LIST(state, value) {
    state.spiderList = value
  },
  SET_ALL_SPIDER_LIST(state, value) {
    state.allSpiderList = value
  },
  SET_ACTIVE_NODE(state, value) {
    state.activeNode = value
  },
  SET_IMPORT_FORM(state, value) {
    state.importForm = value
  },
  SET_OVERVIEW_STATS(state, value) {
    state.overviewStats = value
  },
  SET_STATUS_STATS(state, value) {
    state.statusStats = value
  },
  SET_DAILY_STATS(state, value) {
    state.dailyStats = value
  },
  SET_NODE_STATS(state, value) {
    state.nodeStats = value
  },
  SET_FILTER_SITE(state, value) {
    state.filterSite = value
  },
  SET_PREVIEW_CRAWL_DATA(state, value) {
    state.previewCrawlData = value
  },
  SET_SPIDER_FORM_CONFIG_SETTINGS(state, payload) {
    const settings = {}
    payload.forEach(row => {
      settings[row.name] = row.value
    })
    Vue.set(state.spiderForm.config, 'settings', settings)
  },
  SET_TEMPLATE_LIST(state, value) {
    state.templateList = value
  },
  SET_FILE_TREE(state, value) {
    state.fileTree = value
  },
  SET_SPIDER_SCRAPY_SETTINGS(state, value) {
    state.spiderScrapySettings = value
  },
  SET_SPIDER_SCRAPY_ITEMS(state, value) {
    state.spiderScrapyItems = value
  },
  SET_SPIDER_SCRAPY_PIPELINES(state, value) {
    state.spiderScrapyPipelines = value
  },
  SET_CONFIG_LIST_TS(state, value) {
    state.configListTs = value
  },
  SET_SPIDER_SCRAPY_ERRORS(state, value) {
    for (const key in value) {
      Vue.set(state.spiderScrapyErrors, key, value[key])
    }
  }
}

const actions = {
  getSpiderList({ state, commit }, params = {}) {
    return request.get('/spiders', params)
      .then(response => {
        commit('SET_SPIDER_LIST', response.data.data.list)
        commit('SET_SPIDER_TOTAL', response.data.data.total)
      })
  },
  getAllSpiderList({ state, commit }, params = {}) {
    params.page_num = 1
    params.page_size = 99999999
    return request.get('/spiders', params)
      .then(response => {
        commit('SET_ALL_SPIDER_LIST', response.data.data.list)
      })
  },
  editSpider({ state, dispatch }) {
    return request.post(`/spiders/${state.spiderForm._id}`, state.spiderForm)
  },
  deleteSpider({ state, dispatch }, id) {
    return request.delete(`/spiders/${id}`)
  },
  getSpiderData({ state, commit }, id) {
    return request.get(`/spiders/${id}`)
      .then(response => {
        const data = response.data.data
        commit('SET_SPIDER_FORM', data)
      })
  },
  async getSpiderScrapySpiders({ state, commit }, id) {
    const res = await request.get(`/spiders/${id}/scrapy/spiders`)
    if (res.data.error) {
      commit('SET_SPIDER_SCRAPY_ERRORS', { spiders: res.data.error })
      return
    }
    state.spiderForm.spider_names = res.data.data
    commit('SET_SPIDER_FORM', state.spiderForm)
    commit('SET_SPIDER_SCRAPY_ERRORS', { spiders: '' })
  },
  async getSpiderScrapySettings({ state, commit }, id) {
    const res = await request.get(`/spiders/${id}/scrapy/settings`)
    if (res.data.error) {
      commit('SET_SPIDER_SCRAPY_ERRORS', { settings: res.data.error })
      return
    }
    commit('SET_SPIDER_SCRAPY_SETTINGS', res.data.data.map(d => {
      const key = d.key
      const value = d.value
      let type = typeof value
      if (type === 'object') {
        if (Array.isArray(value)) {
          type = 'array'
        } else {
          type = 'object'
        }
      }
      return {
        key,
        value,
        type
      }
    }))
    commit('SET_SPIDER_SCRAPY_ERRORS', { settings: '' })
  },
  async saveSpiderScrapySettings({ state }, id) {
    return request.post(`/spiders/${id}/scrapy/settings`, state.spiderScrapySettings)
  },
  async getSpiderScrapyItems({ state, commit }, id) {
    const res = await request.get(`/spiders/${id}/scrapy/items`)
    if (res.data.error) {
      commit('SET_SPIDER_SCRAPY_ERRORS', { items: res.data.error })
      return
    }
    let nodeId = 0
    commit('SET_SPIDER_SCRAPY_ITEMS', res.data.data.map(d => {
      d.id = nodeId++
      d.label = d.name
      d.level = 1
      d.isEdit = false
      d.children = d.fields.map(f => {
        return {
          id: nodeId++,
          label: f,
          level: 2,
          isEdit: false
        }
      })
      return d
    }))
    commit('SET_SPIDER_SCRAPY_ERRORS', { items: '' })
  },
  async saveSpiderScrapyItems({ state }, id) {
    return request.post(`/spiders/${id}/scrapy/items`, state.spiderScrapyItems.map(d => {
      d.name = d.label
      d.fields = d.children.map(f => f.label)
      return d
    }))
  },
  async getSpiderScrapyPipelines({ state, commit }, id) {
    const res = await request.get(`/spiders/${id}/scrapy/pipelines`)
    if (res.data.error) {
      commit('SET_SPIDER_SCRAPY_ERRORS', { pipelines: res.data.error })
      return
    }
    commit('SET_SPIDER_SCRAPY_PIPELINES', res.data.data)
    commit('SET_SPIDER_SCRAPY_ERRORS', { pipelines: '' })
  },
  async saveSpiderScrapyPipelines({ state }, id) {
    return request.post(`/spiders/${id}/scrapy/pipelines`, state.spiderScrapyPipelines)
  },
  async getSpiderScrapySpiderFilepath({ state, commit }, payload) {
    const { id, spiderName } = payload
    return request.get(`/spiders/${id}/scrapy/spider/filepath`, { spider_name: spiderName })
  },
  addSpiderScrapySpider({ state }, payload) {
    const { id, form } = payload
    return request.put(`/spiders/${id}/scrapy/spiders`, form)
  },
  crawlSpider({ state, dispatch }, payload) {
    const { spiderId, runType, nodeIds, param } = payload
    return request.put(`/tasks`, {
      spider_id: spiderId,
      run_type: runType,
      node_ids: nodeIds,
      param: param
    })
  },
  crawlSelectedSpiders({ state, dispatch }, payload) {
    const { taskParams, runType, nodeIds } = payload
    return request.post(`/spiders-run`, {
      task_params: taskParams,
      run_type: runType,
      node_ids: nodeIds
    })
  },
  getTaskList({ state, commit }, id) {
    return request.get(`/spiders/${id}/tasks`)
      .then(response => {
        commit('task/SET_TASK_LIST',
          response.data.data ? response.data.data.map(d => {
            return d
          }).sort((a, b) => a.create_ts < b.create_ts ? 1 : -1) : [],
          { root: true })
      })
  },
  getDir({ state, commit }, path) {
    const id = state.spiderForm._id
    return request.get(`/spiders/${id}/dir`)
      .then(response => {
        commit('')
      })
  },
  importGithub({ state }) {
    const url = state.importForm.url
    return request.post('/spiders/import/github', { url })
  },
  getSpiderStats({ state, commit }) {
    return request.get(`/spiders/${state.spiderForm._id}/stats`)
      .then(response => {
        commit('SET_OVERVIEW_STATS', response.data.data.overview)
        // commit('SET_STATUS_STATS', response.data.task_count_by_status)
        commit('SET_DAILY_STATS', response.data.data.daily)
        // commit('SET_NODE_STATS', response.data.task_count_by_node)
      })
  },
  getPreviewCrawlData({ state, commit }) {
    return request.post(`/spiders/${state.spiderForm._id}/preview_crawl`)
      .then(response => {
        commit('SET_PREVIEW_CRAWL_DATA', response.data.items)
      })
  },
  extractFields({ state, commit }) {
    return request.post(`/spiders/${state.spiderForm._id}/extract_fields`)
  },
  postConfigSpiderConfig({ state }) {
    return request.post(`/config_spiders/${state.spiderForm._id}/config`, state.spiderForm.config)
  },
  saveConfigSpiderSpiderfile({ state, rootState }) {
    const content = rootState.file.fileContent
    return request.post(`/config_spiders/${state.spiderForm._id}/spiderfile`, { content })
  },
  addConfigSpider({ state }) {
    return request.put(`/config_spiders`, state.spiderForm)
  },
  addSpider({ state }) {
    return request.put(`/spiders`, state.spiderForm)
  },
  async getTemplateList({ state, commit }) {
    const res = await request.get(`/config_spiders_templates`)
    commit('SET_TEMPLATE_LIST', res.data.data)
  },
  async getScheduleList({ state, commit }, payload) {
    const { id } = payload
    const res = await request.get(`/spiders/${id}/schedules`)
    let data = res.data.data
    if (data) {
      data = data.map(d => {
        const arr = d.cron.split(' ')
        arr.splice(0, 1)
        d.cron = arr.join(' ')
        return d
      })
    }
    commit('schedule/SET_SCHEDULE_LIST', data, { root: true })
  },
  async getFileTree({ state, commit }, payload) {
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
