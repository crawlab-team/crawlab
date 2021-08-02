<template>
  <el-tag :type="type" class="user-role" size="mini">
    <font-awesome-icon :icon="icon" class="icon"/>
    <span>{{ computedLabel }}</span>
  </el-tag>
</template>

<script lang="ts">
import {computed, defineComponent, PropType} from 'vue';
import {ROLE_ADMIN} from '@/constants/user';

export default defineComponent({
  name: 'UserRole',
  props: {
    role: {
      type: String as PropType<UserRole>,
    },
    label: {
      type: String,
    },
  },
  setup(props: UserRoleProps, {emit}) {
    const type = computed<string>(() => {
      const {role} = props;
      return role === ROLE_ADMIN ? 'primary' : 'warning';
    });

    const computedLabel = computed<string>(() => {
      const {role, label} = props;
      if (label) return label;
      return role === ROLE_ADMIN ? 'Admin' : 'Normal';
    });

    const icon = computed<string[]>(() => {
      const {role} = props;
      return role === ROLE_ADMIN ? ['fa', 'star'] : ['fa', 'user'];
    });

    return {
      type,
      computedLabel,
      icon,
    };
  },
});
</script>

<style lang="scss" scoped>
.user-role {
  .icon {
    margin-right: 5px;
  }
}
</style>
