import {
  TABLE_ACTION_ADD,
  TABLE_ACTION_CUSTOMIZE_COLUMNS,
  TABLE_ACTION_DELETE,
  TABLE_ACTION_EDIT,
  TABLE_ACTION_EXPORT,
} from '@/constants/table';

declare global {
  interface TableActionsProps {
    selection: TableData;
    visibleButtons: BuiltInTableActionButtonName[];
  }

  type BuiltInTableActionButtonName =
    TABLE_ACTION_ADD
    | TABLE_ACTION_EDIT
    | TABLE_ACTION_DELETE
    | TABLE_ACTION_EXPORT
    | TABLE_ACTION_CUSTOMIZE_COLUMNS;
}
