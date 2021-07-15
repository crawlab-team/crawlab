<template>
  <ListLayout
      :action-functions="actionFunctions"
      :nav-actions="navActions"
      :no-actions="noActions"
      :pagination="tablePagination"
      :table-columns="tableColumns"
      :table-data="tableData"
      :table-total="tableTotal"
      class="token-list"
  >
    <template #extra>
      <!-- Dialogs (handled by store) -->
      <!-- ./Dialogs -->
    </template>
  </ListLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import ListLayout from '@/layouts/ListLayout.vue';
import useTokenList from '@/views/token/list/tokenList';

export default defineComponent({
  name: 'TokenList',
  props: {
    noActions: {
      type: Boolean,
      default: false,
    }
  },
  components: {
    ListLayout,
  },
  setup() {
    return {
      ...useTokenList(),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
