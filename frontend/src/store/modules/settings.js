import request from '@/api/request'

const system = {
  namespaced: true,
  state: {
    settings: {
      enable_register: false
    }
  },
  getters: {
    config (state) {
      return state.settings
    }
  },
  mutations: {
    SET_SETTINGS: (state, value) => {
      state.settings = value
      return value
    }
  },
  actions: {
    async getSettings ({ state, commit, getters }) {
      const { data: { data } } = await request.get('/settings')

      commit('SET_SETTINGS', data)
      return data
    },
    async generateRoutes ({ state, commit, getters }) {
      const settings = await this.dispatch('settings/getSettings')
      console.log(settings)
    }
  }
}
export default system
