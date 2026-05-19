package targeting

type Audience string

const (
	AudienceNewUsers       Audience = "new_users"
	AudienceReturningUsers Audience = "returning_users"
	AudienceHighValueUsers Audience = "high_value_users"
	AudienceCartAbandoners Audience = "cart_abandoners"
)

type Segment struct {
	ID   string
	Name string
}

type UserProfile struct {
	UserID         string
	Age            int
	Gender         string
	Location       string
	Devices        []string
	Interests      []string
	IsNewUser      bool
	IsHighValue    bool
	IsCartAbandoner bool
	TotalPurchases int
}

type TargetingRule struct {
	ID          string
	CampaignID  string
	Audience    Audience
	Segment     Segment
	Demographic *Demographic
	Devices     []string
	Locations   []string
	Interests   []string
}

type Demographic struct {
	MinAge int
	MaxAge int
	Gender string
}
