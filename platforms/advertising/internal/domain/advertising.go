package domain
import "time"

type Campaign struct { ID string `db:"id" json:"id"`; AdvertiserID string `db:"advertiser_id" json:"advertiser_id"`; Name string `db:"name" json:"name"`; Status string `db:"status" json:"status"`; Budget int64 `db:"budget" json:"budget"`; DailyBudget int64 `db:"daily_budget" json:"daily_budget"`; Spend int64 `db:"spend" json:"spend"`; BidStrategy string `db:"bid_strategy" json:"bid_strategy"`; StartTime time.Time `db:"start_time" json:"start_time"`; EndTime time.Time `db:"end_time" json:"end_time"`; CreatedAt time.Time `db:"created_at" json:"created_at"` }

type AdGroup struct { ID string `db:"id" json:"id"`; CampaignID string `db:"campaign_id" json:"campaign_id"`; Name string `db:"name" json:"name"`; Status string `db:"status" json:"status"`; BidAmount int64 `db:"bid_amount" json:"bid_amount"`; Targeting string `db:"targeting" json:"targeting"` }

type Ad struct { ID string `db:"id" json:"id"`; AdGroupID string `db:"ad_group_id" json:"ad_group_id"`; ProductID string `db:"product_id" json:"product_id"`; Status string `db:"status" json:"status"`; CreatedAt time.Time `db:"created_at" json:"created_at"` }

type Impression struct { ID string `json:"id"`; AdID string `json:"ad_id"`; UserID string `json:"user_id"`; Query string `json:"query"`; Position int `json:"position"`; Timestamp time.Time `json:"timestamp"` }

type Click struct { ID string `json:"id"`; AdID string `json:"ad_id"`; UserID string `json:"user_id"`; Cost int64 `json:"cost"`; Timestamp time.Time `json:"timestamp"` }

const ( CampaignStatusActive = "active"; CampaignStatusPaused = "paused"; CampaignStatusEnded = "ended" )
const ( BidStrategyCPC = "cpc"; BidStrategyCPM = "cpm" )
var ErrCampaignNotFound = ErrAdvertising("campaign_not_found")
type ErrAdvertising string
func (e ErrAdvertising) Error() string { return "advertising: " + string(e) }
