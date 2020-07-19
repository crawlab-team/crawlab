<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
  import {
    mapState
  } from 'vuex'
  import { getToken } from '@/utils/auth'

  export default {
    name: 'App',
    data() {
      return {
        msgPopup: undefined
      }
    },
    computed: {
      ...mapState('setting', ['setting']),
      useStats() {
        return localStorage.getItem('useStats')
      },
      uid() {
        return localStorage.getItem('uid')
      },
      sid() {
        return sessionStorage.getItem('sid')
      }
    },
    async mounted() {
      // set uid if first visit
      if (this.uid === undefined || this.uid === null) {
        localStorage.setItem('uid', this.$utils.encrypt.UUID())
      }

      // set session id if starting a session
      if (this.sid === undefined || this.sid === null) {
        sessionStorage.setItem('sid', this.$utils.encrypt.UUID())
      }

      // get latest version
      await this.$store.dispatch('version/getLatestRelease')
      if (getToken()) {
        // get user info
        await this.$store.dispatch('user/getInfo')
        // remove loading-placeholder
        const elLoading = document.querySelector('#loading-placeholder')
        elLoading.remove()

        // send visit event
        await this.$request.put('/actions', {
          type: 'visit'
        })
      }
    },
    methods: {}
  }
</script>

<style>
  .el-table .cell {
    line-height: 18px;
    font-size: 12px;
  }

  .el-table .el-table__header th,
  .el-table .el-table__body td {
    padding: 3px 0;
  }

  .el-table .el-table__header th .cell,
  .el-table .el-table__body td .cell {
    word-break: break-word;
  }

  .el-select {
    width: 100%;
  }

  .el-table .el-tag {
    font-size: 12px;
    height: 24px;
    line-height: 24px;
    font-weight: 900;
    /*padding: 0;*/
  }

  .pagination {
    margin-top: 10px;
    text-align: right;
  }

  .el-form .el-form-item {
    margin-bottom: 10px;
  }

  .message-btn {
    margin: 0 5px;
    padding: 5px 10px;
    background: transparent;
    color: #909399;
    font-size: 12px;
    border-radius: 4px;
    cursor: pointer;
    border: 1px solid #909399;
  }

  .message-btn:hover {
    opacity: 0.8;
    text-decoration: underline;
  }

  .message-btn.success {
    background: #67c23a;
    border-color: #67c23a;
    color: #fff;
  }

  .message-btn.danger {
    background: #f56c6c;
    border-color: #f56c6c;
    color: #fff;
  }

  .v-tour__target--highlighted {
    box-shadow: none !important;
    /*box-shadow: 0 0 0 4px #f56c6c !important;*/
    border: 3px solid #f56c6c !important;
  }

  .v-step__button {
    background: #67c23a !important;
    border: none !important;
    color: white !important;
  }

  .v-step__button:hover {
    background: #67c23a !important;
    border: none !important;
    color: white !important;
    opacity: 0.9 !important;
  }
</style>
