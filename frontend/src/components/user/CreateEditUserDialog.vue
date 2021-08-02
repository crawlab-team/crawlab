<template>
  <CreateEditDialog
      :action-functions="actionFunctions"
      :batch-form-data="formList"
      :batch-form-fields="batchFormFields"
      :confirm-disabled="confirmDisabled"
      :confirm-loading="confirmLoading"
      :form-rules="formRules"
      :tab-name="createEditDialogTabName"
      :type="activeDialogKey"
      :visible="createEditDialogVisible"
  >
    <template #default>
      <UserForm/>
    </template>
  </CreateEditDialog>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {useStore} from 'vuex';
import CreateEditDialog from '@/components/dialog/CreateEditDialog.vue';
import UserForm from '@/components/user/UserForm.vue';
import useUser from '@/components/user/user';

export default defineComponent({
  name: 'CreateEditUserDialog',
  components: {
    CreateEditDialog,
    UserForm,
  },
  setup() {
    // store
    const store = useStore();

    return {
      ...useUser(store),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
