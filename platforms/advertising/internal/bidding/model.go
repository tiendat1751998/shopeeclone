package bidding

type BidStrategy string

const (
	BidStrategyManual      BidStrategy = "manual_cpc"
	BidStrategyAuto        BidStrategy = "auto_bid"
	BidStrategyEnhancedCPC BidStrategy = "enhanced_cpc"
)

type BidRequest struct {
	CampaignID string
	UserID     string
	Context    BidContext
	MaxBid     float64
}

type BidContext struct {
	Device    string
	Location  string
	Interests []string
	SessionID string
	PageURL   string
}

type BidResponse struct {
	CampaignID   string
	CreativeID   string
	BidAmount    float64
	AdRank       float64
	QualityScore float64
}

type AuctionResult struct {
	Winner      *BidResponse
	SecondPrice float64
	AllBids     []BidResponse
}
