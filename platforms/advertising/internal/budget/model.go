package budget

type BudgetPlan struct {
	CampaignID     string
	DailyBudget    float64
	LifetimeBudget float64
	SpentToday     float64
	LifetimeSpent  float64
	LastResetDate  string
	IsActive       bool
}

type SpendTracker struct {
	CampaignID    string
	Date          string
	HourlySpend   map[int]float64
	TotalSpent    float64
	DailyBudget   float64
	LifetimeSpent float64
}
