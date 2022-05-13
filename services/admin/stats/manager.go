package stats

import (
	"context"

	"github.com/hashicorp/go-multierror"

	"github.com/Yarikf01/graduatedwork/services/admin"
	"github.com/Yarikf01/graduatedwork/services/repo"
)

//go:generate mockery --name Manager --outpkg statsmocks --output ./statsmocks --dir .
type Manager interface {
	GetStats(ctx context.Context) (admin.DataStat, error)
}

type Config struct {
	StatsDB repo.Stats
}

const poolSize = 10

func NewManager(cfg Config) Manager {
	db := cfg.StatsDB
	return &manager{
		extractors: []func(context.Context, *admin.DataStat) error{
			//users
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.PublicUserCount, err = db.CountPublicUsers(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.PrivateUserCount, err = db.CountPrivateUsers(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.ActiveUserCount, err = db.CountActiveUsers(ctx)
				return
			},
			//reviews
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.PublicReadyPOIReviewCount, err = db.CountPublicReadyPOIReviews(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.PrivateReadyPOIReviewCount, err = db.CountPrivateReadyPOIReviews(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.PublicReadyKitchenReviewCount, err = db.CountPublicReadyKitchenReviews(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.PrivateReadyKitchenReviewCount, err = db.CountPrivateReadyKitchenReviews(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.TotalNotReadyReviewCount, err = db.CountIncompleteReviews(ctx)
				return
			},
			//followers
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.PendingActiveFollowersCount, err = db.CountPendingActiveFollowersExcludeAutoFollow(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.AcceptedActiveFollowersCount, err = db.CountAcceptedActiveFollowersExcludeAutoFollow(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.PendingFollowingCount, err = db.CountPendingFollowings(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.AcceptedFollowingCount, err = db.CountAcceptedFollowings(ctx)
				return
			},
			//places
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.FoursquareCompletedPlaceCount, err = db.CountFoursquareCompletePlaces(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.FoursquareUncompletedPlaceCount, err = db.CountFoursquareUncompletedPlaces(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.HereCompletedPlaceCount, err = db.CountHereCompletePlaces(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.HereUncompletedPlaceCount, err = db.CountHereUncompletedPlaces(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.PlaceWithReviewCount, err = db.CountPlacesWithReview(ctx)
				return
			},
			//complaints
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.ReviewComplaintCount, err = db.CountReviewComplaints(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.UserComplaintCount, err = db.CountUserComplaints(ctx)
				return
			},
			//media
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.TotalPhotosCount, err = db.CountPhotos(ctx)
				return
			},
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.TotalVideosCount, err = db.CountVideos(ctx)
				return
			},
			//likes
			func(ctx context.Context, data *admin.DataStat) (err error) {
				data.TotalLikesCount, err = db.CountLikes(ctx)
				return
			},
		},
		c: make(chan func(context.Context, *admin.DataStat) error, poolSize),
	}
}

type manager struct {
	extractors []func(context.Context, *admin.DataStat) error
	c          chan func(context.Context, *admin.DataStat) error
}

func (m *manager) GetStats(ctx context.Context) (admin.DataStat, error) {
	rg := &multierror.Group{}
	var data admin.DataStat
	for _, fn := range m.extractors {
		m.c <- fn
		rg.Go(func() error {
			f := <-m.c
			return f(ctx, &data)
		})
	}
	if err := rg.Wait().ErrorOrNil(); err != nil {
		return admin.DataStat{}, err
	}
	return data, nil
}
