<template>
  <CreateEditDialog
      :action-functions="actionFunctions"
      :batch-form-data="formList"
      :batch-form-fields="batchFormFields"
      :confirm-disabled="confirmDisabled"
      :confirm-loading="confirmLoading"
      :tab-name="createEditDialogTabName"
      :type="activeDialogKey"
      :visible="createEditDialogVisible"
      :no-batch="noBatch"
      :title="title"
  >
    <template #default>
      <TaskForm/>
    </template>
  </CreateEditDialog>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {useStore} from 'vuex';
import CreateEditDialog from '@/components/dialog/CreateEditDialog.vue';
import useTask from '@/components/task/task';
import TaskForm from '@/components/task/TaskForm.vue';

export default defineComponent({
  name: 'CreateTaskDialog',
  components: {
    TaskForm,
    CreateEditDialog,
  },
  props: {
    title: {
      type: String,
    },
    noBatch: {
      type: Boolean,
    },
  },
  setup() {
    // store
    const store = useStore();

    return {
      ...useTask(store),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
