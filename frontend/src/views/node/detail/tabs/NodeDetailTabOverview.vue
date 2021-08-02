<template>
  <div class="node-detail-tab-overview">
    <NodeForm readonly/>
  </div>
</template>
<script lang="ts">
import {defineComponent} from 'vue';
import NodeForm from '@/components/node/NodeForm.vue';

export default defineComponent({
  name: 'NodeDetailTabOverview',
  components: {
    NodeForm,
  },
  setup() {
    return {};
  },
});
</script>
<style lang="scss" scoped>
.node-detail-tab-overview {
  margin: 20px;
}
</style>
