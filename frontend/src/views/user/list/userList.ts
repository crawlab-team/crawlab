import {computed, h} from 'vue';
import {TABLE_COLUMN_NAME_ACTIONS} from '@/constants/table';
import {useStore} from 'vuex';
import {ElMessageBox} from 'element-plus';
import useList from '@/layouts/list';
import useUserService from '@/services/user/userService';
import NavLink from '@/components/nav/NavLink.vue';
import {useRouter} from 'vue-router';
import UserRole from '@/components/user/UserRole.vue';
import {ROLE_ADMIN, ROLE_NORMAL, USERNAME_ADMIN} from '@/constants/user';

const useUserList = () => {
  // router
  const router = useRouter();

  // store
  const ns = 'user';
  const store = useStore<RootStoreState>();
  const {commit} = store;

  // services
  const {
    deleteById,
    getList,
  } = useUserService(store);

  // nav actions
  const navActions = computed<ListActionGroup[]>(() => [
    {
      name: 'common',
      children: [
        {
          buttonType: 'label',
          label: 'New User',
          tooltip: 'New User',
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
  const tableColumns = computed<TableColumns<User>>(() => [
    {
      key: 'username',
      label: 'Username',
      icon: ['fa', 'font'],
      width: '180',
      value: (row: User) => h(NavLink, {
        path: `/users/${row._id}`,
        label: row.username,
      }),
      hasSort: true,
      hasFilter: true,
      allowFilterSearch: true,
    },
    {
      key: 'email',
      label: 'Email',
      icon: ['fa', 'at'],
      width: '180',
      hasSort: true,
      hasFilter: true,
      allowFilterSearch: true,
    },
    {
      key: 'role',
      label: 'Role',
      icon: ['fa', 'font'],
      width: '150',
      value: (row: User) => h(UserRole, {role: row.role} as UserRoleProps),
      hasFilter: true,
      allowFilterItems: true,
      filterItems: [
        {label: 'Admin', value: ROLE_ADMIN},
        {label: 'Normal', value: ROLE_NORMAL},
      ],
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
            router.push(`/users/${row._id}`);
          },
        },
        {
          type: 'danger',
          size: 'mini',
          icon: ['fa', 'trash-alt'],
          tooltip: (row: User) => row.username === USERNAME_ADMIN ? 'Admin user is non-deletable' : 'Delete',
          disabled: (row: User) => row.username === USERNAME_ADMIN,
          onClick: async (row: User) => {
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

  const selectableFunction = (row: User) => {
    return row.username !== USERNAME_ADMIN;
  };

  // options
  const opts = {
    navActions,
    tableColumns,
  } as UseListOptions<User>;

  return {
    ...useList<User>(ns, store, opts),
    selectableFunction,
  };
};

export default useUserList;
