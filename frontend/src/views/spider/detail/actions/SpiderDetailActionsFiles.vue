<template>
  <NavActionGroup>
    <NavActionFaIcon :icon="['fa', 'laptop-code']" tooltip="File Editor Actions" />
    <NavActionItem>
      <FaIconButton :icon="['fa', 'upload']" tooltip="Upload Files" type="primary" @click="onClickUpload" />
      <FaIconButton :icon="['fa', 'cog']" tooltip="File Editor Settings" type="info" @click="onOpenFilesSettings" />
    </NavActionItem>
  </NavActionGroup>

  <Dialog
      :visible="fileUploadVisible"
      title="Files Upload"
      :confirm-loading="confirmLoading"
      :confirm-disabled="confirmDisabled"
      @close="onUploadClose"
      @confirm="onUploadConfirm"
  >
    <FileUpload
        ref="fileUploadRef"
        :mode="mode"
        :get-input-props="getInputProps"
        :open="open"
        @mode-change="onModeChange"
        @files-change="onFilesChange"
    />
  </Dialog>
</template>

<script lang="ts">
import {computed, defineComponent, ref} from 'vue';
import {useStore} from 'vuex';
import NavActionGroup from '@/components/nav/NavActionGroup.vue';
import NavActionItem from '@/components/nav/NavActionItem.vue';
import FaIconButton from '@/components/button/FaIconButton.vue';
import NavActionFaIcon from '@/components/nav/NavActionFaIcon.vue';
import {useDropzone} from 'vue3-dropzone';
import useSpiderService from '@/services/spider/spiderService';
import {useRoute} from 'vue-router';
import FileUpload from '@/components/file/FileUpload.vue';
import Dialog from '@/components/dialog/Dialog.vue';
import {ElMessage} from 'element-plus';
import {FILE_UPLOAD_MODE_DIR, FILE_UPLOAD_MODE_FILES} from '@/constants/file';
import {FileWithPath} from 'file-selector';
import {getOSPathSeparator} from '@/utils/os';

export default defineComponent({
  name: 'SpiderDetailActionsFiles',
  components: {
    Dialog,
    FileUpload,
    NavActionFaIcon,
    FaIconButton,
    NavActionGroup,
    NavActionItem,
  },
  setup() {
    // route
    const route = useRoute();

    // store
    const storeNamespace = 'file';
    const store = useStore();

    const {
      listRootDir,
      saveFileBinary,
    } = useSpiderService(store);

    const mode = ref<string>(FILE_UPLOAD_MODE_DIR);
    const files = ref<FileWithPath[]>([]);

    const id = computed<string>(() => route.params.id as string);

    const fileUploadRef = ref<typeof FileUpload>();

    const confirmLoading = ref<boolean>(false);
    const confirmDisabled = computed<boolean>(() => !files.value?.length);

    const hasMultiDir = computed<boolean>(() => {
      if (!files.value) return false;
      const set = new Set<string>();
      for (const f of files.value) {
        const lv1 = f.path?.split(getOSPathSeparator())[0] as string;
        if (!set.has(lv1)) {
          set.add(lv1);
        }
        if (set.size > 1) {
          return true;
        }
      }
      return false;
    });

    const getFilePath = (path: string) => {
      if (hasMultiDir.value) {
        return path;
      } else {
        return path.split(getOSPathSeparator()).filter((_, i) => i > 0).join(getOSPathSeparator());
      }
    };

    const setInfo = () => {
      // set file upload info
      const info = {
        fileCount: files.value.length,
        filePaths: files.value.map(f => f.path || f.name),
      } as FileUploadInfo;
      if (mode.value === FILE_UPLOAD_MODE_DIR) {
        const f = files.value[0];
        info.dirName = f.path?.split(getOSPathSeparator())[0];
      }
      fileUploadRef.value?.setInfo(info);
    };

    const onOpenFilesSettings = () => {
      store.commit(`${storeNamespace}/setEditorSettingsDialogVisible`, true);
    };

    const uploadFiles = async () => {
      if (!files.value) return;
      await Promise.all(files.value.map((f: FileWithPath) => {
        return saveFileBinary(id.value, getFilePath(f.path as string), f as File);
      }));
      await listRootDir(id.value);
    };

    const {
      getInputProps,
      open,
    } = useDropzone({
      onDrop: (fileList: InputFile[]) => {
        files.value = fileList.map(f => f as FileWithPath);

        // set file upload info
        setInfo();
      },
    });

    const fileUploadVisible = ref<boolean>(false);

    const onClickUpload = () => {
      fileUploadVisible.value = true;
    };

    const onModeChange = (value: string) => {
      mode.value = value;

      // reset file upload info
      fileUploadRef.value?.resetInfo();
    };

    const onFilesChange = (fileList: FileWithPath[]) => {
      if (!fileList.length) return;

      // set files
      files.value = fileList as FileWithPath[];

      // set file upload info
      setInfo();
    };

    const onUploadConfirm = async () => {
      confirmLoading.value = true;
      try {
        await uploadFiles();
        await ElMessage.success('Uploaded successfully');
      } catch (e) {
        await ElMessage.error(e);
      } finally {
        confirmLoading.value = false;
        fileUploadVisible.value = false;
        fileUploadRef.value?.clearFiles();
      }
    };

    const onUploadClose = () => {
      fileUploadVisible.value = false;
    };

    return {
      fileUploadRef,
      confirmLoading,
      confirmDisabled,
      onOpenFilesSettings,
      getInputProps,
      open,
      fileUploadVisible,
      onClickUpload,
      onUploadClose,
      onUploadConfirm,
      mode,
      files,
      onModeChange,
      onFilesChange,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
