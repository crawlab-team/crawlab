import request from '../../api/request'

const state = {
  siteList: [],

  // filter
  filter: {
    category: undefined
  },
  keyword: '',

  // pagination
  pageNum: 1,
  pageSize: 10,
  totalCount: 0
}

const getters = {}

const mutations = {
  SET_KEYWORD (state, value) {
    state.keyword = value
  },
  SET_SITE_LIST (state, value) {
    state.siteList = value
  },
  SET_PAGE_NUM (state, value) {
    state.pageNum = value
  },
  SET_PAGE_SIZE (state, value) {
    state.pageSize = value
  },
  SET_TOTAL_COUNT (state, value) {
    state.totalCount = value
  }
}

const actions = {
  editSite ({ state, dispatch }, payload) {
    const { id, category } = payload
    return request.post(`/sites/${id}`, {
      category
    })
  },
  getSiteList ({ state, commit }) {
    return request.get('/sites', {
      page_num: state.pageNum,
      page_size: state.pageSize,
      keyword: state.keyword || undefined,
      filter: {
        category: state.filter.category || undefined
      }
    })
      .then(response => {
        commit('SET_SITE_LIST', response.data.items)
        commit('SET_TOTAL_COUNT', response.data.total_count)
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
