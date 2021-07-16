<template>
  <el-popover
      :show-arrow="false"
      :visible="visible"
      class="table-header-dialog"
      popper-class="table-header-popper"
      trigger="manual"
  >
    <template #reference>
      <div>
        <slot name="reference"/>
      </div>
    </template>
    <div v-click-outside="onClickOutside" class="content">
      <div class="header">
        <div class="title">{{ column.label }}</div>
      </div>
      <span class="close" @click="onCancel">
        <el-icon name="close"></el-icon>
      </span>
      <div class="body">
        <div class="list">
          <div v-if="column.hasSort" class="item sort">
            <TableHeaderDialogSort :value="internalSort?.d" @change="onSortChange"/>
          </div>
          <div v-if="column.hasFilter" class="item filter">
            <TableHeaderDialogFilter
                :column="column"
                :conditions="conditions"
                :search-string="searchString"
                @change="onFilterChange"
                @enter="onFilterEnter"
            />
          </div>
        </div>
      </div>
      <div class="footer">
        <Button plain size="mini" tooltip="Cancel" type="info" @click="onCancel">Cancel</Button>
        <Button plain size="mini" tooltip="Clear" type="warning" @click="onClear">Clear</Button>
        <Button
            :disabled="isApplyDisabled"
            size="mini"
            tooltip="Apply"
            type="primary"
            @click="onApply"
        >
          Apply
        </Button>
      </div>
    </div>
  </el-popover>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, ref, watch} from 'vue';
import Button from '@/components/button/Button.vue';
import TableHeaderDialogFilter from '@/components/table/TableHeaderDialogFilter.vue';
import TableHeaderDialogSort from '@/components/table/TableHeaderDialogSort.vue';
import variables from '@/styles/variables.scss';
import {plainClone} from '@/utils/object';
import {FILTER_OP_NOT_SET} from '@/constants/filter';
import {ClickOutside} from 'element-plus/lib/directives';

export default defineComponent({
  name: 'TableHeaderFilter',
  components: {
    TableHeaderDialogSort,
    TableHeaderDialogFilter,
    Button,
  },
  directives: {
    ClickOutside,
  },
  props: {
    visible: {
      type: Boolean,
      required: false,
      default: false,
    },
    column: {
      type: Object as PropType<TableColumn>,
      required: true,
    },
    actionStatusMap: {
      type: Object,
      required: true,
    },
    sort: {
      type: Object as PropType<SortData>,
      required: false,
    },
    filter: {
      type: Object,
      required: false,
    },
  },
  emits: [
    'click',
    'cancel',
    'clear',
    'apply',
  ],
  setup(props, {emit}) {
    const defaultInternalSort = {key: props.column.key} as SortData;
    const internalSort = ref<SortData>();
    const internalFilter = ref<TableHeaderDialogFilterData>();

    const searchString = computed<string | undefined>(() => internalFilter.value?.searchString);

    const conditions = computed<FilterConditionData[]>(() => internalFilter.value?.conditions || []);

    const items = computed<string[]>(() => internalFilter.value?.items || []);

    const trueConditions = computed<FilterConditionData[]>(() => {
      return conditions.value?.filter(d => d.op !== FILTER_OP_NOT_SET);
    });

    const isEmptyFilter = computed<boolean>(() => {
      return !searchString.value && trueConditions.value.length == 0 && items.value.length === 0;
    });

    const isApplyDisabled = computed<boolean>(() => {
      for (const cond of trueConditions.value) {
        if (!cond.value) {
          return true;
        }
      }
      return false;
    });

    const cancel = () => {
      emit('cancel');
    };

    const clear = () => {
      if (!internalSort.value) internalSort.value = plainClone(defaultInternalSort);
      internalSort.value.d = undefined;
      internalFilter.value = undefined;
      emit('clear');
    };

    const apply = () => {
      if (!internalSort.value && isEmptyFilter.value) {
        clear();
        return;
      }
      const value: TableHeaderDialogValue = {
        sort: internalSort.value,
        filter: internalFilter.value,
      };
      emit('apply', value);
    };

    const onClickOutside = () => {
      const {visible} = props;
      if (!visible) return;
      cancel();
    };

    const onCancel = () => {
      cancel();
    };

    const onClear = () => {
      clear();
    };

    const onApply = () => {
      apply();
    };

    const onSortChange = (value: string) => {
      if (!internalSort.value) internalSort.value = plainClone(defaultInternalSort);
      internalSort.value.d = value;
    };

    const onFilterChange = (value: TableHeaderDialogFilterData) => {
      internalFilter.value = value;
    };

    const onFilterEnter = () => {
      apply();
    };

    watch(() => {
      const {visible} = props as TableHeaderDialogProps;
      return visible;
    }, () => {
      const {sort, filter, visible} = props as TableHeaderDialogProps;
      if (visible) {
        internalSort.value = (sort ? plainClone(sort) : plainClone(defaultInternalSort)) as SortData;
        internalFilter.value = plainClone(filter) as TableHeaderDialogFilterData;
      }
    });

    return {
      variables,
      internalSort,
      searchString,
      conditions,
      isApplyDisabled,
      onClickOutside,
      onCancel,
      onClear,
      onApply,
      onSortChange,
      onFilterChange,
      onFilterEnter,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.content {
  position: relative;
  min-width: 100%;
  flex: 1;
  display: flex;
  flex-direction: column;

  .close {
    position: absolute;
    top: 0;
    right: 0;
    cursor: pointer;
  }

  .header {
    .title {
      font-size: 16px;
      font-weight: 900;
      color: $infoMediumColor;
      padding-bottom: 10px;
      margin-bottom: 10px;
      border-bottom: 1px solid $infoBorderColor;
    }
  }

  .body {
    flex: 1;
    display: flex;
    flex-direction: column;

    .list {
      flex: 1;
      min-height: 100%;
      display: flex;
      flex-direction: column;

      .item {
        padding: 10px 0;
        border-bottom: 1px solid $infoBorderColor;

        &:first-child {
          padding-top: 0;
        }

        &:last-child {
          border-bottom: none;
        }

        &.sort {
          flex-basis: 100%;
        }

        &.filter {
          flex: 1;
          display: flex;
          flex-direction: column;
        }
      }
    }
  }

  .footer {
    height: 30px;
    text-align: right;
  }
}
</style>
<style>
.table-header-popper {
  min-width: 320px !important;
  display: flex;
  flex-direction: column;
}
</style>
