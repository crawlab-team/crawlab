<template>
  <div class="schedule-detail-tab-overview">
    <ScheduleForm/>
  </div>
</template>
<script lang="ts">
import {defineComponent} from 'vue';
import ScheduleForm from '@/components/schedule/ScheduleForm.vue';

export default defineComponent({
  name: 'ProjectDetailTabOverview',
  components: {
    ScheduleForm,
  },
  setup() {
    return {};
  },
});
</script>
<style lang="scss" scoped>
.schedule-detail-tab-overview {
  margin: 20px;
}
</style>
