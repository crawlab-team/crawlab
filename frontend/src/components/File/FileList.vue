<template>
  <div class="file-list-container">
    <div class="top-part">
      <!--back-->
      <div class="action-container">
        <el-button v-if="showFile" type="primary" size="small" style="margin-right: 10px;" @click="onBackFile">
          <font-awesome-icon :icon="['fa', 'arrow-left']"/>
          {{$t('Back')}}
        </el-button>
        <el-popover v-model="isShowDelete">
          <el-button size="small" type="default" @click="() => this.isShowDelete = false">
            {{$t('Cancel')}}
          </el-button>
          <el-button size="small" type="danger" @click="onFileDelete">
            {{$t('Confirm')}}
          </el-button>
          <template slot="reference">
            <el-button v-if="currentPath !== ''" type="danger" size="small" style="margin-right: 10px;"
                       @click="() => this.isShowDelete = true">
              <font-awesome-icon :icon="['fa', 'trash']"/>
              {{$t('Remove')}}
            </el-button>
          </template>
        </el-popover>
        <el-popover v-model="isShowRename">
          <el-input v-model="name" :placeholder="$t('Name')" style="margin-bottom: 10px"/>
          <div style="text-align: right">
            <el-button size="small" type="warning" @click="onRenameFile">
              {{$t('Confirm')}}
            </el-button>
          </div>
          <template slot="reference">
            <el-button v-if="showFile" type="warning" size="small" style="margin-right: 10px;"
                       @click="onOpenRename">
              <font-awesome-icon :icon="['fa', 'redo']"/>
              {{$t('Rename')}}
            </el-button>
          </template>
        </el-popover>
        <el-button v-if="showFile" type="success" size="small" style="margin-right: 10px;" @click="onFileSave">
          <font-awesome-icon :icon="['fa', 'save']"/>
          {{$t('Save')}}
        </el-button>
        <el-popover v-if="!showFile" v-model="isShowAdd" @hide="onHideAdd">
          <el-input v-model="name" :placeholder="$t('Name')"/>
          <div class="add-type-list">
            <el-button size="small" type="success" icon="el-icon-document-add" @click="onAddFile">
              {{$t('File')}}
            </el-button>
            <el-button size="small" type="primary" icon="el-icon-folder-add" @click="onAddDir">
              {{$t('Directory')}}
            </el-button>
          </div>
          <template slot="reference">
            <el-button type="primary" size="small" style="margin-right: 10px;">
              <font-awesome-icon :icon="['fa', 'plus']"/>
              {{$t('Create')}}
            </el-button>
          </template>
        </el-popover>
      </div>
      <!--./back-->

      <!--file path-->
      <div class="file-path-container">
        <div class="file-path">. / {{currentPath}}</div>
      </div>
      <!--./file path-->
    </div>

    <!--file list-->
    <ul v-if="!showFile" class="file-list">
      <li v-if="currentPath" class="item" @click="onBack">
        <span class="item-icon"></span>
        <span class="item-name">..</span>
      </li>
      <li v-for="(item, index) in fileList" :key="index" class="item" @click="onItemClick(item)">
            <span class="item-icon">
              <font-awesome-icon v-if="item.is_dir" :icon="['fa', 'folder']" color="rgba(3,47,98,.5)"/>
              <font-awesome-icon v-else-if="item.path.match(/\.py$/)" :icon="['fab','python']"
                                 color="rgba(3,47,98,.5)"/>
              <font-awesome-icon v-else-if="item.path.match(/\.js$/)" :icon="['fab','node-js']"
                                 color="rgba(3,47,98,.5)"/>
              <font-awesome-icon v-else-if="item.path.match(/\.(java|jar|class)$/)" :icon="['fab','java']"
                                 color="rgba(3,47,98,.5)"/>
              <font-awesome-icon v-else-if="item.path.match(/\.go$/)" :icon="['fab','go']"
                                 color="rgba(3,47,98,.5)"/>
              <font-awesome-icon v-else-if="item.path.match(/\.zip$/)" :icon="['fa','file-archive']"
                                 color="rgba(3,47,98,.5)"/>
              <font-awesome-icon v-else icon="file-alt" color="rgba(3,47,98,.5)"/>
            </span>
        <span class="item-name">
              {{item.name}}
            </span>
      </li>
    </ul>

    <file-detail v-else/>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import FileDetail from './FileDetail'

