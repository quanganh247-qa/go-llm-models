package models

type ActivityLog struct {
	ActivityType   string   `json:"activity_type"`
	Timestamp      string   `json:"timestamp"`
	Duration       string   `json:"duration"`
	Details        Details  `json:"details"`
	Observations   []string `json:"observations"`
	Concerns       []string `json:"concerns"`
	FollowUpNeeded bool     `json:"follow_up_needed"`
	Notes          string   `json:"notes"`
}

type Details struct {
	Description    string `json:"description"`
	IntensityLevel string `json:"intensity_level"`
	Location       string `json:"location"`
}
