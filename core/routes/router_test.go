package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestInitRoutes(t *testing.T) {
	app := gin.New()
	err := InitRoutes(app)
	require.Nil(t, err)

	srv := &http.Server{
		Handler: app,
		Addr:    "localhost:8000",
	}
	go func() {
		err = srv.ListenAndServe()
		require.Nil(t, err)
	}()

	time.Sleep(5 * time.Second)
}
