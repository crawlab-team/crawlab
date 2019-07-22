<template>
  <div class="general-table-view">
    <el-table
      :data="filteredData"
      :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
      border>
      <template v-for="col in columns">
        <el-table-column :key="col" :label="col" :property="col" min-width="120">
        </el-table-column>
      </template>
    </el-table>
    <div class="pagination">
      <el-pagination
        @current-change="onPageChange"
        @size-change="onPageChange"
        :current-page.sync="pageNum"
        :page-sizes="[10, 20, 50, 100]"
        :page-size.sync="pageSize"
        layout="sizes, prev, pager, next"
        :total="total">
      </el-pagination>
    </div>
  </div>
</template>

<script>
export default {
  name: 'GeneralTableView',
  data () {
    return {}
  },
  props: {
    pageNum: {
      type: Number,
      default: 1
    },
    pageSize: {
      type: Number,
      default: 10
    },
    total: {
      type: Number,
      default: 0
    },
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
      // .map(d => d)
      // .filter((d, index) => {
      //   // pagination
      //   const pageNum = this.pageNum
      //   const pageSize = this.pageSize
      //   return (pageSize * (pageNum - 1) <= index) && (index < pageSize * pageNum)
      // })
    }
  },
  methods: {
    onPageChange () {
      this.$emit('page-change', { pageNum: this.pageNum, pageSize: this.pageSize })
    }
  }
}
</script>

<style scoped>

</style>
