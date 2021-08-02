import {ElMessageBox} from 'element-plus';
import {useStore} from 'vuex';
import {useRouter} from 'vue-router';
import {ACTION_CANCEL, ACTION_CLONE, ACTION_DELETE, ACTION_EDIT, ACTION_RUN, ACTION_VIEW,} from '@/constants/action';
import {TABLE_COLUMN_NAME_ACTIONS} from '@/constants/table';

export const getDefaultTableDataWithTotal = (): TableDataWithTotal => {
  return {
    data: [],
    total: 0,
  };
};

export const getTableWidth = (): number | undefined => {
  const elTable = document.querySelector('.list-layout .table');
  if (!elTable) return;
  const style = getComputedStyle(elTable);
  const widthStr = style.width.replace('px', '');
  const width = Number(widthStr);
  if (isNaN(width)) return;
  return width;
};

export const getColumnWidth = (column: TableColumn): number | undefined => {
  let width: number;
  if (typeof column.width === 'string') {
    width = Number(column.width.replace('px', ''));
    if (isNaN(width)) return;
    return width;
  }
  {
    return column.width;
  }
};

export const getActionColumn = (endpoint: string, ns: ListStoreNamespace, actionNames: TableActionName[]): TableColumn => {
  const store = useStore();
  const router = useRouter();

  const column = {
    key: TABLE_COLUMN_NAME_ACTIONS,
    label: 'Actions',
    fixed: 'right',
    width: '200',
    buttons: [],
  } as TableColumn;

  const buttons = typeof column.buttons === 'function' ? column.buttons() : column.buttons as TableColumnButton[];

  actionNames.forEach(name => {
    if (!buttons) return;
    switch (name) {
      case ACTION_VIEW:
        buttons.push({
          type: 'primary',
          icon: ['fa', 'search'],
          tooltip: 'View',
          onClick: (row: BaseModel) => {
            router.push(`${endpoint}/${row._id}`);
          },
        });
        break;
      case ACTION_EDIT:
        buttons.push({
          type: 'warning',
          icon: ['fa', 'edit'],
          tooltip: 'Edit',
          onClick: (row: BaseModel) => {
            store.commit(`${ns}/setForm`, row);
            store.commit(`${ns}/showDialog`, 'edit');
          },
        },);
        break;
      case ACTION_CLONE:
        buttons.push({
          type: 'info',
          size: 'mini',
          icon: ['fa', 'clone'],
          tooltip: 'Clone',
          onClick: (row: BaseModel) => {
            // TODO: implement
            console.log('clone', row);
          }
        });
        break;
      case ACTION_DELETE:
        buttons.push({
          type: 'danger',
          size: 'mini',
          icon: ['fa', 'trash-alt'],
          tooltip: 'Delete',
          onClick: async (row: BaseModel) => {
            const res = await ElMessageBox.confirm('Are you sure to delete?', 'Delete');
            if (res) {
              await store.dispatch(`${ns}/deleteById`, row._id as string);
            }
            await store.dispatch(`${ns}/getList`);
          },
        });
        break;
      case ACTION_RUN:
        buttons.push({
          type: 'success',
          size: 'mini',
          icon: ['fa', 'play'],
          tooltip: 'Run',
          onClick: async (row: BaseModel) => {
            store.commit(`${ns}/setForm`, row);
            store.commit(`${ns}/showDialog`, 'run');
          },
        });
        break;
      case ACTION_CANCEL:
        buttons.push({
          type: 'info',
          size: 'mini',
          icon: ['fa', 'pause'],
          tooltip: 'Cancel',
          onClick: async (row: BaseModel) => {
            // TODO: implement
            console.log('cancel', row);
          },
        });
        break;
    }
  });

  return column;
};
