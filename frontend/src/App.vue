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
  components: {
    DialogView
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
</style>
