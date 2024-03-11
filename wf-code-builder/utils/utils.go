package utils

import (
	"code-builder/models"
	"fmt"
	"math/rand"
	"net"
)

func GenerateAndCheckPort() (int, bool) {

	// Generate random port between 5000 and 10000
	port := rand.Intn(5001) + 5000

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return port, false
	}

	listener.Close()
	return port, true
}

func CreateBuildDetails(buildInfo *map[string]interface{}) models.BuildRequestDetails {
	return models.BuildRequestDetails{
		BuildId:          (*buildInfo)["build_id"].(string),
		ProjectGithubUrl: (*buildInfo)["project_github_url"].(string),
		BuildCommand:     (*buildInfo)["build_command"].(string),
		BuildOutDir:      (*buildInfo)["build_out_dir"].(string),
	}
}
