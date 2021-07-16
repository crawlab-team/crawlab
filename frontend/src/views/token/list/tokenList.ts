import useList from '@/layouts/list';
import {useStore} from 'vuex';
import {computed} from 'vue';
import {TABLE_COLUMN_NAME_ACTIONS} from '@/constants/table';
import {ElMessage, ElMessageBox} from 'element-plus';
import useClipboard from 'vue-clipboard3';

const useTokenList = () => {
  const ns = 'token';
  const store = useStore<RootStoreState>();

  // use list
  const {
    actionFunctions,
  } = useList<Token>(ns, store);

  // action functions
  const {
    deleteByIdConfirm,
  } = actionFunctions;

  // clipboard
  const {toClipboard} = useClipboard();

  // nav actions
  const navActions = computed<ListActionGroup[]>(() => [
    {
      name: 'common',
      children: [
        {
          buttonType: 'label',
          label: 'New Token',
          tooltip: 'New Token',
          icon: ['fa', 'plus'],
          type: 'success',
          onClick: async () => {
            const res = await ElMessageBox.prompt('Please enter the name of token', 'Create');
            const name = res.value;
            const token = {
              name,
            } as Token;
            await store.dispatch(`${ns}/create`, token);
          },
        }
      ]
    }
  ]);

  // table columns
  const tableColumns = computed<TableColumns<Token>>(() => [
    {
      key: 'name',
      label: 'Name',
      icon: ['fa', 'font'],
      width: '160',
      hasFilter: true,
      allowFilterSearch: true,
    },
    {
      key: 'token',
      label: 'Token',
      icon: ['fa', 'key'],
      width: 'auto',
      value: (row: Token) => {
        if (!row._visible) {
          return (() => {
            const arr = [] as string[];
            for (let i = 0; i < 100; i++) {
              arr.push('*');
            }
            return arr.join('');
          })();
        } else {
          return row.token;
        }
      },
    },
    {
      key: TABLE_COLUMN_NAME_ACTIONS,
      label: 'Actions',
      icon: ['fa', 'tools'],
      width: '180',
      fixed: 'right',
      buttons: (row) => [
        {
          type: 'primary',
          size: 'mini',
          icon: !row._visible ? ['fa', 'eye'] : ['fa', 'eye-slash'],
          tooltip: 'View',
          onClick: async (row: Token) => {
            row._visible = !row._visible;
          },
        },
        {
          type: 'info',
          size: 'mini',
          icon: ['far', 'clipboard'],
          tooltip: 'Copy',
          onClick: async (row: Token) => {
            if (!row.token) return;
            await toClipboard(row.token);
            await ElMessage.success('Copied token to clipboard');
          },
        },
        {
          type: 'danger',
          size: 'mini',
          icon: ['fa', 'trash-alt'],
          tooltip: 'Edit',
          onClick: deleteByIdConfirm,
        },
      ],
      disableTransfer: true,
    },
  ]);

  // options
  const opts = {
    navActions,
    tableColumns,
  } as UseListOptions<Token>;

  return {
    ...useList<Token>(ns, store, opts),
  };
};

export default useTokenList;
