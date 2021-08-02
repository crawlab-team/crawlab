import {useRoute} from 'vue-router';
import {computed} from 'vue';
import {TASK_MODE_RANDOM} from '@/constants/task';
import {Store} from 'vuex';
import useForm from '@/components/form/form';
import useSpiderService from '@/services/spider/spiderService';
import {getDefaultFormComponentData} from '@/utils/form';
import {
  FORM_FIELD_TYPE_INPUT,
  FORM_FIELD_TYPE_INPUT_TEXTAREA,
  FORM_FIELD_TYPE_INPUT_WITH_BUTTON,
  FORM_FIELD_TYPE_SELECT
} from '@/constants/form';
import useProject from '@/components/project/project';
import useRequest from '@/services/request';
import {FILTER_OP_CONTAINS} from '@/constants/filter';
import {getModeOptions} from '@/utils/task';

const {
  getList,
} = useRequest();

// get new spider
export const getNewSpider = (): Spider => {
  return {
    mode: TASK_MODE_RANDOM,
  };
};

// form component data
const formComponentData = getDefaultFormComponentData<Spider>(getNewSpider);

const useSpider = (store: Store<RootStoreState>) => {
  // options for default mode
  const modeOptions = getModeOptions();

  // use project
  const {
    allProjectSelectOptions,
  } = useProject(store);

  // batch form fields
  const batchFormFields = computed<FormTableField[]>(() => [
    {
      prop: 'name',
      label: 'Name',
      width: '150',
      placeholder: 'Spider Name',
      fieldType: FORM_FIELD_TYPE_INPUT,
      required: true,
    },
    {
      prop: 'cmd',
      label: 'Execute Command',
      width: '200',
      placeholder: 'Execute Command',
      fieldType: FORM_FIELD_TYPE_INPUT_WITH_BUTTON,
      required: true,
    },
    {
      prop: 'param',
      label: 'Param',
      width: '200',
      placeholder: 'Param',
      fieldType: FORM_FIELD_TYPE_INPUT_WITH_BUTTON,
    },
    {
      prop: 'mode',
      label: 'Default Run Mode',
      width: '200',
      fieldType: FORM_FIELD_TYPE_SELECT,
      options: modeOptions,
      required: true,
    },
    {
      prop: 'project_id',
      label: 'Project',
      width: '200',
      fieldType: FORM_FIELD_TYPE_SELECT,
      options: allProjectSelectOptions.value,
    },
    {
      prop: 'description',
      label: 'Description',
      width: '200',
      fieldType: FORM_FIELD_TYPE_INPUT_TEXTAREA,
    },
  ]);

  // route
  const route = useRoute();

  // spider id
  const id = computed(() => route.params.id);

  // fetch data collections
  const fetchDataCollection = async (query: string) => {
    const conditions = [{
      key: 'name',
      op: FILTER_OP_CONTAINS,
      value: query,
    }] as FilterConditionData[];
    const res = await getList(`/data/collections`, {conditions});
    return res.data;
  };

  // fetch data collection suggestions
  const fetchDataCollectionSuggestions = (query: string, cb: Function) => {
    fetchDataCollection(query)
      .then(data => {
        cb(data?.map((d: DataCollection) => {
          return {
            _id: d._id,
            value: d.name,
          };
        }));
      });
  };

  return {
    ...useForm('spider', store, useSpiderService(store), formComponentData),
    batchFormFields,
    id,
    modeOptions,
    fetchDataCollection,
    fetchDataCollectionSuggestions,
  };
};

export default useSpider;
