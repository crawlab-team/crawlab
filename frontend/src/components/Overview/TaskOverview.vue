<template>
  <el-row>
    <el-col :span="12" style="padding-right: 20px;">
      <el-row class="task-info-overview-wrapper wrapper">
        <h4 class="title">{{$t('Task Info')}}</h4>
        <task-info-view/>
      </el-row>
      <el-row style="border-bottom:1px solid #e4e7ed;margin:0 0 20px 0;padding-bottom:20px;"/>
    </el-col>

    <el-col :span="12">
      <el-row class="task-info-spider-wrapper wrapper">
        <h4 class="title spider-title" @click="onClickSpiderTitle">
          <i class="fa fa-search" style="margin-right: 5px"></i>
          {{$t('Spider Info')}}</h4>
        <spider-info-view :is-view="true"/>
      </el-row>
      <el-row class="task-info-node-wrapper wrapper">
        <h4 class="title node-title" @click="onClickNodeTitle">
          <i class="fa fa-search" style="margin-right: 5px"></i>
          {{$t('Node Info')}}</h4>
        <node-info-view :is-view="true"/>
      </el-row>
    </el-col>
  </el-row>
</template>

<script>
import {
  mapState
} from 'vuex'
import SpiderInfoView from '../InfoView/SpiderInfoView'
import NodeInfoView from '../InfoView/NodeInfoView'
import TaskInfoView from '../InfoView/TaskInfoView'

export default {
  name: 'SpiderOverview',
  components: {
    NodeInfoView,
    SpiderInfoView,
    TaskInfoView
  },
  data () {
    return {
      // spiderForm: {}
    }
  },
  computed: {
    ...mapState('node', [
      'nodeForm'
    ]),
    ...mapState('spider', [
      'spiderForm'
    ])
  },
  methods: {
    onClickNodeTitle () {
      this.$router.push(`/nodes/${this.nodeForm._id}`)
      this.$st.sendEv('任务详情', '概览', '点击节点详情')
    },
    onClickSpiderTitle () {
      this.$router.push(`/spiders/${this.spiderForm._id}`)
      this.$st.sendEv('任务详情', '概览', '点击爬虫详情')
    }
  },
  created () {
  }
}
</script>

<style scoped>
  .title {
    margin: 10px 0 3px 0;
    text-align: center;
    display: inline-block;
  }

  .wrapper {
    text-align: center;
  }

  .spider-form {
    padding: 10px;
  }

  .button-container {
    padding: 0 10px;
    width: 100%;
    text-align: right;
  }

  .node-title,
  .spider-title {
    cursor: pointer;
  }

  .node-title:hover,
  .spider-title:hover {
    text-decoration: underline;
  }

  .title > i {
    color: grey;
  }
</style>
