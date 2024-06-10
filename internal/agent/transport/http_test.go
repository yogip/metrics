package transport

import (
	"net/http/httptest"
	"testing"

	"metrics/internal/core/model"
	"metrics/internal/infra/api/rest/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSendMetric(t *testing.T) {
	var ten int64 = 10
	var zero int64 = 0
	var big int64 = 10000000000000000
	var minusOne int64 = -1
	var minusBig int64 = -10000000000000000

	tenFloat := 10.
	zeroFloat := 0.
	bigFloat := 10000000000000000.0
	nearZeroFloat := 0.000000000000001
	minusTenFloat := -10.
	minusBigFloat := -10000000000000000.

	expectedMetrics := []*model.MetricsV2{
		{
			ID:    "counter",
			MType: model.CounterType,
			Delta: &ten,
		},
		{
			ID:    "counter_zero",
			MType: model.CounterType,
			Delta: &zero,
		},
		{
			ID:    "counter_negative",
			MType: model.CounterType,
			Delta: &minusOne,
		},
		{
			MType: model.CounterType,
			ID:    "counter_big_negative",
			Delta: &minusBig,
		},
		{
			ID:    "counter_big",
			MType: model.CounterType,
			Delta: &big,
		},
		{
			MType: model.CounterType,
			ID:    "CaunterCapitalizedName",
			Delta: &ten,
		},

		// Gauges
		{
			MType: model.GaugeType,
			ID:    "gauge",
			Value: &tenFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "gauge_zero",
			Value: &zeroFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "gauge_near_zero",
			Value: &nearZeroFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "gauge_negative",
			Value: &minusTenFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "gauge_big_negative",
			Value: &minusBigFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "gauge_big",
			Value: &bigFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "GaugeCapitalizedName",
			Value: &tenFloat,
		},
	}

	// Create a test server
	srv := gin.New()
	srv.Use(middlewares.GzipDecompressMiddleware())
	srv.POST("/updates", func(c *gin.Context) {
		var actualMetrics []*model.MetricsV2
		err := c.BindJSON(&actualMetrics)
		assert.NoError(t, err)

		assert.Equal(t, expectedMetrics, actualMetrics)
	})
	testSrv := httptest.NewServer(srv)

	client := NewClient(testSrv.URL)

	// Call the function being tested
	err := client.SendMetric(expectedMetrics)

	// Verify the result
	assert.NoError(t, err)
}
