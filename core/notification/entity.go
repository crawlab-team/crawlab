package notification

import "github.com/crawlab-team/crawlab/core/models/models/v2"

type VariableData struct {
	Task     *models.TaskV2              `json:"task"`
	TaskStat *models.TaskStatV2          `json:"task_stat"`
	Spider   *models.SpiderV2            `json:"spider"`
	Node     *models.NodeV2              `json:"node"`
	Schedule *models.ScheduleV2          `json:"schedule"`
	Alert    *models.NotificationAlertV2 `json:"alert"`
	Metric   *models.MetricV2            `json:"metric"`
}
