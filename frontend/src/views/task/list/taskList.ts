import useList from '@/layouts/list';
import {useStore} from 'vuex';
import {computed, h} from 'vue';
import NavLink from '@/components/nav/NavLink.vue';
import TaskStatus from '@/components/task/TaskStatus.vue';
import {TABLE_COLUMN_NAME_ACTIONS} from '@/constants/table';
import {useRouter} from 'vue-router';
import {ElMessage, ElMessageBox} from 'element-plus';
import useRequest from '@/services/request';
import TaskPriority from '@/components/task/TaskPriority.vue';
import NodeType from '@/components/node/NodeType.vue';
import Time from '@/components/time/Time.vue';
import Duration from '@/components/time/Duration.vue';
import {setupListComponent} from '@/utils/list';
import {isCancellable} from '@/utils/task';
import TaskResults from '@/components/task/TaskResults.vue';
import useNode from '@/components/node/node';
import useSpider from '@/components/spider/spider';
import {
  TASK_STATUS_CANCELLED,
  TASK_STATUS_ERROR,
  TASK_STATUS_FINISHED,
  TASK_STATUS_PENDING,
  TASK_STATUS_RUNNING
} from '@/constants/task';
import useTask from '@/components/task/task';

const {
  post,
} = useRequest();

const useTaskList = () => {
  const ns = 'task';
  const store = useStore<RootStoreState>();
  const {commit} = store;

  // router
  const router = useRouter();

  // use list
  const {
    actionFunctions,
  } = useList<Task>(ns, store);

  // action functions
  const {
    deleteByIdConfirm,
  } = actionFunctions;

  // all node dict
  const allNodeDict = computed<Map<string, CNode>>(() => store.getters['node/allDict']);

  // all spider dict
  const allSpiderDict = computed<Map<string, Spider>>(() => store.getters['spider/allDict']);

  // nav actions
  const navActions = computed<ListActionGroup[]>(() => [
    {
      name: 'common',
      children: [
        {
          buttonType: 'label',
          label: 'New Task',
          tooltip: 'New Task',
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
    allListSelectOptions: allNodeListSelectOptions,
  } = useNode(store);

  const {
    allListSelectOptions: allSpiderListSelectOptions,
  } = useSpider(store);

  const {
    priorityOptions,
  } = useTask(store);

  // table columns
  const tableColumns = computed<TableColumns<Task>>(() => [
    {
      key: 'node_id',
      label: 'Node',
      icon: ['fa', 'server'],
      width: '160',
      value: (row: Task) => {
        if (!row.node_id) return;
        const node = allNodeDict.value.get(row.node_id);
        if (!node) return;
        return h(NodeType, {
          isMaster: node?.is_master,
          label: node?.name,
          onClick: () => {
            router.push(`/nodes/${node?._id}`);
          }
        } as NodeTypeProps);
      },
      hasFilter: true,
      allowFilterSearch: true,
      allowFilterItems: true,
      filterItems: allNodeListSelectOptions.value,
    },
    {
      key: 'spider_id',
      label: 'Spider',
      icon: ['fa', 'spider'],
      width: '160',
      value: (row: Task) => {
        if (!row.spider_id) return;
        const spider = allSpiderDict.value.get(row.spider_id);
        return h(NavLink, {
          label: spider?.name,
          path: `/spiders/${spider?._id}`,
        });
      },
      hasFilter: true,
      allowFilterSearch: true,
      allowFilterItems: true,
      filterItems: allSpiderListSelectOptions.value,
    },
    {
      key: 'priority',
      label: 'Priority',
      icon: ['fa', 'sort-numeric-down'],
      width: '120',
      value: (row: Task) => {
        return h(TaskPriority, {priority: row.priority} as TaskPriorityProps);
      },
      hasSort: true,
      hasFilter: true,
      allowFilterItems: true,
      filterItems: priorityOptions,
    },
    {
      key: 'status',
      label: 'Status',
      icon: ['fa', 'check-square'],
      width: '120',
      value: (row: Task) => {
        return h(TaskStatus, {status: row.status, error: row.error} as TaskStatusProps);
      },
      hasFilter: true,
      allowFilterItems: true,
      filterItems: [
        {label: 'Pending', value: TASK_STATUS_PENDING},
        {label: 'Running', value: TASK_STATUS_RUNNING},
        {label: 'Finished', value: TASK_STATUS_FINISHED},
        {label: 'Error', value: TASK_STATUS_ERROR},
        {label: 'Cancelled', value: TASK_STATUS_CANCELLED},
      ],
    },
    {
      key: 'stat.create_ts',
      label: 'Created At',
      icon: ['fa', 'clock'],
      width: '120',
      value: (row: Task) => {
        if (!row.stat?.create_ts || row.stat?.create_ts.startsWith('000')) return;
        return h(Time, {time: row.stat?.create_ts as string} as TimeProps);
      },
      defaultHidden: true,
    },
    {
      key: 'stat.start_ts',
      label: 'Started At',
      icon: ['fa', 'clock'],
      width: '120',
      value: (row: Task) => {
        if (!row.stat?.start_ts || row.stat?.start_ts.startsWith('000')) return;
        return h(Time, {time: row.stat?.start_ts as string} as TimeProps);
      },
    },
    {
      key: 'stat.end_ts',
      label: 'Finished At',
      icon: ['fa', 'clock'],
      width: '120',
      value: (row: Task) => {
        if (!row.stat?.end_ts || row.stat?.end_ts.startsWith('000')) return;
        return h(Time, {time: row.stat?.end_ts as string} as TimeProps);
      },
    },
    {
      key: 'stat.wait_duration',
      label: 'Wait Duration',
      icon: ['fa', 'stopwatch'],
      width: '160',
      value: (row: Task) => {
        if (!row.stat?.wait_duration) return;
        return h(Duration, {duration: row.stat?.wait_duration as number} as DurationProps);
      },
      defaultHidden: true,
    },
    {
      key: 'stat.runtime_duration',
      label: 'Runtime Duration',
      icon: ['fa', 'stopwatch'],
      width: '160',
      value: (row: Task) => {
        if (!row.stat?.runtime_duration) return;
        return h(Duration, {duration: row.stat?.runtime_duration as number} as DurationProps);
      },
      defaultHidden: true,
    },
    {
      key: 'stat.total_duration',
      label: 'Total Duration',
      icon: ['fa', 'stopwatch'],
      width: '160',
      value: (row: Task) => {
        if (!row.stat?.total_duration) return;
        return h(Duration, {duration: row.stat?.total_duration as number} as DurationProps);
      },
    },
    {
      key: 'stat.result_count',
      label: 'Results',
      icon: ['fa', 'table'],
      width: '150',
      value: (row: Task) => {
        if (row.stat?.result_count === undefined) return;
        return h(TaskResults, {results: row.stat.result_count, status: row.status} as TaskResultsProps);
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
          icon: ['fa', 'search'],
          tooltip: 'View',
          onClick: (row) => {
            router.push(`/tasks/${row._id}`);
          }
        },
        {
          type: 'warning',
          size: 'mini',
          icon: ['fa', 'redo'],
          tooltip: 'Restart',
          onClick: async (row) => {
            await ElMessageBox.confirm('Are you sure to restart?', 'Restart', {type: 'warning'});
            await post(`/tasks/${row._id}/restart`);
            await ElMessage.success('Restarted successfully');
            await store.dispatch(`task/getList`);
          }
        },
        isCancellable(row.status) ?
          {
            type: 'info',
            size: 'mini',
            icon: ['fa', 'pause'],
            tooltip: 'Cancel',
            onClick: async (row: Task) => {
              await ElMessageBox.confirm('Are you sure to cancel?', 'Cancel', {type: 'warning'});
              await ElMessage.info('Attempt to cancel');
              await post(`/tasks/${row._id}/cancel`);
              await store.dispatch(`task/getList`);
            },
          }
          :
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

  // options
  const opts = {
    navActions,
    tableColumns,
  } as UseListOptions<Task>;

  // init
  setupListComponent(ns, store, ['node', 'project', 'spider']);

  return {
    ...useList<Task>(ns, store, opts),
  };
};

export default useTaskList;
