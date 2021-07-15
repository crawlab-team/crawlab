<template>
  <div class="project-detail-tab-spiders">
    <SpiderList no-actions/>
  </div>
</template>
<script lang="ts">
import {computed, defineComponent, onBeforeMount, onBeforeUnmount} from 'vue';
import {useRoute} from 'vue-router';
import {useStore} from 'vuex';
import {FILTER_OP_EQUAL} from '@/constants/filter';
import SpiderList from '@/views/spider/list/SpiderList.vue';

export default defineComponent({
  name: 'ProjectDetailTabSpiders',
  components: {
    SpiderList,
  },
  setup() {
    // route
    const route = useRoute();

    // store
    const ns = 'project';
    const store = useStore();

    // id
    const id = computed<string>(() => route.params.id as string);

    onBeforeMount(() => {
      store.commit(`spider/setTableListFilter`, [{
        key: 'project_id',
        op: FILTER_OP_EQUAL,
        value: id.value,
      } as FilterConditionData]);
    });

    onBeforeUnmount(() => {
      store.commit(`spider/resetTableListFilter`);
      store.commit(`spider/resetTableData`);
    });

    return {};
  },
});
</script>
<style lang="scss" scoped>
.project-detail-tab-overview {
  margin: 20px;
}
</style>
