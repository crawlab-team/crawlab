<template>
  <ListLayout
      :action-functions="actionFunctions"
      :nav-actions="navActions"
      :no-actions="noActions"
      :pagination="tablePagination"
      :table-columns="tableColumns"
      :table-data="tableData"
      :table-total="tableTotal"
      :selectable-function="selectableFunction"
      class="user-list"
  >
    <template #extra>
      <!-- Dialogs (handled by store) -->
      <CreateEditUserDialog/>
      <!-- ./Dialogs -->
    </template>
  </ListLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import ListLayout from '@/layouts/ListLayout.vue';
import useUserList from './userList';
import CreateEditUserDialog from '@/components/user/CreateEditUserDialog.vue';

export default defineComponent({
  name: 'UserList',
  components: {
    ListLayout,
    CreateEditUserDialog,
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
      selectableFunction,
    } = useUserList();

    return {
      navActions,
      tableColumns,
      tableData,
      tableTotal,
      tablePagination,
      actionFunctions,
      selectableFunction,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
