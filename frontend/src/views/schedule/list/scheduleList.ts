import {computed, h} from 'vue';
import {TABLE_COLUMN_NAME_ACTIONS} from '@/constants/table';
import {useStore} from 'vuex';
import {ElMessage, ElMessageBox} from 'element-plus';
import useList from '@/layouts/list';
import useScheduleService from '@/services/schedule/scheduleService';
import NavLink from '@/components/nav/NavLink.vue';
import {useRouter} from 'vue-router';
import {setupListComponent} from '@/utils/list';
import TaskMode from '@/components/task/TaskMode.vue';
import ScheduleCron from '@/components/schedule/ScheduleCron.vue';
import Switch from '@/components/switch/Switch.vue';
import useSpider from '@/components/spider/spider';
import useTask from '@/components/task/task';

const useScheduleList = () => {
  // router
  const router = useRouter();

  // store
  const ns = 'schedule';
  const store = useStore<RootStoreState>();
  const {commit} = store;

  // services
  const {
    deleteById,
    getList,
  } = useScheduleService(store);

  // all node dict
  const allNodeDict = computed<Map<string, CNode>>(() => store.getters['node/allDict']);

  // all spider dict
  const allSpiderDict = computed<Map<string, Spider>>(() => store.getters['spider/allDict']);

  const {
    allListSelectOptions: allSpiderListSelectOptions,
  } = useSpider(store);

  const {
    modeOptions,
  } = useTask(store);

  // nav actions
  const navActions = computed<ListActionGroup[]>(() => [
    {
      name: 'common',
      children: [
        {
          buttonType: 'label',
          label: 'New Schedule',
          tooltip: 'New Schedule',
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
  const tableColumns = computed<TableColumns<Schedule>>(() => [
    {
      key: 'name',
      label: 'Name',
      icon: ['fa', 'font'],
      width: '150',
      value: (row: Schedule) => h(NavLink, {
        path: `/schedules/${row._id}`,
        label: row.name,
      }),
      hasSort: true,
      hasFilter: true,
      allowFilterSearch: true,
    },
    {
      key: 'spider_id',
      label: 'Spider',
      icon: ['fa', 'spider'],
      width: '160',
      value: (row: Schedule) => {
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
      key: 'mode',
      label: 'Mode',
      icon: ['fa', 'cog'],
      width: '160',
      value: (row: Schedule) => {
        return h(TaskMode, {mode: row.mode} as TaskModeProps);
      },
      hasFilter: true,
      allowFilterItems: true,
      filterItems: modeOptions,
    },
    {
      key: 'cron',
      label: 'Cron Expression',
      icon: ['fa', 'clock'],
      width: '160',
      value: (row: Schedule) => {
        return h(ScheduleCron, {cron: row.cron} as ScheduleCronProps);
      },
      hasFilter: true,
      allowFilterSearch: true,
    },
    {
      key: 'enabled',
      label: 'Enabled',
      icon: ['fa', 'toggle-on'],
      width: '120',
      value: (row: Schedule) => {
        return h(Switch, {
          modelValue: row.enabled,
          'onUpdate:modelValue': async (value: boolean) => {
            if (value) {
              await store.dispatch(`${ns}/enable`, row._id);
              ElMessage.success('Enabled successfully');
            } else {
              await store.dispatch(`${ns}/disable`, row._id);
              ElMessage.success('Disabled successfully');
            }
            await store.dispatch(`${ns}/getList`);
          },
        } as SwitchProps);
      },
      hasFilter: true,
      allowFilterItems: true,
      filterItems: [
        {label: 'Enabled', value: true},
        {label: 'Disabled', value: false},
      ],
    },
    {
      key: 'entry_id',
      label: 'Entry ID',
      icon: ['fa', 'hash'],
      width: '120',
      defaultHidden: true,
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
            router.push(`/schedules/${row._id}`);
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
          onClick: async (row: Schedule) => {
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
  } as UseListOptions<Schedule>;

  // init
  setupListComponent(ns, store, ['node', 'spider']);

  return {
    ...useList<Schedule>(ns, store, opts),
  };
};

export default useScheduleList;
