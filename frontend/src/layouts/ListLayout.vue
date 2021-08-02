<template>
  <div class="list-layout">
    <div class="content">
      <!-- Nav Actions -->
      <NavActions v-if="!noActions" ref="navActions" class="nav-actions">
        <NavActionGroup
            v-for="(grp, i) in navActions"
            :key="i"
        >
          <NavActionItem
              v-for="(btn, j) in grp.children"
              :key="j"
          >
            <NavActionButton
                :button-type="btn.buttonType"
                :disabled="btn.disabled"
                :icon="btn.icon"
                :label="btn.label"
                :size="btn.size"
                :tooltip="btn.tooltip"
                :type="btn.type"
                @click="btn.onClick"
            />
          </NavActionItem>
        </NavActionGroup>
      </NavActions>
      <!-- ./Nav Actions -->

      <!-- Table -->
      <Table
          ref="tableRef"
          :key="tableColumnsHash"
          :columns="tableColumns"
          :data="tableData"
          :total="tableTotal"
          :page="tablePagination.page"
          :page-size="tablePagination.size"
          selectable
          :selectable-function="selectableFunction"
          @selection-change="onSelect"
          @delete="onDelete"
          @edit="onEdit"
          @pagination-change="onPaginationChange"
          @header-change="onHeaderChange"
      >
        <template #actions-prefix>
          <NavActionButton
              v-for="(btn, $index) in tableActionsPrefix"
              :key="$index"
              :button-type="btn.buttonType"
              :disabled="getNavActionButtonDisabled(btn)"
              :icon="btn.icon"
              :label="btn.label"
              :size="btn.size"
              :tooltip="btn.tooltip"
              :type="btn.type"
              @click="btn.onClick"
          />
        </template>
        <template #actions-suffix>
          <NavActionButton
              v-for="(btn, $index) in tableActionsSuffix"
              :key="$index"
              :button-type="btn.buttonType"
              :disabled="btn.disabled"
              :icon="btn.icon"
              :label="btn.label"
              :size="btn.size"
              :tooltip="btn.tooltip"
              :type="btn.type"
              @click="btn.onClick"
          />
        </template>
      </Table>
      <!-- ./Table -->
    </div>

    <slot name="extra"/>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onBeforeMount, PropType, provide, ref, SetupContext, toRefs, watch} from 'vue';
import NavActionGroup from '@/components/nav/NavActionGroup.vue';
import NavActionItem from '@/components/nav/NavActionItem.vue';
import Table from '@/components/table/Table.vue';
import NavActionButton from '@/components/nav/NavActionButton.vue';
import NavActions from '@/components/nav/NavActions.vue';
import {emptyObjectFunc} from '@/utils/func';
import {getMd5} from '@/utils/hash';

export default defineComponent({
  name: 'ListLayout',
  components: {
    NavActions,
    NavActionGroup,
    NavActionItem,
    NavActionButton,
    Table,
  },
  props: {
    navActions: {
      type: Array as PropType<ListActionGroup[]>,
      default: () => {
        return [];
      }
    },
    tableColumns: {
      type: Array as PropType<TableColumns>,
      default: () => {
        return [];
      }
    },
    tableData: {
      type: Array as PropType<TableData>,
      default: () => {
        return [];
      }
    },
    tableTotal: {
      type: Number,
      default: 0,
    },
    tablePagination: {
      type: Object as PropType<TablePagination>,
      default: () => {
        return {
          page: 1,
          size: 10,
        };
      }
    },
    tableActionsPrefix: {
      type: Array as PropType<ListActionButton[]>,
      default: () => {
        return [];
      }
    },
    tableActionsSuffix: {
      type: Array as PropType<ListActionButton[]>,
      default: () => {
        return [];
      }
    },
    actionFunctions: {
      type: Object as PropType<ListLayoutActionFunctions>,
      default: emptyObjectFunc,
    },
    noActions: {
      type: Boolean,
      default: false,
    },
    selectableFunction: {
      type: Function as PropType<TableSelectableFunction>,
      default: () => true,
    },
  },
  emits: [
    'select',
    'edit',
    'delete',
  ],
  setup(props: ListLayoutProps, {emit}: SetupContext) {
    const {
      actionFunctions,
    } = toRefs(props);

    const {
      setPagination,
      getList,
      onHeaderChange,
    } = actionFunctions.value;

    const tableRef = ref();

    const computedTableRef = computed<typeof Table>(() => tableRef.value);

    const onSelect = (value: TableData) => {
      emit('select', value);
    };

    const onEdit = (value: TableData) => {
      emit('edit', value);
    };

    const onDelete = (value: TableData) => {
      emit('delete', value);
    };

    const onPaginationChange = (value: TablePagination) => {
      setPagination(value);
    };

    watch(() => props.tablePagination, getList);

    onBeforeMount(() => {
      getList();
    });

    provide<ListLayoutActionFunctions>('action-functions', actionFunctions.value);

    const getNavActionButtonDisabled = (btn: ListActionButton) => {
      if (typeof btn.disabled === 'boolean') {
        return btn.disabled;
      } else if (typeof btn.disabled === 'function') {
        return btn.disabled(computedTableRef.value);
      } else {
        return false;
      }
    };

    const tableColumnsHash = computed<string>(() => {
      const {tableColumns} = props;
      return getMd5(JSON.stringify(tableColumns));
    });

    return {
      tableRef,
      tableColumnsHash,
      onSelect,
      onPaginationChange,
      onHeaderChange,
      onEdit,
      onDelete,
      getNavActionButtonDisabled,
    };
  },
});
</script>

<style lang="scss" scoped>
@import "../styles/variables.scss";

.list-layout {
  .nav-actions {
    background-color: $containerWhiteBg;
    border-bottom: none;
  }

  .content {
    background-color: $containerWhiteBg;
  }
}
</style>
<style scoped>
.list-layout >>> .tag {
  margin-right: 10px;
}

</style>
