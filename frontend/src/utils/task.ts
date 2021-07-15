import {
  TASK_MODE_ALL,
  TASK_MODE_RANDOM,
  TASK_MODE_SELECTED_NODE_TAGS,
  TASK_MODE_SELECTED_NODES,
  TASK_STATUS_PENDING,
  TASK_STATUS_RUNNING
} from '@/constants/task';

export const getPriorityLabel = (priority: number): string => {
  if (priority <= 2) {
    return `High - ${priority}`;
  } else if (priority <= 4) {
    return `Higher - ${priority}`;
  } else if (priority <= 6) {
    return `Medium - ${priority}`;
  } else if (priority <= 8) {
    return `Lower - ${priority}`;
  } else {
    return `Low - ${priority}`;
  }
};

export const isCancellable = (status: TaskStatus): boolean => {
  switch (status) {
    case TASK_STATUS_PENDING:
    case TASK_STATUS_RUNNING:
      return true;
    default:
      return false;
  }
};

export const getModeOptions = (): SelectOption[] => {
  return [
    {value: TASK_MODE_RANDOM, label: 'Random Node'},
    {value: TASK_MODE_ALL, label: 'All Nodes'},
    {value: TASK_MODE_SELECTED_NODES, label: 'Selected Nodes'},
    {value: TASK_MODE_SELECTED_NODE_TAGS, label: 'Selected Tags'},
  ];
};


export const getModeOptionsDict = (): Map<string, SelectOption> => {
  const modeOptions = getModeOptions();
  const dict = new Map<string, SelectOption>();
  modeOptions.forEach(op => dict.set(op.value, op));
  return dict;
};
