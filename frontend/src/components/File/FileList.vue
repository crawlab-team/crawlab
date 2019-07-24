<template>
  <div class="file-list-container">
    <div class="top-part">
      <!--back-->
      <div class="action-container" v-if="showFile">
        <el-button type="primary" size="small" style="margin-right: 10px;" @click="showFile=false">
          <font-awesome-icon :icon="['fa', 'arrow-left']"/>
          {{$t('Back')}}
        </el-button>
      </div>
      <!--./back-->

      <!--file path-->
      <div class="file-path-container">
        <div class="file-path">. / {{currentPath}}</div>
      </div>
      <!--./file path-->

      <!--action-->
      <div class="action-container">
        <el-button type="primary" size="small">
          <font-awesome-icon :icon="['fa', 'upload']"/>
          {{$t('Upload')}}
        </el-button>
      </div>
      <!--./action-->

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
    getIcon (type) {
      if (type === 1) {
        return 'fa-file-o'
      } else if (type === 2) {
        return 'fa-folder'
      }
    },
    onEdit () {
      this.isEdit = true
    },
    onChange (path) {
      this.$store.commit('file/SET_CURRENT_PATH', path)
    },
    onChangeSubmit () {
      this.isEdit = false
      this.$store.dispatch('file/getFileList', { path: this.currentPath })
    },
    onItemClick (item) {
      if (item.is_dir) {
        // 目录
        this.$store.dispatch('file/getFileList', { path: item.path })
      } else {
        // 文件
        this.showFile = true
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
    }
  }
}
</script>

<style scoped lang="scss">
  .file-list-container {
    height: 100%;

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
</style>