export default {
  name: 'FileList',
  components: { FileDetail },
  data () {
    return {
      code: 'var hello = \'world\'',
      isEdit: false,
      showFile: false,
      name: '',
      isShowAdd: false,
      isShowDelete: false,
      isShowRename: false
    }
  },
  computed: {
    ...mapState('file', [
      'fileList'
    ]),
    currentPath: {
      set (value) {
        this.$store.commit('file/SET_CURRENT_PATH', value)
      },
      get () {
        return this.$store.state.file.currentPath
      }
    }
  },
  methods: {
    onEdit () {
      this.isEdit = true
    },
    onItemClick (item) {
      if (item.is_dir) {
        // 目录
        this.$store.dispatch('file/getFileList', { path: item.path })
      } else {
        // 文件
        this.showFile = true
        this.$store.commit('file/SET_FILE_CONTENT', '')
        this.$store.commit('file/SET_CURRENT_PATH', item.path)
        this.$store.dispatch('file/getFileContent', { path: item.path })
      }
    },
    onBack () {
      const sep = '/'
      let arr = this.currentPath.split(sep)
      arr.splice(arr.length - 1, 1)
      const path = arr.join(sep)
      this.$store.commit('file/SET_CURRENT_PATH', path)
      this.$store.dispatch('file/getFileList', { path: this.currentPath })
    },
    async onFileSave () {
      await this.$store.dispatch('file/saveFileContent', { path: this.currentPath })
      this.$message.success(this.$t('Saved file successfully'))
    },
    onBackFile () {
      this.showFile = false
      this.onBack()
    },
    onHideAdd () {
      this.name = ''
    },
    async onAddFile () {
      if (!this.name) {
        this.$message.error(this.$t('Name cannot be empty'))
        return
      }
      const path = this.currentPath + '/' + this.name
      await this.$store.dispatch('file/addFile', { path })
      await this.$store.dispatch('file/getFileList', { path: this.currentPath })
      this.isShowAdd = false

      this.showFile = true
      this.$store.commit('file/SET_FILE_CONTENT', '')
      this.$store.commit('file/SET_CURRENT_PATH', path)
      await this.$store.dispatch('file/getFileContent', { path })
    },
    async onAddDir () {
      if (!this.name) {
        this.$message.error(this.$t('Name cannot be empty'))
        return
      }
      await this.$store.dispatch('file/addDir', { path: this.currentPath + '/' + this.name })
      await this.$store.dispatch('file/getFileList', { path: this.currentPath })
      this.isShowAdd = false
    },
    async onFileDelete () {
      await this.$store.dispatch('file/deleteFile', { path: this.currentPath })
      this.$message.success(this.$t('Deleted successfully'))
      this.isShowDelete = false
      this.onBackFile()
    },
    onOpenRename () {
      this.isShowRename = true
      const arr = this.currentPath.split('/')
      this.name = arr[arr.length - 1]
    },
    async onRenameFile () {
      let newPath
      if (this.currentPath.split('/').length === 1) {
        newPath = this.name
      } else {
        const arr = this.currentPath.split('/')
        newPath = arr[0] + '/' + this.name
      }
      await this.$store.dispatch('file/renameFile', { path: this.currentPath, newPath })
      this.$store.commit('file/SET_CURRENT_PATH', newPath)
      this.$message.success(this.$t('Renamed successfully'))
      this.isShowRename = false
    }
  },
  created () {
  }
}
</script>

<style scoped lang="scss">
  .file-list-container {
    height: 100%;
    min-height: 100%;

    .top-part {
      display: flex;
      height: 33px;
      margin-bottom: 10px;

      .file-path-container {
        width: 100%;
        padding: 5px;
        margin: 0 10px 0 0;
        border-radius: 5px;
        border: 1px solid #eaecef;
        display: flex;
        justify-content: space-between;
        color: rgba(3, 47, 98, 1);

        .left {
          width: 100%;
          display: flex;

          .el-icon-back {
            margin-right: 10px;
            cursor: pointer;
          }

          .el-input {
            /*height: 22px;*/
            width: 100%;
            line-height: 10px;
          }
        }

        .el-icon-edit {
          cursor: pointer;
        }
      }

      .action-container {
        text-align: right;
        display: flex;
        /*padding: 1px 5px;*/
        /*height: 24px;*/

        .el-button {
          margin: 0;
        }
      }
    }

    .file-list {
      padding: 0;
      margin: 0;
      list-style: none;
      height: 450px;
      overflow-y: auto;
      min-height: 100%;
      border-radius: 5px;
      border: 1px solid #eaecef;

      .item {
        padding: 10px 20px;
        cursor: pointer;
        color: #303133;

        .item-icon {
          .fa-folder {
          }
        }
      }

      .item:hover {
        background-color: rgba(48, 65, 86, 0.1);
      }
    }
  }
</style>

<style scoped>
  .file-path >>> .el-input__inner {
    font-size: 14px;
    line-height: 18px;
    height: 18px;
    border-top: none;
    border-left: none;
    border-right: none;
    border-bottom: 2px solid #409EFF;
    border-radius: 0;
  }

  .CodeMirror-line {
    padding-right: 20px;
  }

  .item {
    border-bottom: 1px solid #eaecef;
  }

  .item-icon {
    display: inline-block;
    width: 18px;
  }

  .item-name {
    font-size: 14px;
    color: rgba(3, 47, 98, 1);
  }

  .add-type-list {
    text-align: right;
    margin-top: 10px;
  }

  .add-type {
    cursor: pointer;
    font-weight: bolder;
  }
</style>
