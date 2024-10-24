package middlewares

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCheckSubnet(t *testing.T) {
	tests := []struct {
		name   string
		subnet net.IPNet
		ip     string
		code   int
	}{
		{
			name: "Test success",
			subnet: net.IPNet{
				IP:   net.IPv4(10, 0, 0, 0),
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
			ip:   "10.0.0.121:",
			code: 200,
		},
		{
			name: "Test success 2",
			subnet: net.IPNet{
				IP:   net.IPv4(10, 0, 0, 0),
				Mask: net.IPv4Mask(255, 255, 0, 0),
			},
			ip:   "10.0.1.121:",
			code: 200,
		},
		{
			name: "Test failed",
			subnet: net.IPNet{
				IP:   net.IPv4(10, 0, 0, 0),
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
			ip:   "10.0.1.121:",
			code: 403,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupSubNetApp(&tt.subnet)
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.ip
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, tt.code, result.StatusCode)
			result.Body.Close()
		})
	}
}

func setupSubNetApp(subnet *net.IPNet) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(CheckSubnet(subnet))
	router.GET(
		"/",
		func(c *gin.Context) {
			c.Status(http.StatusOK)
		},
	)
	return router
}
