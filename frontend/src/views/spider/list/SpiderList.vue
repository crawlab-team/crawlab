<template>
  <ListLayout
      :nav-actions="navActions"
      :action-functions="actionFunctions"
      :pagination="tablePagination"
      :table-columns="tableColumns"
      :table-data="tableData"
      :table-total="tableTotal"
      :table-actions-prefix="tableActionsPrefix"
      :no-actions="noActions"
      class="spider-list"
  >
    <template #extra>
      <!-- Dialogs (handled by store) -->
      <CreateSpiderDialog/>
      <RunSpiderDialog v-if="activeDialogKey === 'run'"/>
      <!-- ./Dialogs -->
    </template>
  </ListLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import CreateSpiderDialog from '@/components/spider/CreateEditSpiderDialog.vue';
import ListLayout from '@/layouts/ListLayout.vue';
import useSpiderList from '@/views/spider/list/spiderList';
import RunSpiderDialog from '@/components/spider/RunSpiderDialog.vue';

export default defineComponent({
  name: 'SpiderList',
  props: {
    noActions: {
      type: Boolean,
    }
  },
  components: {
    RunSpiderDialog,
    ListLayout,
    CreateSpiderDialog,
  },
  setup() {
    const {
      navActions,
      tableColumns,
      tableData,
      tableTotal,
      tablePagination,
      actionFunctions,
      tableActionsPrefix,
      activeDialogKey,
    } = useSpiderList();

    return {
      navActions,
      tableColumns,
      tableData,
      tableTotal,
      tablePagination,
      actionFunctions,
      tableActionsPrefix,
      activeDialogKey,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../../styles/variables.scss";

.spider-list {
  .nav-actions {
    border-bottom: none;
  }

  .content {
    background-color: $containerWhiteBg;
  }
}
</style>
