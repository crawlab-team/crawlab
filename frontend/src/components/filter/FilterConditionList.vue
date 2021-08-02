<template>
  <ul class="filter-condition-list">
    <li
        v-for="(cond, $index) in conditions"
        :key="$index"
        class="filter-condition-item"
    >
      <FilterCondition :condition="cond" @change="onChange($index, $event)" @delete="onDelete($index)"/>
    </li>
  </ul>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import FilterCondition, {getDefaultFilterCondition} from '@/components/filter/FilterCondition.vue';

export default defineComponent({
  name: 'FilterConditionList',
  components: {
    FilterCondition,
  },
  props: {
    conditions: {
      type: Array,
      required: false,
      default: () => {
        return [];
      }
    }
  },
  emits: [
    'change',
  ],
  setup(props, {emit}) {
    const onChange = (index: number, condition: FilterConditionData) => {
      const {conditions} = props as FilterConditionListProps;
      conditions[index] = condition;
      emit('change', conditions);
    };

    const onDelete = (index: number) => {
      const {conditions} = props as FilterConditionListProps;
      conditions.splice(index, 1);
      if (conditions.length === 0) {
        conditions.push(getDefaultFilterCondition());
      }
      emit('change', conditions);
    };

    return {
      onChange,
      onDelete,
    };
  },
});
</script>

<style lang="scss" scoped>
.filter-condition-list {
  list-style: none;
  padding: 0;
  margin: 0;

  .filter-condition-item {
    margin-bottom: 10px;
  }
}
</style>
