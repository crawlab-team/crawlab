<template>
  <div class="spider-detail-tab-overview">
    <!-- Spider Form (handled by store) -->
    <SpiderForm/>
    <!-- ./Spider Form -->
  </div>
</template>
<script lang="ts">
import {defineComponent} from 'vue';
import SpiderForm from '@/components/spider/SpiderForm.vue';

export default defineComponent({
  name: 'SpiderDetailTabOverview',
  components: {
    SpiderForm,
  },
  setup() {
    return {};
  },
});
</script>
<style lang="scss" scoped>
.spider-detail-tab-overview {
  margin: 20px;
}
</style>
