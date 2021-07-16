import {SetupContext} from 'vue';

declare global {
  interface DraggableListProps {
    items: DraggableItemData[];
    itemKey: string;
  }

  interface DraggableListContext {
    ctx: SetupContext<any>;
    props: DraggableListProps;
  }

  interface DraggableListInternalItems {
    draggingItem?: DraggableItemData;
    targetItem?: DraggableItemData;
  }
}
