package models

import "time"

type BuildEvent struct {
	BuildID          string           `json:"build_id"`
	ProjectGitHubURL string           `json:"project_github_url"`
	Events           map[string]Event `json:"events"`
}

type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason,omitempty"`
	URL       string    `json:"url,omitempty"`
}

type BuildRequestDetails struct {
	BuildId          string `json:"build_id"`
	ProjectGithubUrl string `json:"project_github_url"`
	BuildCommand     string `json:"build_command"`
	BuildOutDir      string `json:"build_out_dir"`
}
