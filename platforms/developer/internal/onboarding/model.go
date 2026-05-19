package onboarding

type TaskCategory string

const (
	CategorySetup TaskCategory = "setup"
	CategoryLearn TaskCategory = "learn"
	CategoryBuild TaskCategory = "build"
	CategoryDeploy TaskCategory = "deploy"
)

type OnboardingTask struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Category    TaskCategory `json:"category"`
	Required    bool         `json:"required"`
	IsCompleted bool         `json:"is_completed"`
	Order       int          `json:"order"`
}

type Template struct {
	Name        string           `json:"name"`
	ServiceType string           `json:"service_type"`
	Tasks       []OnboardingTask `json:"tasks"`
}

type Progress struct {
	TotalTasks     int            `json:"total_tasks"`
	CompletedTasks int            `json:"completed_tasks"`
	Percentage     float64        `json:"percentage"`
	Tasks          []OnboardingTask `json:"tasks"`
}
