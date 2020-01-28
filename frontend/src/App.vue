<template>
  <div id="app">
    <dialog-view/>
    <router-view/>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import DialogView from './components/Common/DialogView'

export default {
  name: 'App',
  data () {
    return {
      msgPopup: undefined
    }
  },
  components: {
    DialogView
  },
  computed: {
    ...mapState('setting', ['setting']),
    useStats () {
      return localStorage.getItem('useStats')
    },
    uid () {
      return localStorage.getItem('uid')
    },
    sid () {
      return sessionStorage.getItem('sid')
    }
  },
  methods: {},
  async mounted () {
    window.setUseStats = (value) => {
      document.querySelector('.el-message__closeBtn').click()
      if (value === 1) {
        this.$st.sendPv('/allow_stats')
        this.$st.sendEv('全局', '允许/禁止统计', '允许')
      } else {
        this.$st.sendPv('/disallow_stats')
        this.$st.sendEv('全局', '允许/禁止统计', '禁止')
      }
      localStorage.setItem('useStats', value)
    }

    // first-time user
    if (this.useStats === undefined || this.useStats === null) {
      this.$message({
        type: 'info',
        dangerouslyUseHTMLString: true,
        showClose: true,
        duration: 0,
        message: '<p>' + this.$t('Do you allow us to collect some statistics to improve Crawlab?') + '</p>' +
          '<div style="text-align: center;margin-top: 10px;">' +
          '<button class="message-btn" onclick="setUseStats(1)">' + this.$t('Yes') + '</button>' +
          '<button class="message-btn" onclick="setUseStats(0)">' + this.$t('No') + '</button>' +
          '</div>'
      })
    }

    // set uid if first visit
    if (this.uid === undefined || this.uid === null) {
      localStorage.setItem('uid', this.$utils.encrypt.UUID())
    }

    // set session id if starting a session
    if (this.sid === undefined || this.sid === null) {
      sessionStorage.setItem('sid', this.$utils.encrypt.UUID())
    }
  }
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
</style>
