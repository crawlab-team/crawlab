<template>
  <div class="legend">
    <el-tag type="primary" size="small">
      <i class="el-icon-loading"></i>
      {{$t('Pending')}}
    </el-tag>
    <el-tag type="warning" size="small">
      <i class="el-icon-loading"></i>
      {{$t('Running')}}
    </el-tag>
    <el-tag type="success" size="small">
      <i class="el-icon-check"></i>
      {{$t('Finished')}}
    </el-tag>
    <el-tag type="danger" size="small">
      <i class="el-icon-error"></i>
      {{$t('Error')}}
    </el-tag>
    <el-tag type="info" size="small">
      <i class="el-icon-video-pause"></i>
      {{$t('Cancelled')}}
    </el-tag>
    <el-tag type="danger" size="small">
      <i class="el-icon-warning"></i>
      {{$t('Abnormal')}}
    </el-tag>
  </div>
</template>

<script>
export default {
  name: 'StatusLegend'
}
</script>

<style scoped>
  .legend .el-tag {
    margin-right: 5px;
  }
</style>
