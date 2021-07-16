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
      <ScheduleForm/>
    </template>
  </CreateEditDialog>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {useStore} from 'vuex';
import CreateEditDialog from '@/components/dialog/CreateEditDialog.vue';
import ScheduleForm from '@/components/schedule/ScheduleForm.vue';
import useSchedule from '@/components/schedule/schedule';

export default defineComponent({
  name: 'CreateEditScheduleDialog',
  components: {
    CreateEditDialog,
    ScheduleForm,
  },
  setup() {
    // store
    const store = useStore();

    return {
      ...useSchedule(store),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
