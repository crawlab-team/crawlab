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
      <SpiderForm/>
    </template>
  </CreateEditDialog>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {useStore} from 'vuex';
import CreateEditDialog from '@/components/dialog/CreateEditDialog.vue';
import SpiderForm from '@/components/spider/SpiderForm.vue';
import useSpider from '@/components/spider/spider';

export default defineComponent({
  name: 'CreateSpiderDialog',
  components: {
    CreateEditDialog,
    SpiderForm,
  },
  setup() {
    // store
    const store = useStore();

    return {
      ...useSpider(store),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
