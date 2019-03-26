<template>
  <el-row>
    <el-col :span="12" style="padding-right: 20px;">
      <el-row>
        <h4 class="title">{{$t('Task Info')}}</h4>
        <task-info-view/>
      </el-row>
      <el-row style="border-bottom:1px solid #e4e7ed;margin:0 0 20px 0;padding-bottom:20px;"/>
    </el-col>

    <el-col :span="12">
      <el-row>
        <h4 class="title spider-title" @click="onClickSpiderTitle">{{$t('Spider Info')}}</h4>
        <spider-info-view :is-view="true"/>
      </el-row>
      <el-row>
        <h4 class="title node-title" @click="onClickNodeTitle">{{$t('Node Info')}}</h4>
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
    },
    onClickSpiderTitle () {
      this.$router.push(`/spiders/${this.spiderForm._id}`)
    }
  },
  created () {
  }
}
</script>

<style scoped>
  .title {
    margin: 10px 0 3px 0;
  }

  .spider-form {
    padding: 10px;
  }

  .button-container {
    padding: 0 10px;
    width: 100%;
    text-align: right;
  }

  .title {
    text-align: center;
  }

  .node-title,
  .spider-title {
    text-decoration: underline;
    cursor: pointer;
  }
</style>
