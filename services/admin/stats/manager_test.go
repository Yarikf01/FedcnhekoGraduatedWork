package stats_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Yarikf01/graduatedwork/services/admin"
	"github.com/Yarikf01/graduatedwork/services/admin/stats"
	"github.com/Yarikf01/graduatedwork/services/repo/repomocks"
)

func TestManagerGetStats(t *testing.T) {
	ctx := context.TODO()

	t.Run("happy path", func(t *testing.T) {
		subj, m := managerWithMocks()

		expectedStats := admin.DataStat{
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
				PlaceWithReviewCount:            33,
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
			PhotosStat: admin.PhotosStat{
				TotalPhotosCount: 21,
			},
			VideosStat: admin.VideosStat{
				TotalVideosCount: 5,
			},
			LikesStat: admin.LikesStat{
				TotalLikesCount: 12,
			},
		}

		m.statsDB.On("CountPrivateUsers", mock.Anything).Return(1, nil)
		m.statsDB.On("CountPublicUsers", mock.Anything).Return(2, nil)
		m.statsDB.On("CountActiveUsers", mock.Anything).Return(3, nil)

		m.statsDB.On("CountPublicReadyPOIReviews", mock.Anything).Return(4, nil)
		m.statsDB.On("CountPrivateReadyPOIReviews", mock.Anything).Return(5, nil)
		m.statsDB.On("CountPublicReadyKitchenReviews", mock.Anything).Return(6, nil)
		m.statsDB.On("CountPrivateReadyKitchenReviews", mock.Anything).Return(7, nil)
		m.statsDB.On("CountIncompleteReviews", mock.Anything).Return(8, nil)

		m.statsDB.On("CountHereCompletePlaces", mock.Anything).Return(9, nil)
		m.statsDB.On("CountHereUncompletedPlaces", mock.Anything).Return(10, nil)
		m.statsDB.On("CountFoursquareCompletePlaces", mock.Anything).Return(11, nil)
		m.statsDB.On("CountFoursquareUncompletedPlaces", mock.Anything).Return(12, nil)
		m.statsDB.On("CountPlacesWithReview", mock.Anything).Return(33, nil)

		m.statsDB.On("CountPendingActiveFollowersExcludeAutoFollow", mock.Anything).Return(13, nil)
		m.statsDB.On("CountAcceptedActiveFollowersExcludeAutoFollow", mock.Anything).Return(14, nil)
		m.statsDB.On("CountPendingFollowings", mock.Anything).Return(15, nil)
		m.statsDB.On("CountAcceptedFollowings", mock.Anything).Return(16, nil)

		m.statsDB.On("CountUserComplaints", mock.Anything).Return(17, nil)
		m.statsDB.On("CountReviewComplaints", mock.Anything).Return(18, nil)

		m.statsDB.On("CountPhotos", mock.Anything).Return(21, nil)
		m.statsDB.On("CountVideos", mock.Anything).Return(5, nil)
		m.statsDB.On("CountLikes", mock.Anything).Return(12, nil)

		dataStats, err := subj.GetStats(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedStats, dataStats)

		m.statsDB.AssertCalled(t, "CountPrivateUsers", mock.Anything)
		m.statsDB.AssertCalled(t, "CountPublicUsers", mock.Anything)
		m.statsDB.AssertCalled(t, "CountActiveUsers", mock.Anything)

		m.statsDB.AssertCalled(t, "CountPublicReadyPOIReviews", mock.Anything)
		m.statsDB.AssertCalled(t, "CountPrivateReadyPOIReviews", mock.Anything)
		m.statsDB.AssertCalled(t, "CountPublicReadyKitchenReviews", mock.Anything)
		m.statsDB.AssertCalled(t, "CountPrivateReadyKitchenReviews", mock.Anything)
		m.statsDB.AssertCalled(t, "CountIncompleteReviews", mock.Anything)

		m.statsDB.AssertCalled(t, "CountHereCompletePlaces", mock.Anything)
		m.statsDB.AssertCalled(t, "CountHereUncompletedPlaces", mock.Anything)
		m.statsDB.AssertCalled(t, "CountFoursquareCompletePlaces", mock.Anything)
		m.statsDB.AssertCalled(t, "CountFoursquareUncompletedPlaces", mock.Anything)
		m.statsDB.AssertCalled(t, "CountPlacesWithReview", mock.Anything)

		m.statsDB.AssertCalled(t, "CountPendingActiveFollowersExcludeAutoFollow", mock.Anything)
		m.statsDB.AssertCalled(t, "CountAcceptedActiveFollowersExcludeAutoFollow", mock.Anything)
		m.statsDB.AssertCalled(t, "CountPendingFollowings", mock.Anything)
		m.statsDB.AssertCalled(t, "CountAcceptedFollowings", mock.Anything)

		m.statsDB.AssertCalled(t, "CountUserComplaints", mock.Anything)
		m.statsDB.AssertCalled(t, "CountReviewComplaints", mock.Anything)

		m.statsDB.AssertCalled(t, "CountPhotos", mock.Anything)
		m.statsDB.AssertCalled(t, "CountVideos", mock.Anything)

		m.statsDB.AssertCalled(t, "CountLikes", mock.Anything)
	})

}

// helpers
type mmocks struct {
	statsDB *repomocks.Stats
}

func managerWithMocks() (stats.Manager, *mmocks) {
	m := mmocks{
		statsDB: &repomocks.Stats{},
	}

	cfg := stats.Config{
		StatsDB: m.statsDB,
	}

	return stats.NewManager(cfg), &m
}
