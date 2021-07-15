import {Ref, VNode} from 'vue';
import {AnyObject, Store, StoreMutations, TableColumnCtx} from 'element-plus/lib/el-table/src/table.type';
import {TABLE_ACTION_CUSTOMIZE_COLUMNS, TABLE_ACTION_EXPORT,} from '@/constants/table';

declare global {
  interface TableProps {
    data: TableData;
    columns: TableColumn[];
    selectedColumnKeys: string[];
    total: number;
    page: number;
    pageSize: number;
    rowKey: string;
    selectable: boolean;
    visibleButtons: BuiltInTableActionButtonName[];
    hideFooter: boolean;
    selectableFunction: TableSelectableFunction;
  }

  interface TableColumn<T = any> {
    key: string;
    label: string;
    icon?: string | string[];
    width?: number | string;
    minWidth?: number | string;
    index?: number;
    align?: string;
    sortable?: boolean;
    fixed?: string | boolean;
    rowKey?: string;
    buttons?: TableColumnButton[] | TableColumnButtonsFunction;
    value?: TableValueFunction<T> | any;
    disableTransfer?: boolean;
    defaultHidden?: boolean;
    hasSort?: boolean;
    hasFilter?: boolean;
    filterItems?: SelectOption[];
    allowFilterSearch?: boolean;
    allowFilterItems?: boolean;
    required?: boolean;
    className?: string;
  }

  type TableColumns<T = any> = TableColumn<T>[];

  interface TableAnyRowData {
    [key: string]: any;
  }

  type TableData<T = TableAnyRowData> = T[];

  interface TableDataWithTotal<T = TableAnyRowData> {
    data: TableData<T>;
    total: number;
  }

  interface TableColumnsMap<T = any> {
    [key: string]: TableColumn<T>;
  }

  interface TableColumnCtxMap {
    [key: string]: TableColumnCtx;
  }

  interface TableColumnButton {
    type?: string;
    size?: string;
    icon?: Icon | TableValueFunction;
    tooltip?: string | TableButtonTooltipFunction;
    isHtml?: boolean;
    disabled?: TableButtonDisabledFunction;
    onClick?: TableButtonOnClickFunction;
  }

  type TableColumnButtonsFunction<T = any> = (row?: T) => TableColumnButton[];

  type TableValueFunction<T = any> = (row: T, rowIndex?: number, column?: TableColumn<T>) => VNode;
  type TableButtonOnClickFunction<T = any> = (row: T, rowIndex?: number, column?: TableColumn<T>) => void;
  type TableButtonTooltipFunction<T = any> = (row: T, rowIndex?: number, column?: TableColumn<T>) => string;
  type TableButtonDisabledFunction<T = any> = (row: T, rowIndex?: number, column?: TableColumn<T>) => boolean;
  type TableFilterItemsFunction<T = any> = (filter?: TableHeaderDialogFilterData, column?: TableColumn<T>) => SelectOption[];
  type TableSelectableFunction<T = any> = (row: T, rowIndex?: number) => boolean;

  interface TableStore extends Store {
    mutations: TableStoreMutations;
  }

  interface TableStoreMutations extends StoreMutations {
    setColumns: (states: TableStoreStates, columns: TableColumnCtx[]) => void;
  }

  interface TableStoreStates {
    _data: Ref<AnyObject[]>;
    _columns: Ref<TableColumnCtx[]>;
  }

  interface TablePagination {
    page: number;
    size: number;
  }

  type TableActionName =
    ActionName |
    TABLE_ACTION_EXPORT |
    TABLE_ACTION_CUSTOMIZE_COLUMNS;
}
