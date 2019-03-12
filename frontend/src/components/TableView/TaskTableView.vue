<template>
  <div class="task-table-view">
    <el-row class="title-wrapper">
      <h5 class="title">{{title}}</h5>
      <el-button type="success" plain class="small-btn" size="mini" icon="fa fa-refresh" @click="onRefresh"></el-button>
    </el-row>
    <el-table border height="480px" :data="taskList">
      <el-table-column property="node" label="Node" width="220" align="center">
        <template slot-scope="scope">
          <a class="a-tag" @click="onClickNode(scope.row)">{{scope.row.node_id}}</a>
        </template>
      </el-table-column>
      <el-table-column property="spider_name" label="Spider" width="80" align="center">
        <template slot-scope="scope">
          <a class="a-tag" @click="onClickSpider(scope.row)">{{scope.row.spider_name}}</a>
        </template>
      </el-table-column>
      <el-table-column label="Status"
                       align="center"
                       width="100">
        <template slot-scope="scope">
          <el-tag type="success" v-if="scope.row.status === 'SUCCESS'">SUCCESS</el-tag>
          <el-tag type="warning" v-else-if="scope.row.status === 'PENDING'">PENDING</el-tag>
          <el-tag type="danger" v-else-if="scope.row.status === 'FAILURE'">FAILURE</el-tag>
          <el-tag type="info" v-else>{{scope.row['status']}}</el-tag>
        </template>
      </el-table-column>
      <!--<el-table-column property="create_ts" label="Create Time" width="auto" align="center"></el-table-column>-->
      <el-table-column property="create_ts" label="Create Time" width="auto" align="center">
        <template slot-scope="scope">
          <a href="javascript:" class="a-tag" @click="onClickTask(scope.row)">{{scope.row.create_ts}}</a>
        </template>
      </el-table-column>
    </el-table>
  </div>

</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'TaskTableView',
  props: {
    title: String
  },
  computed: {
    ...mapState('spider', [
      'spiderForm'
    ]),
    ...mapState('task', [
      'taskList'
    ])
  },
  methods: {
    onClickSpider (row) {
      this.$router.push(`/spiders/${row.spider_id}`)
    },
    onClickNode (row) {
      this.$router.push(`/nodes/${row.node_id}`)
    },
    onClickTask (row) {
      this.$router.push(`/tasks/${row._id}`)
    },
    onRefresh () {
      if (this.$route.path.split('/')[1] === 'spiders') {
        this.$store.dispatch('spider/getTaskList', this.$route.params.id)
      } else if (this.$route.path.split('/')[1] === 'nodes') {
        this.$store.dispatch('node/getTaskList', this.$route.params.id)
      }
    }
  }
}
</script>

<style scoped>
  .task-table-view {
    margin-bottom: 10px;
  }

  .el-table .a-tag {
    text-decoration: underline;
  }

  .title {
    margin: 10px 0 3px 0;
    float: left;
  }

  .small-btn {
    float: right;
    width: 24px;
    margin: 0;
    padding: 5px;
  }
</style>
