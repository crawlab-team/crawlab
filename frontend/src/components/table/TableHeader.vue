<template>
  <div class="table-header">
    <span :class="[column.required ? 'required' : '']" class="label">
      <span v-if="column.icon" class="label-icon">
        <Icon :icon="column.icon"/>
      </span>
      {{ column.label }}
    </span>

    <TableHeaderDialog
        v-if="hasDialog"
        :action-status-map="actionStatusMap"
        :column="column"
        :visible="dialogVisible"
        :filter="filterData"
        :sort="sortData"
        @apply="onDialogApply"
        @clear="onDialogClear"
        @cancel="onDialogCancel"
    >
      <template #reference>
        <div class="actions">
          <TableHeaderAction
              v-for="{key, tooltip, isHtml, icon, onClick} in actions"
              :key="key + JSON.stringify(icon)"
              :icon="icon"
              :status="actionStatusMap[key]"
              :tooltip="tooltip"
              :is-html="isHtml"
              @click="onClick"
          />
        </div>
      </template>
    </TableHeaderDialog>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, PropType, reactive, ref} from 'vue';
import TableHeaderDialog from '@/components/table/TableHeaderDialog.vue';
import TableHeaderAction from '@/components/table/TableHeaderAction.vue';
import {conditionTypesMap} from '@/components/filter/FilterCondition.vue';
import {ASCENDING, DESCENDING} from '@/constants/sort';
import variables from '@/styles/variables.scss';
import {FILTER_OP_NOT_SET} from '@/constants/filter';
import Icon from '@/components/icon/Icon.vue';

export default defineComponent({
  name: 'TableHeader',
  components: {
    Icon,
    TableHeaderAction,
    TableHeaderDialog,
  },
  props: {
    column: {
      type: Object as PropType<TableColumn>,
      required: true,
    },
    index: {
      type: Number,
      required: false,
    },
  },
  emits: [
    'change',
  ],
  setup(props, {emit}) {
    const dialogVisible = ref<boolean>(false);

    const actionStatusMap = reactive<TableHeaderActionStatusMap>({
      filter: {active: false, focused: false},
      sort: {active: false, focused: false},
    });

    const sortData = ref<SortData>();
    const filterData = ref<TableHeaderDialogFilterData>();

    const filterItemsMap = computed<Map<any, string | undefined>>(() => {
      const map = new Map<any, string | undefined>();
      const {column} = props;
      const {filterItems} = column;
      if (!filterItems) return map;
      filterItems.forEach(d => {
        const {label, value} = d;
        map.set(value, label);
      });
      return map;
    });

    const actions = computed<TableColumnButton[]>(() => {
      const {column} = props;

      // sort icon and tooltip
      let sortIcon = ['fa', 'sort-amount-down-alt'];
      let sortTooltip = 'Sort';
      if (sortData.value?.d === ASCENDING) {
        sortIcon = ['fa', 'sort-amount-up'];
        sortTooltip = 'Sorted Ascending';
      } else if (sortData.value?.d === DESCENDING) {
        sortIcon = ['fa', 'sort-amount-down-alt'];
        sortTooltip = 'Sorted Descending';
      }

      // filter tooltip
      let filterTooltip = 'Filter';
      let filterIsHtml = false;
      if (filterData.value) {
        const {searchString, conditions, items} = filterData.value;

        // search string
        if (searchString) {
          filterTooltip += `<br><span style="color: ${variables.primaryColor}">Search:</span> <span style="color: ${variables.warningColor};">"${searchString}"</span>`;
          filterIsHtml = true;
        }

        // filter conditions
        if (conditions && conditions.length > 0) {
          filterTooltip += '<br>' + conditions.filter(d => d.op !== FILTER_OP_NOT_SET)
              .map(d => `<span style="color: ${variables.primaryColor};margin-right: 5px">${conditionTypesMap[d.op || '']}:</span> <span style="color: ${variables.warningColor};">"${d.value}"</span>`)
              .join('<br>');
          filterIsHtml = true;
        }

        // filter items
        if (items && items.length > 0) {
          const itemsStr = items.map(value => filterItemsMap.value.get(value)).join(', ');
          filterTooltip += `<br><span style="color: ${variables.primaryColor};margin-right: 5px">Include:</span><span style="color: ${variables.warningColor}">` + itemsStr + '</span>';
          filterIsHtml = true;
        }
      }

      // tooltip items
      const items = [];
      if (column.hasSort) {
        items.push({
          key: 'sort',
          tooltip: sortTooltip,
          icon: sortIcon,
          onClick: () => {
            dialogVisible.value = true;
            actionStatusMap.sort.focused = true;
          }
        });
      }
      if (column.hasFilter) {
        items.push({
          key: 'filter',
          tooltip: filterTooltip,
          isHtml: filterIsHtml,
          icon: ['fa', 'filter'],
          onClick: () => {
            dialogVisible.value = true;
            actionStatusMap.filter.focused = true;
          }
        });
      }

      return items;
    });

    const hideDialog = () => {
      dialogVisible.value = false;
      actionStatusMap.filter.focused = false;
      actionStatusMap.sort.focused = false;
    };

    const clearDialog = () => {
      const {column} = props as TableHeaderProps;

      // set status
      actionStatusMap.filter.active = false;
      actionStatusMap.sort.active = false;

      // set data
      sortData.value = undefined;
      filterData.value = undefined;

      // hide
      hideDialog();

      // emit
      emit('change', column, undefined, undefined);
    };

    const onDialogCancel = () => {
      hideDialog();
    };

    const onDialogClear = () => {
      clearDialog();
    };

    const onDialogApply = (value: TableHeaderDialogValue) => {
      const {column} = props as TableHeaderProps;
      const {sort, filter} = value;

      // set status
      if (sort) actionStatusMap.sort.active = true;
      if (filter) actionStatusMap.filter.active = true;

      // set data
      sortData.value = sort;
      filterData.value = filter;

      // if no data set, clear
      if (!sortData.value && !filterData.value) {
        clearDialog();
        return;
      }

      // hide
      hideDialog();

      // emit
      emit('change', column, sortData.value, filterData.value);
    };

    const hasDialog = computed<boolean>(() => {
      const {
        hasSort,
        hasFilter,
      } = props.column;

      return !!hasSort || !!hasFilter;
    });

    return {
      dialogVisible,
      actionStatusMap,
      actions,
      sortData,
      filterData,
      onDialogCancel,
      onDialogClear,
      onDialogApply,
      hasDialog,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../../styles/variables.scss";

.table-header {
  display: flex;
  position: relative;

  .label {
    display: flex;
    align-items: center;

    &.required:before {
      content: "*";
      color: $red;
      margin-right: 4px;
    }

    .label-icon {
      color: $infoMediumLightColor;
      font-size: 8px;
      margin-right: 5px;
    }
  }

  .actions {
    position: absolute;
    right: 0;
    height: 100%;
    align-items: center;
    cursor: pointer;
  }
}
</style>
<style scoped>
.table-header:hover .actions >>> .action {
  display: inline;
}

.table-header .actions >>> .action {
  display: none;
}
</style>
