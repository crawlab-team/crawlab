<template>
  <div class="app-container">
    <!--add popup-->
    <el-dialog
      :visible.sync="dialogVisible"
      width="640px"
      :before-close="onDialogClose">
      <el-form label-width="180px"
               class="add-form"
               :model="projectForm"
               :inline-message="true"
               ref="projectForm"
               label-position="right">
        <el-form-item :label="$t('Project Name')" prop="name" required>
          <el-input id="name" v-model="projectForm.name" :placeholder="$t('Project Name')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Project Description')" prop="description">
          <el-input
            id="description"
            type="textarea"
            v-model="projectForm.description"
            :placeholder="$t('Project Description')"
          />
        </el-form-item>
        <el-form-item :label="$t('Tags')" prop="tags">
          <el-select
            id="tags"
            v-model="projectForm.tags"
            :placeholder="$t('Enter Tags')"
            allow-create
            filterable
            multiple
          >
          </el-select>
        </el-form-item>
      </el-form>
      <!--取消、保存-->
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="onDialogClose">{{$t('Cancel')}}</el-button>
        <el-button id="btn-submit" size="small" type="primary" @click="onAddSubmit">{{$t('Submit')}}</el-button>
      </span>
    </el-dialog>
    <!--./add popup-->

    <div class="action-wrapper">
      <div class="left">
        <el-select
          v-model="filter.tag"
          size="small"
          :placeholder="$t('Select Tag')"
          @change="onFilterChange"
        >
          <el-option value="" :label="$t('All Tags')"/>
          <el-option
            v-for="tag in projectTags"
            :key="tag"
            :label="tag"
            :value="tag"
          />
        </el-select>
      </div>
      <div class="right">
        <el-button
          icon="el-icon-plus"
          type="primary"
          size="small"
          @click="onAdd"
        >
          {{$t('Add Project')}}
        </el-button>
      </div>
    </div>
    <div class="content">
      <div v-if="projectList.length === 0" class="empty-list">
        {{ $t('You have no projects created. You can create a project by clicking the "Add" button.')}}
      </div>
      <ul v-else class="list">
        <li
          class="item"
          v-for="item in projectList.filter(d => d._id !== '000000000000000000000000')"
          :key="item._id"
          @click="onView(item)"
        >
          <el-card
            class="item-card"
          >
            <i v-if="!isNoProject(item)" class="btn-edit fa fa-edit" @click="onEdit(item)"></i>
            <i v-if="!isNoProject(item)" class="btn-close fa fa-trash-o" @click="onRemove(item)"></i>
            <el-row>
              <h4 class="title">{{ item.name }}</h4>
            </el-row>
            <el-row>
              <div style="display: flex; justify-content: space-between">
                <span class="spider-count">
                {{$t('Spider Count')}}: {{ item.spiders.length }}
                </span>
                <span class="owner">
                  {{item.username}}
                </span>
              </div>
            </el-row>
            <el-row class="description-wrapper">
              <div class="description">
                {{ item.description }}
              </div>
            </el-row>
            <el-row class="tags-wrapper">
              <div class="tags">
                <el-tag
                  v-for="(tag, index) in item.tags"
                  :key="index"
                  size="mini"
                  class="tag"
                >
                  {{ tag }}
                </el-tag>
              </div>
            </el-row>
          </el-card>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'ProjectList',
  data () {
    return {
      defaultTags: [],
      dialogVisible: false,
      isClickAction: false,
      filter: {
        tag: ''
      }
    }
  },
  computed: {
    ...mapState('project', [
      'projectForm',
      'projectList',
      'projectTags'
    ])
  },
  methods: {
    onDialogClose () {
      this.dialogVisible = false
    },
    onFilterChange () {
      this.$store.dispatch('project/getProjectList', this.filter)
      this.$st.sendEv('项目', '筛选项目')
    },
    onAdd () {
      this.isEdit = false
      this.dialogVisible = true
      this.$store.commit('project/SET_PROJECT_FORM', { tags: [] })
      this.$st.sendEv('项目', '添加项目')
    },
    onAddSubmit () {
      this.$refs.projectForm.validate(res => {
        if (res) {
          const form = JSON.parse(JSON.stringify(this.projectForm))
          if (this.isEdit) {
            this.$request.post(`/projects/${this.projectForm._id}`, form).then(response => {
              if (response.data.error) {
                this.$message.error(response.data.error)
                return
              }
              this.dialogVisible = false
              this.$store.dispatch('project/getProjectList')
              this.$message.success(this.$t('The project has been saved'))
            })
          } else {
            this.$request.put('/projects', form).then(response => {
              if (response.data.error) {
                this.$message.error(response.data.error)
                return
              }
              this.dialogVisible = false
              this.$store.dispatch('project/getProjectList')
              this.$message.success(this.$t('The project has been added'))
            })
          }
        }
      })
      this.$st.sendEv('项目', '提交项目')
    },
    onEdit (row) {
      this.isClickAction = true
      setTimeout(() => {
        this.isClickAction = false
      }, 100)

      this.$store.commit('project/SET_PROJECT_FORM', row)
      this.dialogVisible = true
      this.isEdit = true
      this.$st.sendEv('项目', '修改项目')
    },
    onRemove (row) {
      this.isClickAction = true
      setTimeout(() => {
        this.isClickAction = false
      }, 100)

      this.$confirm(this.$t('Are you sure to delete the project?'), this.$t('Notification'), {
        confirmButtonText: this.$t('Confirm'),
        cancelButtonText: this.$t('Cancel'),
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('project/removeProject', row._id)
          .then(() => {
            setTimeout(() => {
              this.$store.dispatch('project/getProjectList')
              this.$message.success(this.$t('The project has been removed'))
            }, 100)
          })
      }).catch(() => {
      })
      this.$st.sendEv('项目', '删除项目')
    },
    onView (row) {
      if (this.isClickAction) return

      this.$router.push({
        name: 'SpiderList',
        params: {
          project_id: row._id
        }
      })
    },
    isNoProject (row) {
      return row._id === '000000000000000000000000'
    }
  },
  async created () {
    await this.$store.dispatch('project/getProjectList', this.filter)
    await this.$store.dispatch('project/getProjectTags')
  }
}
</script>

