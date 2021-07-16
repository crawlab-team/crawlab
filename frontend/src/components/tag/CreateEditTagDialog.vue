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
  >
    <template #default>
      <TagForm/>
    </template>
  </CreateEditDialog>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {useStore} from 'vuex';
import CreateEditDialog from '@/components/dialog/CreateEditDialog.vue';
import TagForm from '@/components/tag/TagForm.vue';
import useTag from '@/components/tag/tag';

export default defineComponent({
  name: 'CreateEditTagDialog',
  components: {
    CreateEditDialog,
    TagForm,
  },
  setup() {
    // store
    const store = useStore();

    return {
      ...useTag(store),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
