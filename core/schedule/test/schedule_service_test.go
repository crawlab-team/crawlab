package test

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestScheduleService_Enable_Disable(t *testing.T) {
	var err error
	T.Setup(t)

	time.Sleep(1 * time.Second)
	err = T.scheduleSvc.Enable(T.TestSchedule)
	require.Nil(t, err)
	time.Sleep(1 * time.Second)

	require.True(t, T.TestSchedule.GetEnabled())
	require.Greater(t, int(T.TestSchedule.GetEntryId()), -1)
	e := T.scheduleSvc.GetCron().Entry(T.TestSchedule.GetEntryId())
	require.Equal(t, T.TestSchedule.GetEntryId(), e.ID)
	time.Sleep(1 * time.Second)

	err = T.scheduleSvc.Disable(T.TestSchedule)
	require.False(t, T.TestSchedule.GetEnabled())
	require.Equal(t, 0, len(T.scheduleSvc.GetCron().Entries()))
}

func TestScheduleService_Run(t *testing.T) {
	var err error
	T.Setup(t)

	time.Sleep(1 * time.Second)
	err = T.scheduleSvc.Enable(T.TestSchedule)
	require.Nil(t, err)
	time.Sleep(1 * time.Minute)

	tasks, err := T.modelSvc.GetTaskList(nil, nil)
	require.Nil(t, err)
	require.Greater(t, len(tasks), 0)
	for _, task := range tasks {
		require.False(t, task.ScheduleId.IsZero())
	}
}
