<template>
  <div class="empty">
    <ImgEmpty/>
    <div class="description">
      {{ description }}
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import ImgEmpty from '@/components/empty/ImgEmpty.vue';

export default defineComponent({
  name: 'Empty',
  components: {ImgEmpty},
  props: {
    description: {
      type: String,
      required: false,
      default: 'No Data Available'
    }
  },
  setup() {
    return {};
  },
});
</script>

<style lang="scss" scoped>
.empty {
  max-height: 240px;
  display: flex;
  flex-direction: column;

  .description {
    margin-top: 10px;
    display: flex;
    justify-content: center;
  }
}
</style>
