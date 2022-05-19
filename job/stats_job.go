package job

import (
	"context"

	"github.com/Yarikf01/graduatedwork/api/stats"
	"github.com/Yarikf01/graduatedwork/api/utils"
	"github.com/Yarikf01/graduatedwork/metric"
	"github.com/Yarikf01/graduatedwork/metric/business"
)

var extractors = []func(stats.DataStat) (map[string]string, map[string]interface{}){
	//users
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "user_account", "visibility": "public"}, map[string]interface{}{"count": data.PublicUserCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "user_account", "visibility": "private"}, map[string]interface{}{"count": data.PrivateUserCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "user_account", "type": "active"}, map[string]interface{}{"count": data.ActiveUserCount}
	},
	//reviews
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "review", "type": "poi", "visibility": "public", "status": "ready"}, map[string]interface{}{"count": data.PublicReadyPOIReviewCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "review", "type": "poi", "visibility": "private", "status": "ready"}, map[string]interface{}{"count": data.PrivateReadyPOIReviewCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "review", "type": "kitchen", "visibility": "public", "status": "ready"}, map[string]interface{}{"count": data.PublicReadyKitchenReviewCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "review", "type": "kitchen", "visibility": "private", "status": "ready"}, map[string]interface{}{"count": data.PrivateReadyKitchenReviewCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "review", "status": "not_ready"}, map[string]interface{}{"count": data.TotalNotReadyReviewCount}
	},
	//places
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "place", "status": "completed", "provider": "foursquare"}, map[string]interface{}{"count": data.FoursquareCompletedPlaceCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "place", "status": "uncompleted", "provider": "foursquare"}, map[string]interface{}{"count": data.FoursquareUncompletedPlaceCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "place", "status": "completed", "provider": "here"}, map[string]interface{}{"count": data.HereCompletedPlaceCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "place", "status": "uncompleted", "provider": "here"}, map[string]interface{}{"count": data.HereUncompletedPlaceCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "place_with_review"}, map[string]interface{}{"count": data.PlaceWithReviewCount}
	},
	//followers
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "active_followers", "status": "pending"}, map[string]interface{}{"count": data.PendingActiveFollowersCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "active_followers", "status": "accepted"}, map[string]interface{}{"count": data.AcceptedActiveFollowersCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "followings", "status": "pending"}, map[string]interface{}{"count": data.PendingFollowingCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "followings", "status": "accepted"}, map[string]interface{}{"count": data.AcceptedFollowingCount}
	},
	//complaints
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "complaint", "type": "user"}, map[string]interface{}{"count": data.UserComplaintCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "complaint", "type": "review"}, map[string]interface{}{"count": data.ReviewComplaintCount}
	},
	//media
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "photos"}, map[string]interface{}{"count": data.TotalPhotosCount}
	},
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "videos"}, map[string]interface{}{"count": data.TotalVideosCount}
	},
	//likes
	func(data stats.DataStat) (map[string]string, map[string]interface{}) {
		return map[string]string{"unit": "likes"}, map[string]interface{}{"count": data.TotalLikesCount}
	},
}

func GetStatsJob(ctx context.Context, statsManager stats.Manager, writer business.MetricWriter) {
	logger := log.FromContext(ctx)

	dataStats, err := statsManager.GetStats(ctx)
	if err != nil {
		log.WithError(logger, err).Errorf("failed to get stats from db")
		return
	}

	logger.Info("sending statistics to GCP")
	send2GCP(ctx, dataStats)

	logger.Info("sending statistics to influx cloud")
	for _, f := range extractors {
		writer.WriteStatPoint(f(dataStats))
	}
}

func send2GCP(ctx context.Context, data stats.DataStat) {
	//users
	metric.RecordSumMetric(ctx, "recon_stats_user_account_public_count", data.PublicUserCount)
	metric.RecordSumMetric(ctx, "recon_stats_user_account_private_count", data.PrivateUserCount)
	metric.RecordSumMetric(ctx, "recon_stats_user_active_count", data.ActiveUserCount)
	//reviews
	metric.RecordSumMetric(ctx, "recon_stats_public_ready_review_count", data.PublicReadyPOIReviewCount+data.PublicReadyKitchenReviewCount)
	metric.RecordSumMetric(ctx, "recon_stats_private_ready_review_count", data.PrivateReadyPOIReviewCount+data.PrivateReadyKitchenReviewCount)
	metric.RecordSumMetric(ctx, "recon_stats_total_not_ready_review_count", data.TotalNotReadyReviewCount)
	//places
	metric.RecordSumMetric(ctx, "recon_stats_foursquare_completed_place_count", data.FoursquareCompletedPlaceCount)
	metric.RecordSumMetric(ctx, "recon_stats_foursquare_uncompleted_place_count", data.FoursquareUncompletedPlaceCount)
	metric.RecordSumMetric(ctx, "recon_stats_here_completed_place_count", data.HereCompletedPlaceCount)
	metric.RecordSumMetric(ctx, "recon_stats_here_uncompleted_place_count", data.HereUncompletedPlaceCount)
	metric.RecordSumMetric(ctx, "recon_stats_place_with_review_count", data.PlaceWithReviewCount)
	//followers
	metric.RecordSumMetric(ctx, "recon_stats_pending_followers_count", data.PendingActiveFollowersCount)
	metric.RecordSumMetric(ctx, "recon_stats_accepted_followers_count", data.AcceptedActiveFollowersCount)
	//complaints
	metric.RecordSumMetric(ctx, "recon_stats_user_complaint_count", data.UserComplaintCount)
	metric.RecordSumMetric(ctx, "recon_stats_review_complaint_count", data.ReviewComplaintCount)
	//media
	metric.RecordSumMetric(ctx, "recon_stats_photos_count", data.TotalPhotosCount)
	metric.RecordSumMetric(ctx, "recon_stats_videos_count", data.TotalVideosCount)
	//likes
	metric.RecordSumMetric(ctx, "recon_stats_likes_count", data.TotalLikesCount)
}
