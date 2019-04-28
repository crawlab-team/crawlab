<template>
  <el-card class="metric-card">
    <el-col :span="6" class="icon-col">
      <i :class="icon" :style="{color:color}"></i>
    </el-col>
    <el-col :span="18" class="text-col">
      <el-row>
        <label class="label">{{$t(label)}}</label>
      </el-row>
      <el-row>
        <div class="value">{{value}}</div>
      </el-row>
    </el-col>
  </el-card>
</template>

<script>
export default {
  name: 'MetricCard',
  props: {
    icon: {
      type: String,
      default: ''
    },
    label: {
      type: String,
      default: ''
    },
    value: {
      type: String,
      default: ''
    },
    type: {
      type: String,
      default: 'default'
    }
  },
  computed: {
    color () {
      if (this.type === 'primary') {
        return '#409EFF'
      } else if (this.type === 'warning') {
        return '#e6a23c'
      } else if (this.type === 'success') {
        return '#67c23a'
      } else if (this.type === 'danger') {
        return '#f56c6c'
      } else {
        return 'grey'
      }
    }
  }
}
</script>

<style scoped>
  .metric-card {
    margin-right: 20px;
  }

  .metric-card:last-child {
    margin-right: 0;
  }

  .metric-card .icon-col i {
    margin-bottom: 15px;
    font-size: 56px;
  }

  .metric-card .text-col {
    padding-left: 10px;
    height: 76px;
    text-align: center;
  }

  .metric-card .text-col label {
    font-size: 16px;
    display: block;
    height: 24px;
    color: grey;
    font-weight: 900;
  }

  .metric-card .text-col .value {
    font-size: 24px;
    display: block;
    height: 32px;
  }
</style>
