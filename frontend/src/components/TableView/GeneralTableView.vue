<template>
  <div class="general-table-view">
    <el-table
      :data="filteredData"
      :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
      border>
      <template v-for="col in columns">
        <el-table-column :key="col" :label="col" :property="col" align="center">
        </el-table-column>
      </template>
    </el-table>
    <div class="pagination">
      <el-pagination
        :current-page.sync="pagination.pageNum"
        :page-sizes="[10, 20, 50, 100]"
        :page-size.sync="pagination.pageSize"
        layout="sizes, prev, pager, next"
        :total="data.length">
      </el-pagination>
    </div>
  </div>
</template>

<script>
export default {
  name: 'GeneralTableView',
  data () {
    return {
      pagination: {
        pageNum: 1,
        pageSize: 10
      }
    }
  },
  props: {
    columns: {
      type: Array,
      default () {
        return []
      }
    },
    data: {
      type: Array,
      default () {
        return []
      }
    }
  },
  computed: {
    filteredData () {
      return this.data
        .map(d => {
          for (let k in d) {
            if (d.hasOwnProperty(k)) {
              if (d[k] === undefined || d[k] === null) continue
              if (typeof d[k] === 'object') {
                if (d[k].$oid) {
                  d[k] = d[k].$oid
                }
              }
            }
          }
          return d
        })
        .filter((d, index) => {
          // pagination
          const { pageNum, pageSize } = this.pagination
          return (pageSize * (pageNum - 1) <= index) && (index < pageSize * pageNum)
        })
    }
  }
}
</script>

<style scoped>

</style>
