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

    const mode = ref<string>(FILE_UPLOAD_MODE_FILES);
    const files = ref<File[]>();

    const id = computed<string>(() => route.params.id as string);

    const fileUploadRef = ref<typeof FileUpload>();

    const confirmLoading = ref<boolean>(false);
    const confirmDisabled = computed<boolean>(() => !files.value?.length);

    const onOpenFilesSettings = () => {
      store.commit(`${storeNamespace}/setEditorSettingsDialogVisible`, true);
    };

    const uploadFiles = async () => {
      if (!files.value) return;
      await Promise.all(files.value.map(f => {
        return saveFileBinary(id.value, f.name, f as File);
      }));
      await listRootDir(id.value);
    };

    const {
      getInputProps,
      open,
    } = useDropzone({
      onDrop: async (fileList: InputFile[]) => {
        if (mode.value === FILE_UPLOAD_MODE_DIR) {
          if (!fileList.length) return;
          const f = fileList[0];
          const dirName = f.path?.split('/')[0];
          const fileCount = fileList.length;
          const dirInfo = {
            dirName,
            fileCount,
          } as FileUploadDirInfo;
          console.debug(fileList, dirInfo);
          fileUploadRef.value?.setDirInfo(dirInfo);
        }
        files.value = fileList as File[];
      },
    });

    const fileUploadVisible = ref<boolean>(false);

    const onClickUpload = () => {
      fileUploadVisible.value = true;
    };

    const onModeChange = (value: string) => {
      mode.value = value;
    };

    const onFilesChange = (fileList: File[]) => {
      files.value = fileList;
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
