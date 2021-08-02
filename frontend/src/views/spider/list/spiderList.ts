import {useRouter} from 'vue-router';
import {useStore} from 'vuex';
import {computed, h} from 'vue';
import TaskStatus from '@/components/task/TaskStatus.vue';
import {TABLE_COLUMN_NAME_ACTIONS} from '@/constants/table';
import useList from '@/layouts/list';
import NavLink from '@/components/nav/NavLink.vue';
import Time from '@/components/time/Time.vue';
import SpiderStat from '@/components/spider/SpiderStat.vue';
import {setupListComponent} from '@/utils/list';
import useProject from '@/components/project/project';

const useSpiderList = () => {
  // router
  const router = useRouter();

  // store
  const ns = 'spider';
  const store = useStore<RootStoreState>();
  const {commit} = store;

  // use list
  const {
    actionFunctions,
  } = useList<Task>(ns, store);

  // action functions
  const {
    deleteByIdConfirm,
  } = actionFunctions;

  // nav actions
  const navActions = computed<ListActionGroup[]>(() => [
    {
      name: 'common',
      children: [
        {
          buttonType: 'label',
          label: 'New Spider',
          tooltip: 'New Spider',
          icon: ['fa', 'plus'],
          type: 'success',
          onClick: () => {
            commit(`${ns}/showDialog`, 'create');
          }
        }
      ]
    }
  ]);

  const {
    allListSelectOptions: allProjectListSelectOptions,
  } = useProject(store);
  // const allProjectList = computed<Project[]>(() => store.state.project.allList);

  // all project dict
  const allProjectDict = computed<Map<string, Project>>(() => store.getters['project/allDict']);

  // table columns
  const tableColumns = computed<TableColumns<Spider>>(() => [
    {
      key: 'name',
      label: 'Name',
      icon: ['fa', 'font'],
      width: '160',
      align: 'left',
      value: (row: Spider) => h(NavLink, {
        path: `/spiders/${row._id}`,
        label: row.name,
      }),
      hasSort: true,
      hasFilter: true,
      allowFilterSearch: true,
    },
    // {
    //   key: 'spider_type',
    //   label: 'Spider Type',
    //   icon: ['fa', 'list'],
    //   width: '120',
    //   filterItems: [
    //     {label: 'Customized', value: 'customized'},
    //     {label: 'Configurable', value: 'configurable'},
    //   ],
    //   value: (row: Spider) => {
    //     return h(SpiderType, {type: row.spider_type});
    //   },
    //   hasFilter: true,
    // },
    {
      key: 'project_id',
      label: 'Project',
      icon: ['fa', 'project-diagram'],
      width: '120',
      value: (row: Spider) => {
        if (!row.project_id) return;
        const p = allProjectDict.value.get(row.project_id);
        return h(NavLink, {
          label: p?.name,
          path: `/projects/${row.project_id}`,
        });
      },
      hasFilter: true,
      allowFilterSearch: true,
      allowFilterItems: true,
      filterItems: allProjectListSelectOptions.value,
    },
    // {
    //   key: 'is_long_task',
    //   label: 'Is Long Task',
    //   width: '80',
    // },
    // {
    //   key: 'latest_tasks',
    //   label: 'Latest Tasks',
    //   icon: ['fa', 'project-diagram'],
    //   width: '180',
    //   defaultHidden: true,
    // },
    {
      key: 'last_status',
      label: 'Last Status',
      icon: ['fa', 'heartbeat'],
      width: '120',
      value: (row: Spider) => {
        const status = row.stat?.last_task?.status;
        if (!status) return;
        return h(TaskStatus, {status} as TaskStatusProps);
      }
    },
    {
      key: 'last_run_ts',
      label: 'Last Run At',
      icon: ['fa', 'clock'],
      width: '160',
      value: (row: Spider) => {
        const time = row.stat?.last_task?.stat?.start_ts;
        if (!time) return;
        return h(Time, {time} as TaskStatusProps);
      },
    },
    {
      key: 'stats',
      label: 'Stats',
      icon: ['fa', 'chart-pie'],
      width: '240',
      hasFilter: true,
      value: (row: Spider) => {
        const stat = row.stat;
        if (!stat || !stat.tasks) return;
        return h(SpiderStat, {stat} as SpiderStatProps);
      }
    },
    {
      key: 'create_ts',
      label: 'Created At',
      icon: ['far', 'calendar-plus'],
      width: '160',
      defaultHidden: true,
    },
    {
      key: 'update_ts',
      label: 'Updated At',
      icon: ['far', 'calendar-check'],
      width: '160',
      defaultHidden: true,
    },
    // {
    //   key: 'create_username',
    //   label: 'Created By',
    //   icon: ['fa', 'user'],
    //   width: '100',
    //   hasFilter: true,
    //   defaultHidden: true,
    // },
    {
      key: 'description',
      label: 'Description',
      icon: ['fa', 'comment-alt'],
      width: 'auto',
    },
    {
      key: TABLE_COLUMN_NAME_ACTIONS,
      label: 'Actions',
      icon: ['fa', 'tools'],
      width: '180',
      fixed: 'right',
      buttons: [
        {
          type: 'success',
          size: 'mini',
          icon: ['fa', 'play'],
          tooltip: 'Run',
          onClick: (row) => {
            store.commit(`${ns}/setForm`, row);
            store.commit(`${ns}/showDialog`, 'run');
          },
        },
        {
          type: 'primary',
          size: 'mini',
          icon: ['fa', 'search'],
          tooltip: 'View',
          onClick: (row) => {
            router.push(`/spiders/${row._id}`);
          }
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
          onClick: deleteByIdConfirm,
        },
      ],
      disableTransfer: true,
    },
  ]);

  // table actions prefix
  const tableActionsPrefix = computed<ListActionButton[]>(() => {
    return [
      // {
      //   buttonType: 'fa-icon',
      //   tooltip: 'Run',
      //   size: 'mini',
      //   icon: ['fa', 'play'],
      //   type: 'success',
      //   disabled: (table: typeof Table) => {
      //     return !table?.internalSelection?.length;
      //   },
      // }
    ];
  });

  // const onClickCreate = () => {
  //   commit(`${ns}/showDialog`, 'create');
  // };
  //
  // const onClickEdit = () => {
  //   commit(`${ns}/showDialog`, 'edit');
  // };
  //
  // const onClickClone = () => {
  //   commit(`${ns}/showDialog`, 'clone');
  // };

  // const onClickRun = () => {
  //   commit(`${ns}/showDialog`, 'run');
  // };

  // options
  const opts = {
    navActions,
    tableColumns,
  } as UseListOptions<Spider>;

  // init
  setupListComponent(ns, store, ['node', 'project', 'dataCollection']);

  return {
    ...useList<Spider>(ns, store, opts),
    tableActionsPrefix,
  };
};

export default useSpiderList;
