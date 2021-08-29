<template>
  App
  <hello-world/>
</template>

<script lang="ts">
import HelloWorld from './HelloWorld.vue';

export default {
  name: 'App',
  components: {
    HelloWorld,
  },
  setup() {
    return {};
  },
};
</script>
