const dialogView = {
  namespaced: true,
  state: {
    dialogType: '',
    dialogVisible: false
  },
  getters: {},
  mutations: {
    SET_DIALOG_TYPE (state, value) {
      state.dialogType = value
    },
    SET_DIALOG_VISIBLE (state, value) {
      state.dialogVisible = value
    }
  },
  actions: {}
}

export default dialogView
