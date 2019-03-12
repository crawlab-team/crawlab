import dayjs from 'dayjs'
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
  // pagination
  pageNum: 0,
  pageSize: 10,
  // results
  resultsPageNum: 0,
  resultsPageSize: 10
}

const getters = {}

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
        let data = response.data
        if (data.create_ts && data.finish_ts) {
          data.duration = dayjs(data.finish_ts.$date).diff(dayjs(data.create_ts.$date), 'second')
        }
        if (data.create_ts) data.create_ts = dayjs(data.create_ts.$date).format('YYYY-MM-DD HH:mm:ss')
        if (data.finish_ts) data.finish_ts = dayjs(data.finish_ts.$date).format('YYYY-MM-DD HH:mm:ss')
        commit('SET_TASK_FORM', data)
        dispatch('spider/getSpiderData', data.spider_id, { root: true })
        dispatch('node/getNodeData', data.node_id, { root: true })
      })
  },
  getTaskList ({ state, commit }) {
    return request.get('/tasks', {
      page_num: state.pageNum,
      page_size: state.pageSize
    })
      .then(response => {
        commit('SET_TASK_LIST', response.data.items)
        commit('SET_TASK_LIST_TOTAL_COUNT', response.data.total_count)
      })
  },
  deleteTask ({ state, dispatch }, id) {
    return request.delete(`/tasks/${id}`)
      .then(() => {
        dispatch('getTaskList')
      })
  },
  stopTask ({ state, dispatch }, id) {
    return request.post(`/tasks/${id}/stop`)
      .then(() => {
        dispatch('getTaskList')
      })
  },
  getTaskLog ({ state, commit }, id) {
    return request.get(`/tasks/${id}/get_log`)
      .then(response => {
        commit('SET_TASK_LOG', response.data.log)
      })
  },
  getTaskResults ({ state, commit }, id) {
    return request.get(`/tasks/${id}/get_results`, {
      page_num: state.resultsPageNum,
      page_size: state.resultsPageSize
    })
      .then(response => {
        commit('SET_TASK_RESULTS_DATA', response.data.items)
        commit('SET_TASK_RESULTS_COLUMNS', response.data.fields)
        commit('SET_TASK_RESULTS_COLUMNS', response.data.fields)
        commit('SET_TASK_RESULTS_TOTAL_COUNT', response.data.total_count)
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
