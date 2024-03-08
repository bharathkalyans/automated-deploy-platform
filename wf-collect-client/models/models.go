package models

import "encoding/json"

type CollectRequest struct {
	ProjectGithubURL string `json:"project_github_url"`
	BuildCommand     string `json:"build_command"`
	BuildOutDir      string `json:"build_out_dir"`
}

type CollectResponse struct {
	BuildID string `json:"build_id"`
}

type BuildEvent struct {
	BuildID string `json:"build_id"`
	Event   string `json:"event"`
}

type BuildInfo struct {
	BuildID          string          `json:"build_id"`
	ProjectGithubURL string          `json:"project_github_url"`
	Events           json.RawMessage `json:"events"`
}
