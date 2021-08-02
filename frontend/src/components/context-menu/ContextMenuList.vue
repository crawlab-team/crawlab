<template>
  <ul class="context-menu-list">
    <li
        v-for="(item, $index) in items"
        :key="$index"
        class="context-menu-item"
        @click="onClick(item)"
    >
      <span class="prefix">
        <template v-if="item.icon">
          <font-awesome-icon v-if="Array.isArray(item.icon)" :icon="item.icon"/>
          <atom-material-icon v-else-if="typeof item.icon === 'string'" :is-dir="false" :name="item.icon"/>
        </template>
      </span>
      <span class="title">
        {{ item.title }}
      </span>
    </li>
  </ul>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import AtomMaterialIcon from '@/components/icon/AtomMaterialIcon.vue';

export default defineComponent({
  name: 'ContextMenuList',
  components: {AtomMaterialIcon},
  props: {
    items: {
      type: [Array, String],
      default: () => {
        return [];
      },
    },
  },
  setup(props, {emit}) {
    const onClick = (item: ContextMenuItem) => {
      if (!item.action) return;
      item.action();
      emit('hide');
    };

    return {
      onClick,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.context-menu-list {
  list-style: none;
  margin: 0;
  padding: 0;
  min-width: auto;

  .context-menu-item {
    height: $contextMenuItemHeight;
    max-width: $contextMenuItemMaxWidth;
    display: flex;
    align-items: center;
    margin: 0;
    padding: 10px;
    cursor: pointer;

    &:hover {
      background-color: $primaryPlainColor;
    }

    .title {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .prefix {
      width: 24px;
      display: flex;
      align-items: center;
    }
  }
}
</style>
