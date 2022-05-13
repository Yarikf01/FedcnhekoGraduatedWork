package job_test

import (
	"context"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"

	"github.com/Yarikf01/graduatedwork/job"
	"github.com/Yarikf01/graduatedwork/metric/business/businessmocks"
	"github.com/Yarikf01/graduatedwork/services/admin"
	"github.com/Yarikf01/graduatedwork/services/admin/stats"
	"github.com/Yarikf01/graduatedwork/services/admin/stats/statsmocks"
)

func TestGetStatsJob(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		subj := prepareTest()
		ctx := context.TODO()

		subj.manager.On("GetStats", mock.Anything).Return(admin.DataStat{}, nil)
		subj.metricWriter.On("WriteStatPoint", mock.Anything, mock.Anything)

		job.GetStatsJob(ctx, subj.manager, subj.metricWriter)

		subj.manager.AssertCalled(t, "GetStats", mock.Anything)
		subj.metricWriter.AssertCalled(t, "WriteStatPoint", mock.Anything, mock.Anything)
	})
}

type mocks struct {
	manager      *statsmocks.Manager
	metricWriter *businessmocks.MetricWriter
}

func prepareTest() *mocks {
	ech := echo.New()
	manager := &statsmocks.Manager{}

	stats.Assemble(ech.Group(admin.Prefix), manager)

	return &mocks{
		manager:      manager,
		metricWriter: &businessmocks.MetricWriter{},
	}
}
