package models

type ModelMap struct {
	Artifact          Artifact
	Tag               Tag
	Node              Node
	Project           Project
	Spider            Spider
	Task              Task
	Job               Job
	Schedule          Schedule
	User              User
	Setting           Setting
	Token             Token
	Variable          Variable
	TaskQueueItem     TaskQueueItem
	TaskStat          TaskStat
	SpiderStat        SpiderStat
	DataSource        DataSource
	DataCollection    DataCollection
	Result            Result
	Password          Password
	ExtraValue        ExtraValue
	Git               Git
	Role              Role
	UserRole          UserRole
	Permission        Permission
	RolePermission    RolePermission
	Environment       Environment
	DependencySetting DependencySetting
}

type ModelListMap struct {
	Artifacts          ArtifactList
	Tags               TagList
	Nodes              NodeList
	Projects           ProjectList
	Spiders            SpiderList
	Tasks              TaskList
	Jobs               JobList
	Schedules          ScheduleList
	Users              UserList
	Settings           SettingList
	Tokens             TokenList
	Variables          VariableList
	TaskQueueItems     TaskQueueItemList
	TaskStats          TaskStatList
	SpiderStats        SpiderStatList
	DataSources        DataSourceList
	DataCollections    DataCollectionList
	Results            ResultList
	Passwords          PasswordList
	ExtraValues        ExtraValueList
	Gits               GitList
	Roles              RoleList
	UserRoles          UserRoleList
	PermissionList     PermissionList
	RolePermissionList RolePermissionList
	Environments       EnvironmentList
	DependencySettings DependencySettingList
}

func NewModelMap() (m *ModelMap) {
	return &ModelMap{}
}

func NewModelListMap() (m *ModelListMap) {
	return &ModelListMap{
		Artifacts:          ArtifactList{},
		Tags:               TagList{},
		Nodes:              NodeList{},
		Projects:           ProjectList{},
		Spiders:            SpiderList{},
		Tasks:              TaskList{},
		Jobs:               JobList{},
		Schedules:          ScheduleList{},
		Users:              UserList{},
		Settings:           SettingList{},
		Tokens:             TokenList{},
		Variables:          VariableList{},
		TaskQueueItems:     TaskQueueItemList{},
		TaskStats:          TaskStatList{},
		SpiderStats:        SpiderStatList{},
		DataSources:        DataSourceList{},
		DataCollections:    DataCollectionList{},
		Results:            ResultList{},
		Passwords:          PasswordList{},
		ExtraValues:        ExtraValueList{},
		Gits:               GitList{},
		Roles:              RoleList{},
		PermissionList:     PermissionList{},
		RolePermissionList: RolePermissionList{},
		Environments:       EnvironmentList{},
	}
}
