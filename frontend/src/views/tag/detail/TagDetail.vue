<template>
  <DetailLayout store-namespace="tag">
  </DetailLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import DetailLayout from '@/layouts/DetailLayout.vue';

export default defineComponent({
  name: 'TagDetail',
  components: {DetailLayout},
  setup(props, {emit}) {
    return {};
  },
});
</script>

<style lang="scss" scoped>

</style>
