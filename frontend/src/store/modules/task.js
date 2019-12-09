import request from '../../api/request'

const state = {
  // TaskList
  taskList: [],
  taskListTotalCount: 0,
  taskForm: {},
  taskLog: '',
  taskResultsData: [],
  taskResultsColumns: [],
  taskResultsTotalCount: 0,
  // filter
  filter: {
    node_id: '',
    spider_id: ''
  },
  // pagination
  pageNum: 1,
  pageSize: 10,
  // results
  resultsPageNum: 1,
  resultsPageSize: 10
}

const getters = {
  taskResultsColumns (state) {
    if (!state.taskResultsData || !state.taskResultsData.length) {
      return []
    }
    const keys = []
    const item = state.taskResultsData[0]
    for (const key in item) {
      if (item.hasOwnProperty(key)) {
        keys.push(key)
      }
    }
    return keys
  }
}

const mutations = {
  SET_TASK_FORM (state, value) {
    state.taskForm = value
  },
  SET_TASK_LIST (state, value) {
    state.taskList = value
  },
  SET_TASK_LOG (state, value) {
    state.taskLog = value
  },
  SET_TASK_RESULTS_DATA (state, value) {
    state.taskResultsData = value
  },
  SET_TASK_RESULTS_COLUMNS (state, value) {
    state.taskResultsColumns = value
  },
  SET_PAGE_NUM (state, value) {
    state.pageNum = value
  },
  SET_PAGE_SIZE (state, value) {
    state.pageSize = value
  },
  SET_TASK_LIST_TOTAL_COUNT (state, value) {
    state.taskListTotalCount = value
  },
  SET_RESULTS_PAGE_NUM (state, value) {
    state.resultsPageNum = value
  },
  SET_RESULTS_PAGE_SIZE (state, value) {
    state.resultsPageSize = value
  },
  SET_TASK_RESULTS_TOTAL_COUNT (state, value) {
    state.taskResultsTotalCount = value
  }
}

const actions = {
  getTaskData ({ state, dispatch, commit }, id) {
    return request.get(`/tasks/${id}`)
      .then(response => {
        let data = response.data.data
        commit('SET_TASK_FORM', data)
        dispatch('spider/getSpiderData', data.spider_id, { root: true })
        dispatch('node/getNodeData', data.node_id, { root: true })
      })
  },
  getTaskList ({ state, commit }) {
    return request.get('/tasks', {
      page_num: state.pageNum,
      page_size: state.pageSize,
      node_id: state.filter.node_id || undefined,
      spider_id: state.filter.spider_id || undefined
    })
      .then(response => {
        commit('SET_TASK_LIST', response.data.data || [])
        commit('SET_TASK_LIST_TOTAL_COUNT', response.data.total)
      })
  },
  deleteTask ({ state, dispatch }, id) {
    return request.delete(`/tasks/${id}`)
      .then(() => {
        dispatch('getTaskList')
      })
  },
  deleteTaskMultiple ({ state }, ids) {
    return request.delete(`/tasks_multiple`, {
      ids: ids
    })
  },
  getTaskLog ({ state, commit }, id) {
    commit('SET_TASK_LOG', '')
    return request.get(`/tasks/${id}/log`)
      .then(response => {
        commit('SET_TASK_LOG', response.data.data)
      })
  },
  getTaskResults ({ state, commit }, id) {
    return request.get(`/tasks/${id}/results`, {
      page_num: state.resultsPageNum,
      page_size: state.resultsPageSize
    })
      .then(response => {
        commit('SET_TASK_RESULTS_DATA', response.data.data)
        // commit('SET_TASK_RESULTS_COLUMNS', response.data.fields)
        commit('SET_TASK_RESULTS_TOTAL_COUNT', response.data.total)
      })
  },
  async getTaskResultExcel ({ state, commit }, id) {
    const { data } = await request.request('GET', '/tasks/' + id + '/results/download', {}, {
      responseType: 'blob' // important
    })
    const downloadUrl = window.URL.createObjectURL(new Blob([data]))

    const link = document.createElement('a')

    link.href = downloadUrl

    link.setAttribute('download', 'data.csv') // any other extension

    document.body.appendChild(link)
    link.click()
    link.remove()
  },
  cancelTask ({ state, dispatch }, id) {
    return request.post(`/tasks/${id}/cancel`)
      .then(() => {
        dispatch('getTaskData', id)
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
