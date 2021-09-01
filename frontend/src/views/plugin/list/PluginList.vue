<template>
  <ListLayout
      :action-functions="actionFunctions"
      :nav-actions="navActions"
      :pagination="tablePagination"
      :table-columns="tableColumns"
      :table-data="tableData"
      :table-total="tableTotal"
      class="plugin-list"
  >
    <template #extra>
      <!-- Dialogs (handled by store) -->
      <CreateEditPluginDialog/>
      <!-- ./Dialogs -->
    </template>
  </ListLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import ListLayout from '@/layouts/ListLayout.vue';
import usePluginList from '@/views/plugin/list/pluginList';
import CreateEditPluginDialog from '@/components/plugin/CreateEditPluginDialog.vue';

export default defineComponent({
  name: 'PluginList',
  components: {
    ListLayout,
    CreateEditPluginDialog,
  },
  setup() {
    return {
      ...usePluginList(),
    };
  }
});
</script>

<style lang="scss" scoped>

</style>
