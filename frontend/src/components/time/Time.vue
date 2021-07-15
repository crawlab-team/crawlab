<template>
  <div class="time">
    {{ label }}
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import TimeAgo from 'javascript-time-ago';
import en from 'javascript-time-ago/locale/en';
// import zh from 'javascript-time-ago/locale/zh';

TimeAgo.addDefaultLocale(en);
// TimeAgo.addDefaultLocale(zh);

export default defineComponent({
  name: 'Time',
  props: {
    time: {
      type: [Date, String] as PropType<Date | string>,
      required: false,
      default: () => new Date(),
    },
    ago: {
      type: Boolean,
      required: false,
      default: true,
    }
  },
  setup(props: TimeProps, {emit}) {
    const timeAgo = new TimeAgo();

    const label = computed<string | undefined>(() => {
      const {time} = props;
      if (!time) return;
      return timeAgo.format(new Date(time));
    });

    return {
      label,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
