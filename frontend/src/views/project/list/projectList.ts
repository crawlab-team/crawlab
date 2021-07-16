import {computed, h} from 'vue';
import {TABLE_COLUMN_NAME_ACTIONS} from '@/constants/table';
import {useStore} from 'vuex';
import {ElMessageBox} from 'element-plus';
import useList from '@/layouts/list';
import useProjectService from '@/services/project/projectService';
import NavLink from '@/components/nav/NavLink.vue';
import {useRouter} from 'vue-router';
import TagList from '@/components/tag/TagList.vue';
import {TAB_NAME_SPIDERS} from '@/constants/tab';

const useProjectList = () => {
  // router
  const router = useRouter();

  // store
  const ns = 'project';
  const store = useStore<RootStoreState>();
  const {commit} = store;

  // services
  const {
    deleteById,
    getList,
  } = useProjectService(store);

  // nav actions
  const navActions = computed<ListActionGroup[]>(() => [
    {
      name: 'common',
      children: [
        {
          buttonType: 'label',
          label: 'New Project',
          tooltip: 'New Project',
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
  const tableColumns = computed<TableColumns<Project>>(() => [
    {
      key: 'name',
      label: 'Name',
      icon: ['fa', 'font'],
      width: '150',
      value: (row: Project) => h(NavLink, {
        path: `/projects/${row._id}`,
        label: row.name,
      }),
      hasSort: true,
      hasFilter: true,
      allowFilterSearch: true,
    },
    {
      key: 'spiders',
      label: 'Spiders',
      icon: ['fa', 'spider'],
      value: (row: Project) => h(NavLink, {
        path: `/projects/${row._id}/${TAB_NAME_SPIDERS}`,
        label: row.spiders,
      }),
      width: '120',
    },
    {
      key: 'tags',
      label: 'Tags',
      icon: ['fa', 'hashtag'],
      value: ({tags}: Project) => {
        return h(TagList, {tags});
      },
      width: '200',
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
            router.push(`/projects/${row._id}`);
          },
        },
        {
          type: 'warning',
          icon: ['fa', 'edit'],
          tooltip: 'Edit',
          onClick: (row) => {
            store.commit(`${ns}/setForm`, row);
            store.commit(`${ns}/showDialog`, 'edit');
          },
        },
        {
          type: 'info',
          size: 'mini',
          icon: ['fa', 'clone'],
          tooltip: 'Clone',
          onClick: (row) => {
            console.log('clone', row);
          }
        },
        {
          type: 'danger',
          size: 'mini',
          icon: ['fa', 'trash-alt'],
          tooltip: 'Delete',
          onClick: async (row: Project) => {
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
  const opts = {
    navActions,
    tableColumns,
  } as UseListOptions<Project>;

  return {
    ...useList<Project>(ns, store, opts),
  };
};

export default useProjectList;
