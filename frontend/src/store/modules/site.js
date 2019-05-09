import request from '../../api/request'

const state = {
  tableData: [],

  // filter
  filter: {},
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
  SET_TABLE_DATA (state, value) {
    state.tableData = value
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
        commit('SET_TABLE_DATA', response.data.items)
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