<style scoped>
  .action-wrapper {
    display: flex;
    justify-content: space-between;
    padding-bottom: 10px;
    border-bottom: 1px solid #EBEEF5;
  }

  .list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-wrap: wrap;
  }

  .list .item {
    width: 320px;
    margin: 10px;
  }

  .list .item .item-card {
    position: relative;
    cursor: pointer;
  }

  .list .item .item-card .title {
    margin: 10px 0 0 0;
  }

  .list .item .item-card .spider-count,
  .list .item .item-card .owner {
    font-size: 12px;
    color: grey;
    font-weight: bolder;
  }

  .list .item .item-card .description-wrapper {
    padding-bottom: 5px;
    margin-bottom: 0;
    border-bottom: 1px solid #EBEEF5;
  }

  .list .item .item-card .description {
    font-size: 12px;
    line-height: 16px;
    color: grey;
  }

  .list .item .item-card .tags {
    margin-bottom: -5px;
  }

  .list .item .item-card .tags .tag {
    margin: 0 5px 5px 0;
  }

  .list .item .item-card .el-row {
    margin-bottom: 5px;
  }

  .list .item .item-card .el-row:last-child {
    margin-bottom: 0;
  }

  .list .item .item-card .btn-edit {
    z-index: 1;
    color: grey;
    position: absolute;
    top: 11px;
    right: 40px;
  }

  .list .item .item-card .btn-close {
    z-index: 1;
    color: grey;
    position: absolute;
    top: 10px;
    right: 10px;
  }

  .empty-list {
    font-size: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    height: calc(100vh - 240px);
  }

</style>
