const state = {
  lang: window.localStorage.getItem('lang') || 'en'
}

const getters = {
  lang (state) {
    if (state.lang === 'en') {
      return 'English'
    } else if (state.lang === 'zh') {
      return '中文'
    } else {
      return state.lang
    }
  }
}

const mutations = {
  SET_LANG (state, value) {
    state.lang = value
  }
}

const actions = {}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
