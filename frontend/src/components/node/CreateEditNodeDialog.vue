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
      :form-rules="formRules"
  >
    <template #default>
      <NodeForm/>
    </template>
  </CreateEditDialog>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {useStore} from 'vuex';
import CreateEditDialog from '@/components/dialog/CreateEditDialog.vue';
import NodeForm from '@/components/node/NodeForm.vue';
import useNode from '@/components/node/node';

export default defineComponent({
  name: 'CreateEditProjectDialog',
  components: {
    CreateEditDialog,
    NodeForm,
  },
  setup() {
    // store
    const store = useStore();

    return {
      ...useNode(store),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
