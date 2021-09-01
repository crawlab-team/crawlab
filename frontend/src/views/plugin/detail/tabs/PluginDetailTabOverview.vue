<template>
  <div class="plugin-detail-tab-overview">
    <PluginForm/>
  </div>
</template>
<script lang="ts">
import {defineComponent} from 'vue';
import PluginForm from '@/components/plugin/PluginForm.vue';

export default defineComponent({
  name: 'PluginDetailTabOverview',
  components: {
    PluginForm,
  },
  setup() {
    return {};
  },
});
</script>
<style lang="scss" scoped>
.plugin-detail-tab-overview {
  margin: 20px;
}
</style>
