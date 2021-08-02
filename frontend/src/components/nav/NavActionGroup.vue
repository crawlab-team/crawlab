<template>
  <div class="nav-action-group">
    <div class="border"/>
    <slot></slot>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue';

export default defineComponent({
  name: 'NavActionGroup',
  setup() {
    return {};
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.nav-action-group {
  height: auto;
  display: flex;
  align-items: center;
  box-sizing: border-box;

  & + .nav-action-group {
    //padding-left: 10px;
    margin-left: 10px;

    .border {
      margin-left: -10px;
      margin-right: 10px;
      border-left: 1px solid $navActionsGroupBorderColor;
      height: calc(100% - 20px);
    }
  }
}
</style>

<style scoped>
.nav-action-group >>> .nav-action-item:last-child {
  margin-right: 0;
}
</style>
