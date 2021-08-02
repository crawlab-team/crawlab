<template>
  <ListLayout
      :nav-actions="navActions"
      :table-columns="tableColumns"
      :table-data="tableData"
      :table-total="tableTotal"
      :pagination="tablePagination"
      :action-functions="actionFunctions"
      class="project-list"
  >
    <template #extra>
      <!-- Dialogs (handled by store) -->
      <CreateEditProjectDialog/>
      <!-- ./Dialogs -->
    </template>
  </ListLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import ListLayout from '@/layouts/ListLayout.vue';
import useProjectList from './projectList';
import CreateEditProjectDialog from '@/components/project/CreateEditProjectDialog.vue';

export default defineComponent({
  name: 'ProjectList',
  components: {
    ListLayout,
    CreateEditProjectDialog,
  },
  setup() {
    const {
      navActions,
      tableColumns,
      tableData,
      tableTotal,
      tablePagination,
      actionFunctions,
    } = useProjectList();

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
