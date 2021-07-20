<template>
  <div class="file-upload">
    <div class="mode-select">
      <el-radio-group v-model="internalMode" @change="onModeChange">
        <el-radio
            v-for="{value, label} in modeOptions"
            :key="value"
            :label="value"
        >
          {{ label }}
        </el-radio>
      </el-radio-group>
    </div>

    <template v-if="mode === FILE_UPLOAD_MODE_FILES">
      <el-upload
          ref="uploadRef"
          :on-change="onFileChange"
          :http-request="() => {}"
          drag
          multiple
      >
        <i class="el-icon-upload"></i>
        <div class="el-upload__text">Drag files here, or <em>click to upload</em></div>
      </el-upload>
      <input v-bind="getInputProps()">
    </template>
    <template v-else-if="mode === FILE_UPLOAD_MODE_DIR">
      <div class="folder-upload">
        <Button @click="open">
          <i class="fa fa-folder"></i>
          Click to Select Folder to Upload
        </Button>
        <template v-if="!!dirInfo?.dirName && dirInfo?.fileCount">
          <Tag
              type="primary"
              class="info-tag"
              :label="dirInfo?.dirName"
              :icon="['fa', 'folder']"
              tooltip="Folder Name"
          />
          <Tag
              type="success"
              class="info-tag"
              :label="dirInfo?.fileCount"
              :icon="['fa', 'hashtag']"
              tooltip="Files Count"
          />
        </template>
      </div>
      <input v-bind="getInputProps()" webkitdirectory>
    </template>
  </div>
</template>

<script lang="ts">
import {defineComponent, onBeforeMount, ref, watch} from 'vue';
import {FILE_UPLOAD_MODE_DIR, FILE_UPLOAD_MODE_FILES} from '@/constants/file';
import {ElUpload} from 'element-plus/lib/el-upload/src/upload.type';
import {UploadFile} from 'element-plus/packages/upload/src/upload.type';
import Button from '@/components/button/Button.vue';
import Tag from '@/components/tag/Tag.vue';
import {plainClone} from '@/utils/object';

export default defineComponent({
  name: 'FileUpload',
  components: {
    Tag,
    Button,
  },
  props: {
    mode: {
      type: String,
    },
    getInputProps: {
      type: Function,
    },
    open: {
      type: Function,
    },
  },
  emits: [
    'mode-change',
    'files-change',
  ],
  setup(props: FileUploadProps, {emit}) {
    const modeOptions: FileUploadModeOption[] = [
      {
        label: 'Files',
        value: FILE_UPLOAD_MODE_FILES,
      },
      {
        label: 'Folder',
        value: FILE_UPLOAD_MODE_DIR,
      },
    ];
    const internalMode = ref<string>();

    const uploadRef = ref<ElUpload>();

    const dirPath = ref<string>();

    watch(() => props.mode, () => {
      internalMode.value = props.mode;
      uploadRef.value?.clearFiles();
    });

    const onFileChange = (file: UploadFile, fileList: UploadFile[]) => {
      emit('files-change', fileList.map(f => f.raw));
    };

    const clearFiles = () => {
      uploadRef.value?.clearFiles();
    };

    const onModeChange = (mode: string) => {
      emit('mode-change', mode);
    };

    onBeforeMount(() => {
      const {mode} = props;
      internalMode.value = mode;
    });

    const dirInfo = ref<FileUploadDirInfo>();

    const setDirInfo = (info: FileUploadDirInfo) => {
      console.debug(info);
      dirInfo.value = plainClone(info);
    };

    return {
      uploadRef,
      FILE_UPLOAD_MODE_FILES,
      FILE_UPLOAD_MODE_DIR,
      modeOptions,
      internalMode,
      dirPath,
      onFileChange,
      clearFiles,
      onModeChange,
      dirInfo,
      setDirInfo,
    };
  },
});
</script>

<style scoped lang="scss">
.file-upload {
  .mode-select {
    margin-bottom: 20px;
  }

  .el-upload {
    width: 100%;
  }

  .folder-upload {
    display: flex;
    align-items: center;
  }
}
</style>

<style scoped>
.file-upload >>> .el-upload,
.file-upload >>> .el-upload .el-upload-dragger {
  width: 100%;
}

.file-upload >>> .folder-upload .info-tag {
  margin-left: 10px;
}
</style>
