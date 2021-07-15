<template>
  <DetailLayout store-namespace="node">
  </DetailLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import DetailLayout from '@/layouts/DetailLayout.vue';

export default defineComponent({
  name: 'NodeDetail',
  components: {DetailLayout},
  setup(props, {emit}) {
    return {};
  },
});
</script>

<style lang="scss" scoped>

</style>
