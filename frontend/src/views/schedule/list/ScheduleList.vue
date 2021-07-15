<template>
  <ListLayout
      :action-functions="actionFunctions"
      :nav-actions="navActions"
      :pagination="tablePagination"
      :table-columns="tableColumns"
      :table-data="tableData"
      :table-total="tableTotal"
      :no-actions="noActions"
      class="schedule-list"
  >
    <template #extra>
      <!-- Dialogs (handled by store) -->
      <CreateEditScheduleDialog/>
      <!-- ./Dialogs -->
    </template>
  </ListLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import ListLayout from '@/layouts/ListLayout.vue';
import useScheduleList from './scheduleList';
import CreateEditScheduleDialog from '@/components/schedule/CreateEditScheduleDialog.vue';

export default defineComponent({
  name: 'ScheduleList',
  components: {
    ListLayout,
    CreateEditScheduleDialog,
  },
  props: {
    noActions: {
      type: Boolean,
      default: false,
    }
  },
  setup() {
    const {
      navActions,
      tableColumns,
      tableData,
      tableTotal,
      tablePagination,
      actionFunctions,
    } = useScheduleList();

    return {
      navActions,
      tableColumns,
      tableData,
      tableTotal,
      tablePagination,
      actionFunctions,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
