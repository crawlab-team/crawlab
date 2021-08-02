<template>
  <div class="table-header-dialog-sort">
    <div class="title">
      <span>Sort</span>
      <el-tooltip v-if="value" content="Clear sort">
        <span class="icon" @click="onClear">
          <el-icon name="circle-close"/>
        </span>
      </el-tooltip>
    </div>
    <el-radio-group :model-value="value" size="mini" type="primary" @change="onChange">
      <el-radio-button :label="ASCENDING" class="sort-btn">
        <font-awesome-icon :icon="['fa', 'sort-amount-up']"/>
        Ascending
      </el-radio-button>
      <el-radio-button :label="DESCENDING" class="sort-btn">
        <font-awesome-icon :icon="['fa', 'sort-amount-down-alt']"/>
        Descending
      </el-radio-button>
    </el-radio-group>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {ASCENDING, DESCENDING, UNSORTED} from '@/constants/sort';

export default defineComponent({
  name: 'TableHeaderDialogSort',
  props: {
    value: {
      type: String,
      required: false,
    },
  },
  emits: [
    'change'
  ],
  setup(props, {emit}) {
    const onChange = (value: SortDirection) => {
      if (value === UNSORTED) {
        emit('change', undefined);
        return;
      }
      emit('change', value);
    };

    const onClear = () => {
      emit('change');
    };

    return {
      onChange,
      onClear,
      UNSORTED,
      ASCENDING,
      DESCENDING,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.table-header-dialog-sort {
  .title {
    font-size: 14px;
    font-weight: 900;
    margin-bottom: 10px;
    color: $infoMediumColor;
    display: flex;
    align-items: center;

    .icon {
      cursor: pointer;
      margin-left: 5px;
    }
  }

  .el-radio-group {
    width: 100%;
    display: flex;

    .sort-btn.el-radio-button {
      &:not(.unsorted) {
        flex: 1;
      }

      &.unsorted {
        flex-basis: 20px;
      }
    }
  }
}
</style>
<style scoped>
.table-header-dialog-sort >>> .el-radio-group .el-radio-button .el-radio-button__inner {
  width: 100%;
}
</style>
