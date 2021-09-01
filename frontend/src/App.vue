<template>
  <router-view />
</template>
<script lang="ts">
import {defineComponent, onBeforeMount} from 'vue';
import {initPlugins} from '@/utils/plugin';
import {useStore} from 'vuex';

export default defineComponent({
  name: 'App',
  setup() {
    onBeforeMount(() => {
      const store = useStore();

      initPlugins(store);
    });

    return {};
  },
});
</script>
