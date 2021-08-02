import {
  TASK_MODE_ALL_NODES,
  TASK_MODE_RANDOM,
  TASK_MODE_SELECTED_NODE_TAGS,
  TASK_MODE_SELECTED_NODES,
  TASK_STATUS_ABNORMAL,
  TASK_STATUS_CANCELLED,
  TASK_STATUS_ERROR,
  TASK_STATUS_FINISHED,
  TASK_STATUS_PENDING,
  TASK_STATUS_RUNNING,
} from '@/constants/task';

declare global {
  interface Task extends BaseModel {
    spider_id?: string;
    spider_name?: string;
    status?: TaskStatus;
    node_id?: string;
    node_name?: string;
    pid?: number;
    schedule_id?: string;
    schedule_name?: string;
    type?: string;
    mode?: TaskMode;
    parent_id?: string;
    cmd?: string;
    param?: string;
    error?: string;
    stat?: TaskStat;
    priority?: number;

    // view model
    spider?: Spider;
  }

  type TaskStatus =
    TASK_STATUS_PENDING
    | TASK_STATUS_RUNNING
    | TASK_STATUS_FINISHED
    | TASK_STATUS_ERROR
    | TASK_STATUS_CANCELLED
    | TASK_STATUS_ABNORMAL;

  interface TaskStat {
    create_ts?: string;
    start_ts?: string;
    end_ts?: string;
    result_count?: number;
    error_log_count?: number;
    wait_duration?: number;
    runtime_duration?: number;
    total_duration?: number;
  }

  type TaskMode =
    TASK_MODE_RANDOM
    | TASK_MODE_ALL_NODES
    | TASK_MODE_SELECTED_NODES
    | TASK_MODE_SELECTED_NODE_TAGS;
}
