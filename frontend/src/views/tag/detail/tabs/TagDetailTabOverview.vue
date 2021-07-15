<template>
  <div class="tag-detail-tab-overview">
    <TagForm/>
  </div>
</template>
<script lang="ts">
import {defineComponent} from 'vue';
import TagForm from '@/components/tag/TagForm.vue';

export default defineComponent({
  name: 'TagDetailTabOverview',
  components: {
    TagForm,
  },
  setup() {
    return {};
  },
});
</script>
<style lang="scss" scoped>
.tag-detail-tab-overview {
  margin: 20px;
}
</style>
