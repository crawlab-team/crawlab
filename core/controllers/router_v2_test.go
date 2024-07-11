package controllers_test

import (
	"github.com/crawlab-team/crawlab/core/controllers"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouterGroups(t *testing.T) {
	router := gin.Default()
	groups := controllers.NewRouterGroups(router)

	assertions := []struct {
		group *gin.RouterGroup
		name  string
	}{
		{groups.AuthGroup, "AuthGroup"},
		{groups.AnonymousGroup, "AnonymousGroup"},
	}

	for _, a := range assertions {
		assert.NotNil(t, a.group, a.name+" should not be nil")
	}
}

func TestRegisterController_Routes(t *testing.T) {
	router := gin.Default()
	groups := controllers.NewRouterGroups(router)
	ctr := controllers.NewControllerV2[models.TestModelV2]()
	basePath := "/testmodels"

	controllers.RegisterController(groups.AuthGroup, basePath, ctr)

	// Check if all routes are registered
	routes := router.Routes()

	var methodPaths []string
	for _, route := range routes {
		methodPaths = append(methodPaths, route.Method+" - "+route.Path)
	}

	expectedRoutes := []gin.RouteInfo{
		{Method: "GET", Path: basePath},
		{Method: "GET", Path: basePath + "/:id"},
		{Method: "POST", Path: basePath},
		{Method: "PUT", Path: basePath + "/:id"},
		{Method: "PATCH", Path: basePath},
		{Method: "DELETE", Path: basePath + "/:id"},
		{Method: "DELETE", Path: basePath},
	}

	assert.Equal(t, len(expectedRoutes), len(routes))
	for _, route := range expectedRoutes {
		assert.Contains(t, methodPaths, route.Method+" - "+route.Path)
	}
}

func TestInitRoutes_ProjectsRoute(t *testing.T) {
	router := gin.Default()

	controllers.InitRoutes(router)

	// Check if the projects route is registered
	routes := router.Routes()

	var methodPaths []string
	for _, route := range routes {
		methodPaths = append(methodPaths, route.Method+" - "+route.Path)
	}

	expectedRoutes := []gin.RouteInfo{
		{Method: "GET", Path: "/projects"},
		{Method: "GET", Path: "/projects/:id"},
		{Method: "POST", Path: "/projects"},
		{Method: "PUT", Path: "/projects/:id"},
		{Method: "PATCH", Path: "/projects"},
		{Method: "DELETE", Path: "/projects/:id"},
		{Method: "DELETE", Path: "/projects"},
	}

	for _, route := range expectedRoutes {
		assert.Contains(t, methodPaths, route.Method+" - "+route.Path)
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}
