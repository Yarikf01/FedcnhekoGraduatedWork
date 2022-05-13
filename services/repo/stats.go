package repo

import (
	"context"
)

//go:generate mockery --name Stats --outpkg repomocks --output ./repomocks --dir .
type Stats interface {
	CountPublicUsers(context.Context) (int, error)
	CountPrivateUsers(context.Context) (int, error)
	CountActiveUsers(context.Context) (int, error)

	CountPublicReadyPOIReviews(context.Context) (int, error)
	CountPrivateReadyPOIReviews(context.Context) (int, error)
	CountPublicReadyKitchenReviews(context.Context) (int, error)
	CountPrivateReadyKitchenReviews(context.Context) (int, error)
	CountIncompleteReviews(context.Context) (int, error)

	CountPendingActiveFollowersExcludeAutoFollow(context.Context) (int, error)
	CountAcceptedActiveFollowersExcludeAutoFollow(context.Context) (int, error)
	CountPendingFollowings(context.Context) (int, error)
	CountAcceptedFollowings(context.Context) (int, error)

	CountFoursquareCompletePlaces(context.Context) (int, error)
	CountFoursquareUncompletedPlaces(context.Context) (int, error)
	CountHereCompletePlaces(context.Context) (int, error)
	CountHereUncompletedPlaces(context.Context) (int, error)
	CountPlacesWithReview(ctx context.Context) (int, error)

	CountReviewComplaints(context.Context) (int, error)
	CountUserComplaints(context.Context) (int, error)

	CountPhotos(ctx context.Context) (int, error)
	CountVideos(ctx context.Context) (int, error)

	CountLikes(ctx context.Context) (int, error)
}

func (db *DB) CountPublicUsers(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM user_profile WHERE account_type = 'public'")
}

func (db *DB) CountPrivateUsers(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM user_profile WHERE account_type = 'private'")
}

func (db *DB) CountActiveUsers(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(DISTINCT user_id) FROM review WHERE ready AND created >= now() - INTERVAL '7 DAYS'")
}

func (db *DB) CountPublicReadyPOIReviews(ctx context.Context) (int, error) {
	query := "SELECT count(*)" +
		" FROM review r" +
		" JOIN user_profile u ON u.id = r.user_id AND u.account_type = 'public'" +
		" WHERE NOT is_kitchen AND ready"
	return db.selectInt(ctx, query)
}

func (db *DB) CountPrivateReadyPOIReviews(ctx context.Context) (int, error) {
	query := "SELECT count(*)" +
		" FROM review r" +
		" JOIN user_profile u ON u.id = r.user_id AND u.account_type = 'private'" +
		" WHERE NOT is_kitchen AND ready"
	return db.selectInt(ctx, query)
}

func (db *DB) CountPublicReadyKitchenReviews(ctx context.Context) (int, error) {
	query := "SELECT count(*)" +
		" FROM review r" +
		" JOIN user_profile u ON u.id = r.user_id AND u.account_type = 'public'" +
		" WHERE is_kitchen AND ready"
	return db.selectInt(ctx, query)
}

func (db *DB) CountPrivateReadyKitchenReviews(ctx context.Context) (int, error) {
	query := "SELECT count(*)" +
		" FROM review r" +
		" JOIN user_profile u ON u.id = r.user_id AND u.account_type = 'private'" +
		" WHERE is_kitchen AND ready"
	return db.selectInt(ctx, query)
}

func (db *DB) CountIncompleteReviews(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM review WHERE NOT ready")
}

func (db *DB) CountPendingActiveFollowersExcludeAutoFollow(ctx context.Context) (int, error) {
	query := "SELECT count(DISTINCT f.follower_id) FROM user_follower f " +
		"JOIN user_profile u on f.follower_id = u.id " +
		"JOIN review r on r.user_id = u.id " +
		"WHERE f.status = 'pending' AND r.ready AND r.created >= now() - INTERVAL '7 DAYS' AND u.nick NOT IN ('spencer', 'sophia')"
	return db.selectInt(ctx, query)
}

func (db *DB) CountAcceptedActiveFollowersExcludeAutoFollow(ctx context.Context) (int, error) {
	query := "SELECT count(DISTINCT f.follower_id) FROM user_follower f " +
		"JOIN user_profile u on f.follower_id = u.id " +
		"JOIN review r on r.user_id = u.id " +
		"WHERE f.status = 'accepted' AND r.ready AND r.created >= now() - INTERVAL '7 DAYS' AND u.nick NOT IN ('spencer', 'sophia')"
	return db.selectInt(ctx, query)
}

func (db *DB) CountPendingFollowings(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM user_follower WHERE status = 'pending'")
}

func (db *DB) CountAcceptedFollowings(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM user_follower WHERE status = 'accepted'")
}

func (db *DB) CountFoursquareCompletePlaces(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM place WHERE location_type='foursquare' AND completed")
}

func (db *DB) CountFoursquareUncompletedPlaces(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM place WHERE location_type='foursquare' AND NOT completed")
}

func (db *DB) CountHereCompletePlaces(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM place WHERE location_type='here' AND completed")
}

func (db *DB) CountHereUncompletedPlaces(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM place WHERE location_type='here' AND NOT completed")
}

func (db *DB) CountPlacesWithReview(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(DISTINCT place_id) FROM review WHERE ready")
}

func (db *DB) CountReviewComplaints(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM complaint WHERE complaint_type = 'review'")
}

func (db *DB) CountUserComplaints(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM complaint WHERE complaint_type = 'user'")
}

func (db *DB) CountPhotos(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM review_media WHERE media_type='image'")
}

func (db *DB) CountVideos(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM review_media WHERE media_type='video'")
}

func (db *DB) CountLikes(ctx context.Context) (int, error) {
	return db.selectInt(ctx, "SELECT count(*) FROM review_like")
}
