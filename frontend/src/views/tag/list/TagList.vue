<template>
  <ListLayout
      :action-functions="actionFunctions"
      :nav-actions="navActions"
      :pagination="tablePagination"
      :table-columns="tableColumns"
      :table-data="tableData"
      :table-total="tableTotal"
      class="tag-list"
  >
    <template #extra>
      <!-- Dialogs (handled by store) -->
      <CreateTagDialog/>
      <!-- ./Dialogs -->
    </template>
  </ListLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import CreateTagDialog from '@/components/tag/CreateEditTagDialog.vue';
import ListLayout from '@/layouts/ListLayout.vue';
import useTagList from '@/views/tag/list/tagList';

export default defineComponent({
  name: 'TagList',
  components: {
    ListLayout,
    CreateTagDialog,
  },
  setup() {
    const {
      navActions,
      tableColumns,
      tableData,
      tableTotal,
      tablePagination,
      actionFunctions,
    } = useTagList();

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
@import "../../../styles/variables.scss";

.tag-list {
  .nav-actions {
    border-bottom: none;
  }

  .content {
    background-color: $containerWhiteBg;
  }
}
</style>
