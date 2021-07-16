<template>
  <NavActionGroup>
    <NavActionFaIcon :icon="['fa', 'laptop-code']" tooltip="File Editor Actions"/>
    <NavActionItem>
      <FaIconButton :icon="['fa', 'upload']" tooltip="Upload Files" type="primary" @click="onOpenFiles"/>
      <input v-bind="getInputProps()">
      <FaIconButton :icon="['fa', 'cog']" tooltip="File Editor Settings" type="info" @click="onOpenFilesSettings"/>
    </NavActionItem>
  </NavActionGroup>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';
import {useStore} from 'vuex';
import NavActionGroup from '@/components/nav/NavActionGroup.vue';
import NavActionItem from '@/components/nav/NavActionItem.vue';
import FaIconButton from '@/components/button/FaIconButton.vue';
import NavActionFaIcon from '@/components/nav/NavActionFaIcon.vue';
import {useDropzone} from 'vue3-dropzone';
import useSpiderService from '@/services/spider/spiderService';
import {useRoute} from 'vue-router';

export default defineComponent({
  name: 'SpiderDetailActionsFiles',
  components: {
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

    const id = computed<string>(() => route.params.id as string);

    const {
      listRootDir,
      saveFileBinary,
    } = useSpiderService(store);

    const onOpenFilesSettings = () => {
      store.commit(`${storeNamespace}/setEditorSettingsDialogVisible`, true);
    };

    const {
      getInputProps,
      open: onOpenFiles,
    } = useDropzone({
      onDrop: async (files: InputFile[]) => {
        await Promise.all(files.map(f => {
          return saveFileBinary(id.value, f.path as string, f as File);
        }));
        await listRootDir(id.value);
      },
    });

    return {
      onOpenFilesSettings,
      getInputProps,
      onOpenFiles,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
