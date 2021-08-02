<template>
  <ListLayout
      :action-functions="actionFunctions"
      :nav-actions="navActions"
      :pagination="tablePagination"
      :table-columns="tableColumns"
      :table-data="tableData"
      :table-total="tableTotal"
      class="node-list"
  >
    <template #extra>
      <!-- Dialogs (handled by store) -->
      <CreateEditNodeDialog/>
      <!-- ./Dialogs -->
    </template>
  </ListLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import ListLayout from '@/layouts/ListLayout.vue';
import useNodeList from '@/views/node/list/nodeList';
import CreateEditNodeDialog from '@/components/node/CreateEditNodeDialog.vue';

export default defineComponent({
  name: 'NodeList',
  components: {
    ListLayout,
    CreateEditNodeDialog,
  },
  setup() {
    return {
      ...useNodeList(),
    };
  }
});
</script>

<style lang="scss" scoped>

</style>
