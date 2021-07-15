<template>
  <DetailLayout store-namespace="spider">
    <template #actions>
      <SpiderDetailActionsCommon/>
      <SpiderDetailActionsFiles v-if="activeTabName === 'files'"/>
    </template>
  </DetailLayout>
</template>
<script lang="ts">
import {defineComponent} from 'vue';
import SpiderDetailActionsFiles from '@/views/spider/detail/actions/SpiderDetailActionsFiles.vue';
import SpiderDetailActionsCommon from '@/views/spider/detail/actions/SpiderDetailActionsCommon.vue';
import DetailLayout from '@/layouts/DetailLayout.vue';
import useSpiderDetail from '@/views/spider/detail/spiderDetail';

export default defineComponent({
  name: 'SpiderDetail',
  components: {
    DetailLayout,
    SpiderDetailActionsCommon,
    SpiderDetailActionsFiles,
  },
  setup() {
    return {
      ...useSpiderDetail(),
    };
  },
});
</script>
<style scoped lang="scss">
</style>
