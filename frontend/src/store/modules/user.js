import request from '../../api/request'

const user = {
  namespaced: true,

  state: {
    // token: getToken(),
    name: '',
    avatar: '',
    roles: []
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
    // getInfo ({ commit, state }) {
    //   return new Promise((resolve, reject) => {
    //     getInfo(state.token).then(response => {
    //       const data = response.data
    //       if (data.roles && data.roles.length > 0) { // 验证返回的roles是否是一个非空数组
    //         commit('SET_ROLES', data.roles)
    //       } else {
    //         reject(new Error('getInfo: roles must be a non-null array !'))
    //       }
    //       commit('SET_NAME', data.name)
    //       commit('SET_AVATAR', data.avatar)
    //       resolve(response)
    //     }).catch(error => {
    //       reject(error)
    //     })
    //   })
    // },

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
        commit('SET_TOKEN', '')
        commit('SET_ROLES', [])
        resolve()
      })
    }
  }
}

export default user
