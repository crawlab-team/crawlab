import {useStore} from 'vuex';
import {computed, h} from 'vue';
import useList from '@/layouts/list';
import NavLink from '@/components/nav/NavLink.vue';
import ColorPicker from '@/components/color/ColorPicker.vue';
import {getActionColumn} from '@/utils/table';
import {ACTION_DELETE, ACTION_VIEW} from '@/constants/action';
import Tag from '@/components/tag/Tag.vue';

const useTagList = () => {
  // store
  const ns = 'tag';
  const store = useStore<RootStoreState>();
  const {commit} = store;

  // nav actions
  const navActions = computed<ListActionGroup[]>(() => [
    {
      name: 'common',
      children: [
        {
          buttonType: 'label',
          label: 'New Tag',
          tooltip: 'New Tag',
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
  const tableColumns = computed<TableColumns<Tag>>(() => [
    {
      key: 'name',
      label: 'Name',
      icon: ['fa', 'font'],
      width: '160',
      align: 'left',
      value: (row: Tag) => h(NavLink, {
        path: `/tags/${row._id}`,
        label: row.name,
      }),
    },
    {
      key: 'color',
      label: 'Color',
      icon: ['fa', 'palette'],
      width: '120',
      value: ({color}: Tag) => {
        return h(ColorPicker, {
          modelValue: color,
          disabled: true,
        });
      }
    },
    {
      key: 'col',
      label: 'Model',
      icon: ['fa', 'table'],
      width: '120',
      // value: ({color}: Tag) => {
      //   return h(Tag, {
      //     modelValue: color,
      //     disabled: true,
      //   });
      // }
    },
    {
      key: 'description',
      label: 'Description',
      icon: ['fa', 'comment-alt'],
      width: 'auto',
    },
    getActionColumn('/tags', ns, [ACTION_VIEW, ACTION_DELETE]),
  ]);

  // options
  const opts = {
    navActions,
    tableColumns,
  } as UseListOptions<Tag>;

  return {
    ...useList<Tag>(ns, store, opts),
  };
};

export default useTagList;
