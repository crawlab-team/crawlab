<template>
  <div class="spider-detail-tab-schedules">
    <ScheduleList no-actions/>
  </div>
</template>
<script lang="ts">
import {computed, defineComponent, onBeforeMount, onBeforeUnmount} from 'vue';
import {useRoute} from 'vue-router';
import {useStore} from 'vuex';
import {FILTER_OP_EQUAL} from '@/constants/filter';
import ScheduleList from '@/views/schedule/list/ScheduleList.vue';

export default defineComponent({
  name: 'SpiderDetailTabSchedules',
  components: {
    ScheduleList,
  },
  setup() {
    // route
    const route = useRoute();

    // store
    const store = useStore();

    // id
    const id = computed<string>(() => route.params.id as string);

    onBeforeMount(() => {
      // set filter
      store.commit(`schedule/setTableListFilter`, [{
        key: 'spider_id',
        op: FILTER_OP_EQUAL,
        value: id.value,
      }]);
    });

    onBeforeUnmount(() => {
      store.commit(`schedule/resetTableListFilter`);
      store.commit(`schedule/resetTableData`);
    });

    return {};
  },
});
</script>

<style scoped>
.spider-detail-tab-schedules >>> .el-table {
  border: none;
}
</style>
