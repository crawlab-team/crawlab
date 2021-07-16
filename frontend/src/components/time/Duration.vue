<template>
  <div class="duration">
    {{ label }}
  </div>
</template>

<script lang="ts">
import {computed, defineComponent} from 'vue';
import humanizeDuration from 'humanize-duration';

export default defineComponent({
  name: 'Duration',
  props: {
    duration: {
      type: Number,
      required: false,
    },
  },
  setup(props: DurationProps, {emit}) {
    const label = computed<string>(() => {
      const {duration} = props;
      const d = Math.ceil((duration as number) / 1000) * 1000;
      return humanizeDuration(d, {
        spacer: ' ',
      });
    });

    return {
      label,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
