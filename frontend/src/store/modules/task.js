import request from '../../api/request'
import utils from '../../utils'

const state = {
  // TaskList
  taskList: [],
  taskListTotalCount: 0,
  taskForm: {},
  taskLog: '',
  currentLogIndex: 0,
  taskResultsData: [],
  taskResultsColumns: [],
  taskResultsTotalCount: 0,
  // filter
  filter: {
    node_id: '',
    spider_id: '',
    status: ''
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
  },
  logData (state) {
    const data = state.taskLog.split('\n')
      .map((d, i) => {
        return {
          index: i + 1,
          data: d,
          active: state.currentLogIndex === i + 1
        }
      })
    if (state.taskForm && state.taskForm.status === 'running') {
      data.push({
        index: data.length + 1,
        data: '###LOG_END###'
      })
      data.push({
        index: data.length + 1,
        data: ''
      })
    }
    return data
  },
  errorLogData (state, getters) {
    return getters.logData.filter(d => {
      return d.data.match(utils.log.errorRegex)
    })
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
  SET_CURRENT_LOG_INDEX (state, value) {
    state.currentLogIndex = value
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
      spider_id: state.filter.spider_id || undefined,
      status: state.filter.status || undefined
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
    return new Promise(resolve => {
      request.post(`/tasks/${id}/cancel`)
        .then(res => {
          dispatch('getTaskData', id)
          resolve(res)
        })
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
