<template>
  <div class="app-container">
    <!--filter-->
    <div class="filter">
      <el-input
        v-model="keyword"
        prefix-icon="el-icon-search"
        :placeholder="$t('Search')"
        class="filter-search"
      />
      <el-select
        v-model="filter.mainCategory"
        class="filter-category"
        :placeholder="$t('Select Main Category')"
        clearable
        filterable
        @change="onSelectMainCategory"
      >
        <el-option v-for="op in mainCategoryList" :key="op" :value="op" :label="op" />
      </el-select>
      <el-select
        v-model="filter.category"
        class="filter-category"
        :placeholder="$t('Select Category')"
        clearable
        filterable
        @change="onSelectCategory"
      >
        <el-option v-for="op in categoryList" :key="op" :value="op" :label="op" />
      </el-select>
      <el-button
        type="success"
        icon="el-icon-search"
        class="btn refresh"
        @click="onSearch"
      >
        {{ $t('Search') }}
      </el-button>
    </div>

    <!--table list-->
    <el-table
      :data="siteList"
      class="table"
      :cell-class-name="getCellClassName"
      :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
      border
    >
      <template v-for="col in columns">
        <!--<el-table-column v-if="col.name === 'category'"-->
        <!--:key="col.name"-->
        <!--:label="$t(col.label)"-->
        <!--:width="col.width"-->
        <!--:align="col.align">-->
        <!--<template slot-scope="scope">-->
        <!--<el-select v-model="scope.row[col.name]"-->
        <!--:placeholder="$t('Select')"-->
        <!--@change="onRowChange(scope.row)">-->
        <!--<el-option v-for="op in categoryList"-->
        <!--:key="op"-->
        <!--:value="op"-->
        <!--:label="op">-->
        <!--</el-option>-->
        <!--</el-select>-->
        <!--</template>-->
        <!--</el-table-column>-->
        <el-table-column
          v-if="col.name === 'domain'"
          :key="col.name"
          :label="$t(col.label)"
          :width="col.width"
          :align="col.align"
        >
          <template slot-scope="scope">
            <a
              class="domain"
              :href="'http://' + scope.row[col.name]"
              target="_blank"
              @click="onClickDomain(scope.row._id)"
            >
              {{ scope.row[col.name] }}
            </a>
          </template>
        </el-table-column>
        <el-table-column
          v-else-if="col.name === 'spider_count'"
          :key="col.name"
          :label="$t(col.label)"
          :width="col.width"
          :align="col.align"
        >
          <template slot-scope="scope">
            <div>
              <template v-if="scope.row[col.name] > 0">
                <a href="javascript:" @click="goToSpiders(scope.row._id)">
                  {{ scope.row[col.name] }}
                </a>
              </template>
              <template v-else>
                {{ scope.row[col.name] }}
              </template>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          v-else-if="col.name === 'has_robots'"
          :key="col.name"
          :label="$t(col.label)"
          :width="col.width"
          :align="col.align"
        >
          <template slot-scope="scope">
            <div>
              <template v-if="scope.row[col.name]">
                <a :href="`http://${scope.row._id}/robots.txt`" target="_blank" @click="onClickRobots(scope.row._id)">
                  Y
                </a>
              </template>
              <template v-else>
                {{ scope.row[col.name] === undefined ? 'N/A' : 'N' }}
              </template>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          v-else-if="col.name === 'home_response_time'"
          :key="col.name"
          :label="$t(col.label)"
          :width="col.width"
          :align="col.align"
        >
          <template slot-scope="scope">
            {{ scope.row[col.name] ? scope.row[col.name].toFixed(1) : 'N/A' }}
          </template>
        </el-table-column>
        <el-table-column
          v-else-if="col.name === 'home_http_status'"
          :key="col.name"
          :label="$t(col.label)"
          :width="col.width"
          :align="col.align"
        >
          <template slot-scope="scope">
            {{ scope.row[col.name] ? scope.row[col.name].toFixed(0) : 'N/A' }}
          </template>
        </el-table-column>
        <el-table-column
          v-else
          :key="col.name"
          :property="col.name"
          :label="$t(col.label)"
          :sortable="col.sortable"
          :align="col.align || 'center'"
          :width="col.width"
        />
      </template>
      <el-table-column v-if="false" :label="$t('Action')" align="left" width="120">
        <template slot-scope="scope">
          <el-tooltip :content="$t('View')" placement="top">
            <el-button type="primary" icon="el-icon-search" size="mini" @click="onView(scope.row)" />
          </el-tooltip>
          <!--<el-tooltip :content="$t('Remove')" placement="top">-->
          <!--<el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)"></el-button>-->
          <!--</el-tooltip>-->
        </template>
      </el-table-column>
    </el-table>
    <div class="pagination">
      <el-pagination
        :current-page.sync="pageNum"
        :page-sizes="[10, 20, 50, 100]"
        :page-size.sync="pageSize"
        layout="sizes, prev, pager, next"
        :total="totalCount"
        @current-change="onPageChange"
        @size-change="onPageChange"
      />
    </div>
  </div>
</template>

