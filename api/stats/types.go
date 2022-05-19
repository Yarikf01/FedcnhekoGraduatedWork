package stats

type CodesStat struct {
	TotalCodes     int `json:"total_codes"`
	ActivatedCodes int `json:"activated_codes"`
	SentCodes      int `json:"sent_codes"`
}

type UsersStat struct {
	PublicUserCount  int `json:"user_account_public_count"`
	PrivateUserCount int `json:"user_account_private_count"`
	ActiveUserCount  int `json:"user_account_active_count"`
}

type ReviewsStat struct {
	PublicReadyPOIReviewCount      int `json:"public_ready_poi_review_count"`
	PrivateReadyPOIReviewCount     int `json:"private_ready_poi_review_count"`
	PublicReadyKitchenReviewCount  int `json:"public_ready_kitchen_review_count"`
	PrivateReadyKitchenReviewCount int `json:"private_ready_kitchen_review_count"`
	TotalNotReadyReviewCount       int `json:"total_not_ready_review_count"`
}

type PlacesStat struct {
	FoursquareCompletedPlaceCount   int `json:"foursquare_completed_place_count"`
	FoursquareUncompletedPlaceCount int `json:"foursquare_uncompleted_place_count"`
	HereCompletedPlaceCount         int `json:"here_completed_place_count"`
	HereUncompletedPlaceCount       int `json:"here_uncompleted_place_count"`
	PlaceWithReviewCount            int `json:"place_with_review_count"`
}

type FollowersStat struct {
	PendingActiveFollowersCount  int `json:"pending_active_followers_count"`
	AcceptedActiveFollowersCount int `json:"accepted_active_followers_count"`
	PendingFollowingCount        int `json:"pending_following_count"`
	AcceptedFollowingCount       int `json:"accepted_following_count"`
}

type ComplaintsStat struct {
	UserComplaintCount   int `json:"user_complaint_count"`
	ReviewComplaintCount int `json:"review_complaint_count"`
}

type LikesStat struct {
	TotalLikesCount int `json:"total_likes_count"`
}

type PhotosStat struct {
	TotalPhotosCount int `json:"total_photos_count"`
}

type VideosStat struct {
	TotalVideosCount int `json:"total_videos_count"`
}

type DataStat struct {
	UsersStat
	ReviewsStat
	PlacesStat
	FollowersStat
	ComplaintsStat
	LikesStat
	PhotosStat
	VideosStat
}

type UserInfoView struct {
	Name         string
	Created      string
	AccountType  string
	Email        string
	Phone        string
	Bio          string
	Nick         string
	Avatar       string
}
