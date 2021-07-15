<template>
  <CreateEditDialog
      :type="activeDialogKey"
      :tab-name="createEditDialogTabName"
      :visible="createEditDialogVisible"
      :action-functions="actionFunctions"
      :batch-form-data="formList"
      :confirm-disabled="confirmDisabled"
      :confirm-loading="confirmLoading"
      :batch-form-fields="batchFormFields"
      :form-rules="formRules"
  >
    <template #default>
      <ProjectForm/>
    </template>
  </CreateEditDialog>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {useStore} from 'vuex';
import CreateEditDialog from '@/components/dialog/CreateEditDialog.vue';
import ProjectForm from '@/components/project/ProjectForm.vue';
import useProject from '@/components/project/project';

export default defineComponent({
  name: 'CreateEditProjectDialog',
  components: {
    CreateEditDialog,
    ProjectForm,
  },
  setup() {
    // store
    const store = useStore();

    return {
      ...useProject(store),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
