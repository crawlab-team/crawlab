<template>
  <DetailLayout store-namespace="schedule">
  </DetailLayout>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import DetailLayout from '@/layouts/DetailLayout.vue';
import useScheduleDetail from '@/views/schedule/detail/scheduleDetail';

export default defineComponent({
  name: 'ScheduleDetail',
  components: {DetailLayout},
  setup() {
    return {
      ...useScheduleDetail(),
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
