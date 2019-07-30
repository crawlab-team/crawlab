import request from '../../api/request'

const user = {
  namespaced: true,

  state: {
    // token: getToken(),
    name: '',
    avatar: '',
    roles: [],
    userList: [],
    userForm: {},
    userInfo: undefined,
    adminPaths: [
      '/users'
    ],
    pageNum: 1,
    pageSize: 10,
    totalCount: 0
  },

  getters: {
    userInfo (state) {
      if (state.userInfo) return state.userInfo
      const userInfoStr = window.localStorage.getItem('user_info')
      if (!userInfoStr) return {}
      return JSON.parse(userInfoStr)
    },
    token () {
      return window.localStorage.getItem('token')
    }
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token
    },
    SET_NAME: (state, name) => {
      state.name = name
    },
    SET_AVATAR: (state, avatar) => {
      state.avatar = avatar
    },
    SET_ROLES: (state, roles) => {
      state.roles = roles
    },
    SET_USER_LIST: (state, value) => {
      state.userList = value
    },
    SET_USER_FORM: (state, value) => {
      state.userForm = value
    },
    SET_USER_INFO: (state, value) => {
      state.userInfo = value
    },
    SET_PAGE_NUM: (state, value) => {
      state.pageNum = value
    },
    SET_PAGE_SIZE: (state, value) => {
      state.pageSize = value
    },
    SET_TOTAL_COUNT: (state, value) => {
      state.totalCount = value
    }
  },

  actions: {
    // 登录
    login ({ commit }, userInfo) {
      const username = userInfo.username.trim()
      return new Promise((resolve, reject) => {
        request.post('/login', { username, password: userInfo.password })
          .then(response => {
            const token = response.data.data
            commit('SET_TOKEN', token)
            window.localStorage.setItem('token', token)
            resolve()
          })
          .catch(error => {
            reject(error)
          })
      })
    },

    // 获取用户信息
    getInfo ({ commit, state }) {
      return request.get('/me')
        .then(response => {
          commit('SET_USER_INFO', response.data.data)
          window.localStorage.setItem('user_info', JSON.stringify(response.data.data))
        })
    },

    // 注册
    register ({ dispatch, commit, state }, userInfo) {
      return new Promise((resolve, reject) => {
        request.put('/users', { username: userInfo.username, password: userInfo.password })
          .then(() => {
            resolve()
          })
          .catch(r => {
            reject(r.response.data.error)
          })
      })
    },

    // 登出
    logout ({ commit, state }) {
      return new Promise((resolve, reject) => {
        window.localStorage.removeItem('token')
        window.localStorage.removeItem('user_info')
        commit('SET_USER_INFO', undefined)
        commit('SET_TOKEN', '')
        commit('SET_ROLES', [])
        resolve()
      })
    },

    // 获取用户列表
    getUserList ({ commit, state }) {
      return new Promise((resolve, reject) => {
        request.get('/users', {
          page_num: state.pageNum,
          page_size: state.pageSize
        })
          .then(response => {
            commit('SET_USER_LIST', response.data.data)
            commit('SET_TOTAL_COUNT', response.data.total)
          })
      })
    },

    // 删除用户
    deleteUser ({ state }, id) {
      return request.delete(`/users/${id}`)
    },

    // 编辑用户
    editUser ({ state }) {
      return request.post(`/users/${state.userForm._id}`, state.userForm)
    }
  }
}

export default user
