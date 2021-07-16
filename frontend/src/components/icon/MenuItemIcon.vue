<template>
  <template v-if="!item || !item.icon">
    <i :style="{'font-size': fontSize}" class="menu-item-icon fa fa-circle-o"></i>
  </template>
  <template v-else-if="Array.isArray(item.icon)">
    <font-awesome-icon
        :icon="item.icon"
        :style="{'font-size': fontSize}"
        class="menu-item-icon"
    />
  </template>
  <template v-else>
    <i :class="item.icon" :style="{'font-size': fontSize}" class="menu-item-icon"></i>
  </template>
</template>
<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import useIcon from '@/components/icon/icon';

export default defineComponent({
  name: 'MenuItemIcon',
  props: {
    item: {
      type: Object as PropType<MenuItem>,
    },
    size: {
      type: String as PropType<IconSize>,
      default: 'mini',
    }
  },
  setup(props: MenuItemIconProps) {
    const {
      getFontSize,
    } = useIcon();

    const fontSize = computed(() => {
      const {size} = props as MenuItemIconProps;
      return getFontSize(size);
    });

    return {
      fontSize,
    };
  },
});
</script>
<style lang="scss" scoped>
.menu-item-icon {
  width: 20px;
}
</style>
