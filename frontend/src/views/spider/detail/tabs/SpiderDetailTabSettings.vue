<template>
  Settings
</template>
<script lang="ts">
import {computed, defineComponent} from 'vue';
import {useRoute} from 'vue-router';

export default defineComponent({
  name: 'SpiderDetailTabSettings',
  setup() {
    const route = useRoute();

    const id = computed(() => route.params.id);

    return {
      id,
    };
  },
});
</script>
<style scoped lang="scss">
</style>