<script>
  import {
    mapState
  } from 'vuex'

  export default {
    name: 'SiteList',
    data() {
      return {
        // categoryList: [
        //   '新闻',
        //   '搜索引擎',
        //   '综合',
        //   '金融',
        //   '购物',
        //   '社交',
        //   '视频',
        //   '音乐',
        //   '资讯',
        //   '政企官网',
        //   '其他'
        // ],
        columns: [
          { name: 'rank', label: 'Rank', align: 'center', width: '80' },
          { name: 'name', label: 'Name', align: 'left', width: 'auto' },
          { name: 'domain', label: 'Domain', align: 'left', width: '150' },
          // { name: 'description', label: 'Description', align: 'left', width: 'auto' },
          { name: 'main_category', label: 'Main Category', align: 'center', width: '100' },
          { name: 'category', label: 'Category', align: 'center', width: '100' },
          { name: 'spider_count', label: 'Spider Count', align: 'center', width: '60' },
          { name: 'has_robots', label: 'Robots Protocol', align: 'center', width: '65' },
          { name: 'home_response_time', label: 'Home Page Response Time (sec)', align: 'right', width: '80' },
          { name: 'home_http_status', label: 'Home Page Response Status Code', align: 'right', width: '80' }
        ]
      }
    },
    computed: {
      ...mapState('site', [
        'filter',
        'siteList',
        'mainCategoryList',
        'categoryList',
        'totalCount'
      ]),
      keyword: {
        get() {
          return this.$store.state.site.keyword
        },
        set(value) {
          this.$store.commit('site/SET_KEYWORD', value)
        }
      },
      pageNum: {
        get() {
          return this.$store.state.site.pageNum
        },
        set(value) {
          this.$store.commit('site/SET_PAGE_NUM', value)
        }
      },
      pageSize: {
        get() {
          return this.$store.state.site.pageSize
        },
        set(value) {
          this.$store.commit('site/SET_PAGE_SIZE', value)
        }
      },
      mainCategory() {
        return this.filter.mainCategory
      },
      category() {
        return this.filter.category
      }
    },
    watch: {
      mainCategory() {
        // reset category
        this.filter.category = undefined

        // get category list
        this.$store.dispatch('site/getCategoryList')

        // get site list
        this.$store.dispatch('site/getSiteList')
      },
      category() {
        // get site list
        this.$store.dispatch('site/getSiteList')
      }
    },
    created() {
      this.$store.dispatch('site/getSiteList')

      this.$store.dispatch('site/getMainCategoryList')

      this.$store.dispatch('site/getCategoryList')
    },
    methods: {
      onSearch() {
        setTimeout(() => {
          this.$store.dispatch('site/getSiteList')
        }, 0)
        this.$st.sendEv('网站', '搜索')
      },
      onSelectMainCategory() {
        this.$st.sendEv('网站', '选择主类别')
      },
      onSelectCategory() {
        this.$st.sendEv('网站', '选择类别')
      },
      onClickDomain(domain) {
        this.$st.sendEv('网站', '点击域名', 'domain', domain)
      },
      onPageChange() {
        setTimeout(() => {
          this.$store.dispatch('site/getSiteList')
        }, 0)
      },
      onRowChange(row) {
        this.$store.dispatch('site/editSite', {
          id: row.domain,
          category: row.category
        })
      },
      getCellClassName({ row, columnIndex }) {
        const cls = []
        if (columnIndex === this.getColumnIndex('has_robots')) {
          cls.push('status')
          if (row.has_robots === undefined) {
            cls.push('info')
          } else if (row.has_robots) {
            cls.push('danger')
          } else {
            cls.push('success')
          }
        } else if (columnIndex === this.getColumnIndex('home_response_time')) {
          cls.push('status')
          if (row.home_response_time === undefined) {
            cls.push('info')
          } else if (row.home_response_time <= 1) {
            cls.push('success')
          } else if (row.home_response_time <= 5) {
            cls.push('primary')
          } else if (row.home_response_time <= 10) {
            cls.push('warning')
          } else {
            cls.push('danger')
          }
        } else if (columnIndex === this.getColumnIndex('home_http_status')) {
          cls.push('status')
          if (row.home_http_status === undefined) {
            cls.push('info')
          } else if (row.home_http_status >= 200 && row.home_http_status < 300) {
            cls.push('success')
          } else {
            cls.push('danger')
          }
        } else if (columnIndex === this.getColumnIndex('spider_count')) {
          cls.push('status')
          if (row.spider_count > 0) {
            cls.push('success')
          } else {
            cls.push('info')
          }
        }
        return cls.join(' ')
      },
      getColumnIndex(columnName) {
        return this.columns.map(d => d.name).indexOf(columnName)
      },
      goToSpiders(domain) {
        this.$router.push({ name: 'SpiderList', params: { domain }})
        this.$st.sendEv('网站', '点击爬虫数', 'domain', domain)
      },
      onClickRobots(domain) {
        this.$st.sendEv('网站', '点击Robots协议', 'domain', domain)
      }
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

  .table >>> .domain {
    text-decoration: underline;
  }

  .table >>> .status {
  }

  .table >>> .status.info {
    color: #909399;
    background: rgba(144, 147, 153, .1);
  }

  .table >>> .status.danger {
    color: #f56c6c;
    background: rgba(245, 108, 108, .1);
  }

  .table >>> .status.success {
    color: #67c23a;
    background: rgba(103, 194, 58, .1);
  }

  .table >>> .status.primary {
    color: #409eff;
    background: rgba(64, 158, 255, .1);
  }

  .table >>> .status.warning {
    color: #e6a23c;
    background: rgba(230, 162, 60, .1);
  }

  .table >>> a {
    text-decoration: underline;
    display: inline-block;
    width: 100%;
    height: 100%;
  }

</style>
