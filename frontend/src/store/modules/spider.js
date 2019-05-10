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
  nodeStats: []
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
      lang: state.spiderForm.lang,
      col: state.spiderForm.col,
      cron: state.spiderForm.cron,
      cron_enabled: state.spiderForm.cron_enabled ? 1 : 0,
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
      cron: state.spiderForm.cron,
      cron_enabled: state.spiderForm.cron_enabled ? 1 : 0,
      site: state.spiderForm.site
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
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
