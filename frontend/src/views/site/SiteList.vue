<template>
  <div class="app-container">
    <!--filter-->
    <div class="filter">
      <el-input prefix-icon="el-icon-search"
                :placeholder="$t('Search')"
                class="filter-search"
                v-model="keyword">
      </el-input>
      <el-select v-model="filter.category" class="filter-category" :placeholder="$t('Select Category')" clearable>
        <el-option v-for="op in categoryList" :key="op" :value="op" :label="op"></el-option>
      </el-select>
      <el-button type="success"
                 icon="el-icon-refresh"
                 class="btn refresh"
                 @click="onSearch">
        {{$t('Search')}}
      </el-button>
    </div>

    <!--table list-->
    <el-table :data="tableData"
              class="table"
              :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
              border>
      <template v-for="col in columns">
        <el-table-column v-if="col.name === 'category'"
                         :key="col.name"
                         :label="$t(col.label)"
                         :width="col.width"
                         :align="col.align">
          <template slot-scope="scope">
            <el-select v-model="scope.row[col.name]"
                       :placeholder="$t('Select')"
                       @change="onRowChange(scope.row)">
              <el-option v-for="op in categoryList"
                         :key="op"
                         :value="op"
                         :label="op">
              </el-option>
            </el-select>
          </template>
        </el-table-column>
        <el-table-column v-else
                         :key="col.name"
                         :property="col.name"
                         :label="$t(col.label)"
                         :sortable="col.sortable"
                         :align="col.align || 'center'"
                         :width="col.width">
        </el-table-column>
      </template>
      <el-table-column :label="$t('Action')" align="left" width="120">
        <template slot-scope="scope">
          <el-tooltip :content="$t('View')" placement="top">
            <el-button type="primary" icon="el-icon-search" size="mini" @click="onView(scope.row)"></el-button>
          </el-tooltip>
          <!--<el-tooltip :content="$t('Remove')" placement="top">-->
          <!--<el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)"></el-button>-->
          <!--</el-tooltip>-->
        </template>
      </el-table-column>
    </el-table>
    <div class="pagination">
      <el-pagination
        @current-change="onPageChange"
        @size-change="onPageChange"
        :current-page.sync="pageNum"
        :page-sizes="[10, 20, 50, 100]"
        :page-size.sync="pageSize"
        layout="sizes, prev, pager, next"
        :total="totalCount">
      </el-pagination>
    </div>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'SiteList',
  data () {
    return {
      categoryList: [
        '新闻',
        '搜索引擎',
        '综合',
        '金融',
        '购物',
        '社交',
        '视频',
        '音乐',
        '资讯',
        '政企官网',
        '其他'
      ],
      columns: [
        { name: 'rank', label: 'Rank', align: 'center', width: '80' },
        { name: 'name', label: 'Name', align: 'left', width: '120' },
        { name: 'domain', label: 'Domain', align: 'left', width: '150' },
        { name: 'description', label: 'Description', align: 'left' },
        { name: 'category', label: 'Category', align: 'center', width: '180' }
      ]
    }
  },
  computed: {
    ...mapState('site', [
      'filter',
      'tableData',
      'totalCount'
    ]),
    keyword: {
      get () {
        return this.$store.state.site.keyword
      },
      set (value) {
        this.$store.commit('site/SET_KEYWORD', value)
      }
    },
    pageNum: {
      get () {
        return this.$store.state.site.pageNum
      },
      set (value) {
        this.$store.commit('site/SET_PAGE_NUM', value)
      }
    },
    pageSize: {
      get () {
        return this.$store.state.site.pageSize
      },
      set (value) {
        this.$store.commit('site/SET_PAGE_SIZE', value)
      }
    }
  },
  methods: {
    onSearch () {
      this.$store.dispatch('site/getSiteList')
    },
    onPageChange () {
      this.$store.dispatch('site/getSiteList')
    },
    onRowChange (row) {
      this.$store.dispatch('site/editSite', {
        id: row.domain,
        category: row.category
      })
    }
  },
  created () {
    this.$store.dispatch('site/getSiteList')
  }
}
</script>

<style scoped>
  .filter {
    display: flex;
  }

  .filter .filter-search {
    width: 180px;
  }

  .filter .filter-category {
    width: 180px;
    margin-left: 20px;
  }

  .filter .btn {
    margin-left: 20px;
  }

  .table {
    margin-top: 20px;
  }

  .table >>> .el-select .el-input__inner {
    height: 32px;
  }

  .table >>> .el-select .el-select__caret {
    line-height: 32px;
  }
</style>
