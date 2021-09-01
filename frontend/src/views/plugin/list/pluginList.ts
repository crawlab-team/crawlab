import useList from '@/layouts/list';
import {useStore} from 'vuex';
import {getDefaultUseListOptions, setupListComponent} from '@/utils/list';
import {computed, h} from 'vue';
import {TABLE_COLUMN_NAME_ACTIONS} from '@/constants/table';
import {ElMessageBox} from 'element-plus';
import usePluginService from '@/services/plugin/pluginService';
import NavLink from '@/components/nav/NavLink.vue';
import {useRouter} from 'vue-router';

type Plugin = CPlugin;

const usePluginList = () => {
  // router
  const router = useRouter();

  // store
  const ns = 'plugin';
  const store = useStore<RootStoreState>();
  const {commit} = store;

  // services
  const {
    getList,
    deleteById,
  } = usePluginService(store);

  // nav actions
  const navActions = computed<ListActionGroup[]>(() => [
    {
      name: 'common',
      children: [
        {
          buttonType: 'label',
          label: 'New Plugin',
          tooltip: 'New Plugin',
          icon: ['fa', 'plus'],
          type: 'success',
          onClick: () => {
            commit(`${ns}/showDialog`, 'create');
          }
        }
      ]
    }
  ]);

  // table columns
  const tableColumns = computed<TableColumns<Plugin>>(() => [
    {
      key: 'name', // name
      label: 'Name',
      icon: ['fa', 'font'],
      width: '150',
      value: (row: Plugin) => h(NavLink, {
        path: `/plugins/${row._id}`,
        label: row.name,
      }),
      hasSort: true,
      hasFilter: true,
      allowFilterSearch: true,
    },
    {
      key: 'description',
      label: 'Description',
      icon: ['fa', 'comment-alt'],
      width: 'auto',
      hasFilter: true,
      allowFilterSearch: true,
    },
    {
      key: TABLE_COLUMN_NAME_ACTIONS,
      label: 'Actions',
      fixed: 'right',
      width: '200',
      buttons: [
        {
          type: 'primary',
          icon: ['fa', 'search'],
          tooltip: 'View',
          onClick: (row) => {
            router.push(`/plugins/${row._id}`);
          },
        },
        // {
        //   type: 'info',
        //   size: 'mini',
        //   icon: ['fa', 'clone'],
        //   tooltip: 'Clone',
        //   onClick: (row) => {
        //     console.log('clone', row);
        //   }
        // },
        {
          type: 'danger',
          size: 'mini',
          icon: ['fa', 'trash-alt'],
          tooltip: 'Delete',
          disabled: (row: Plugin) => !!row.active,
          onClick: async (row: Plugin) => {
            const res = await ElMessageBox.confirm('Are you sure to delete?', 'Delete');
            if (res) {
              await deleteById(row._id as string);
            }
            await getList();
          },
        },
      ],
      disableTransfer: true,
    }
  ]);

  // options
  const opts = getDefaultUseListOptions<Plugin>(navActions, tableColumns);

  // init
  setupListComponent(ns, store, []);

  return {
    ...useList<Plugin>(ns, store, opts)
  };
};

export default usePluginList;
