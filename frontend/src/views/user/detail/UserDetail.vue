<template>
  <DetailLayout nav-item-name-key="username" store-namespace="user">
  </DetailLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import DetailLayout from '@/layouts/DetailLayout.vue';
import useUserDetail from '@/views/user/detail/userDetail';

export default defineComponent({
  name: 'UserDetail',
  components: {DetailLayout},
  setup() {
    return {
      ...useUserDetail(),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
