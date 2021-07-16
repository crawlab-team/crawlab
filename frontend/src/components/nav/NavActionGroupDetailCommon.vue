<template>
  <NavActionGroup>
    <NavActionButton
        :disabled="disabled"
        :icon="['fa', 'undo']"
        button-type="label"
        label="Back"
        type="primary"
        @click="() => $emit('back')"
    />
    <NavActionButton
        :disabled="disabled"
        :icon="['fa', 'save']"
        button-type="label"
        label="Save"
        type="success"
        @click="() => $emit('save')"
    />
  </NavActionGroup>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import NavActionGroup from '@/components/nav/NavActionGroup.vue';
import NavActionButton from '@/components/nav/NavActionButton.vue';

export default defineComponent({
  name: 'NavActionGroupDetailCommon',
  components: {NavActionButton, NavActionGroup},
  props: {
    disabled: {
      type: Boolean,
      default: false
    }
  },
  emits: [
    'back',
  ],
  setup() {
    return {};
  },
});
</script>

<style lang="scss" scoped>

</style>
