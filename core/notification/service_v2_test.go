package notification

import (
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTemplateVariables_WithValidTemplate_ReturnsVariables(t *testing.T) {
	svc := ServiceV2{}
	template := "Dear ${user:name}, your task ${task:id} is ${task:status}."
	expected := []entity.NotificationVariable{
		{Category: "user", Name: "name"},
		{Category: "task", Name: "id"},
		{Category: "task", Name: "status"},
	}
	setting := models.NotificationSettingV2{Template: template}

	variables := svc.parseTemplateVariables(&setting)

	assert.Equal(t, expected, variables)
}

func TestParseTemplateVariables_WithRepeatedVariables_ReturnsUniqueVariables(t *testing.T) {
	svc := ServiceV2{}
	template := "Dear ${user:name}, your task ${task:id} is ${task:status}. Again, ${user:name} and ${task:id}."
	expected := []entity.NotificationVariable{
		{Category: "user", Name: "name"},
		{Category: "task", Name: "id"},
		{Category: "task", Name: "status"},
	}
	setting := models.NotificationSettingV2{Template: template}

	variables := svc.parseTemplateVariables(&setting)

	assert.Equal(t, expected, variables)
}
