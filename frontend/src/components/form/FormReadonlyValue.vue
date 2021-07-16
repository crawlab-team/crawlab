<template>
  <div class="form-readonly-value">
    {{ value }}
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue';

export default defineComponent({
  name: 'FormReadonlyValue',
  props: {
    value: {
      type: String,
      required: true,
    },
  },
  setup(props: FormReadonlyValueProps, {emit}) {
    return {};
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables";

.form-readonly-value {
  font-size: 14px;
  color: $infoColor;
}
</style>
