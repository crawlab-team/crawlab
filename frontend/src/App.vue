<template>
  <div id="app">
    <dialog-view/>
    <router-view/>
  </div>
</template>

<script>
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
    useStats () {
      return localStorage.getItem('useStats')
    }
  },
  methods: {},
  mounted () {
    window.setUseStats = (value) => {
      localStorage.setItem('useStats', value)
      document.querySelector('.el-message__closeBtn').click()
      if (value === 1) {
        window._hmt.push(['_trackPageview', '/allow_stats'])
      } else {
        window._hmt.push(['_trackPageview', '/disallow_stats'])
      }
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
