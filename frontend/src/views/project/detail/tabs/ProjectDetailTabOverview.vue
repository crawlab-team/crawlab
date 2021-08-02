<template>
  <div class="project-detail-tab-overview">
    <ProjectForm/>
  </div>
</template>
<script lang="ts">
import {defineComponent} from 'vue';
import ProjectForm from '@/components/project/ProjectForm.vue';

export default defineComponent({
  name: 'ProjectDetailTabOverview',
  components: {
    ProjectForm,
  },
  setup() {
    return {};
  },
});
</script>
<style lang="scss" scoped>
.project-detail-tab-overview {
  margin: 20px;
}
</style>
