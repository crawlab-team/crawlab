<template>
  <div class="file-list-container">
    <div class="top-part">
      <!--file path-->
      <div class="file-path-container">
        <div class="left">
          <i class="el-icon-back" @click="onBack"></i>
          <div class="file-path" v-show="!isEdit">{{currentPath}}</div>
          <el-input class="file-path"
                    v-show="isEdit"
                    v-model="currentPath"
                    @change="onChange"
                    @keypress.enter.native="onChangeSubmit">
          </el-input>
        </div>
        <i class="el-icon-edit" @click="onEdit"></i>
      </div>
      <!--action-->
      <div class="action-container">
        <el-button type="success" size="mini">Choose Folder</el-button>
      </div>

    </div>

    <!--file list-->
    <template v-if="true">
      <!--<code-mirror v-model="code"/>-->
      <ul class="file-list">
        <li v-for="(item, index) in fileList" :key="index" class="item" @click="onItemClick(item)">
        <span class="item-icon">
          <i class="fa" :class="getIcon(item.type)"></i>
        </span>
          <span class="item-name">
          {{item.path}}
        </span>
        </li>
      </ul>
    </template>
    <template v-else>
    </template>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import path from 'path'
// import { codemirror } from 'vue-codemirror-lite'

export default {
  name: 'FileList',
  components: {
    // CodeMirror: codemirror
  },
  data () {
    return {
      code: 'var hello = \'world\'',
      isEdit: false
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
      this.$store.dispatch('file/getFileList', this.currentPath)
    },
    onItemClick (item) {
      if (item.type === 2) {
        this.$store.commit('file/SET_CURRENT_PATH', path.join(this.currentPath, item.path))
        this.$store.dispatch('file/getFileList', this.currentPath)
      }
    },
    onBack () {
      const sep = '/'
      let arr = this.currentPath.split(sep)
      arr.splice(arr.length - 1, 1)
      const path = arr.join(sep)
      this.$store.commit('file/SET_CURRENT_PATH', path)
      this.$store.dispatch('file/getFileList', this.currentPath)
    }
  }
}
</script>

<style scoped lang="scss">
  .file-list-container {
    height: 100%;

    .top-part {
      display: flex;
      margin-bottom: 10px;

      .file-path-container {
        width: 100%;
        padding: 5px;
        margin: 0 10px;
        border-radius: 5px;
        border: 1px solid rgba(48, 65, 86, 0.4);
        display: flex;
        justify-content: space-between;

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
        padding: 1px 5px;
        height: 24px;

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
</style>
