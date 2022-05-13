package stats_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Yarikf01/graduatedwork/services/admin"
	"github.com/Yarikf01/graduatedwork/services/admin/stats"
	"github.com/Yarikf01/graduatedwork/services/admin/stats/statsmocks"
)

func TestHandlerGetStats(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		subj := prepareTest()
		expected := admin.DataStat{
			UsersStat: admin.UsersStat{
				PrivateUserCount: 1,
				PublicUserCount:  2,
				ActiveUserCount:  3,
			},
			ReviewsStat: admin.ReviewsStat{
				PublicReadyPOIReviewCount:      4,
				PrivateReadyPOIReviewCount:     5,
				PublicReadyKitchenReviewCount:  6,
				PrivateReadyKitchenReviewCount: 7,
				TotalNotReadyReviewCount:       8,
			},
			PlacesStat: admin.PlacesStat{
				HereCompletedPlaceCount:         9,
				HereUncompletedPlaceCount:       10,
				FoursquareCompletedPlaceCount:   11,
				FoursquareUncompletedPlaceCount: 12,
			},
			FollowersStat: admin.FollowersStat{
				PendingActiveFollowersCount:  13,
				AcceptedActiveFollowersCount: 14,
				PendingFollowingCount:        15,
				AcceptedFollowingCount:       16,
			},
			ComplaintsStat: admin.ComplaintsStat{
				UserComplaintCount:   17,
				ReviewComplaintCount: 18,
			},
		}
		subj.manager.On("GetStats", mock.Anything).Return(expected, nil)
		subj.req.GET("/admin/v1/stats").
			SetDebug(true).
			Run(subj.ech, func(resp gofight.HTTPResponse, req gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var actual admin.DataStat
				err := json.Unmarshal(resp.Body.Bytes(), &actual)
				assert.NoError(t, err)

				assert.Equal(t, expected, actual)
			})
	})

	t.Run("manager failed", func(t *testing.T) {
		subj := prepareTest()

		subj.manager.On("GetStats", mock.Anything).Return(admin.DataStat{}, fmt.Errorf(""))
		subj.req.GET("/admin/v1/stats").
			SetDebug(true).
			Run(subj.ech, func(resp gofight.HTTPResponse, req gofight.HTTPRequest) {
				assert.Equal(t, http.StatusConflict, resp.Code)
			})
	})
}

// helpers
type mocks struct {
	req     *gofight.RequestConfig
	ech     *echo.Echo
	manager *statsmocks.Manager
}

func prepareTest() *mocks {
	req := gofight.New()
	ech := echo.New()
	manager := &statsmocks.Manager{}

	stats.Assemble(ech.Group(admin.Prefix), manager)

	return &mocks{
		req:     req,
		ech:     ech,
		manager: manager,
	}
}
