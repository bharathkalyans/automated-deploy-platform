# wf-collect-client





#### Steps to follow
- `make run` to start the application


#### API's 
- http://localhost:8080/api/v1/collect (POST)
```
Request Body
{
    "project_github_url": "https://github.com/bharathkalyans/sampl.app",
    "build_command": "npm run build",
    "build_out_dir": "dist"
}
```
```
Response Body
{
    "build_id": "99873458-0baf-42a4-9e9b-e5793b24f116"
}
```

- http://localhost:8080/api/v1/build/{buildId} (GET)
```
Reponse Body
{
    "build_id": "aa793c42-be5a-41f1-8ed3-13b505bfd719",
    "project_github_url": "https://github.com/bharathkalyans/github-profile-viewer-react-app",
    "events": {
        "BUILD_FAILED": {
            "reason": "GITHUB_PRIVATE",
            "timestamp": "2024-02-27T12:34:56"
        },
        "BUILD_QUEUED": {
            "timestamp": "2024-02-27T12:34:56"
        },
        "BUILD_PASSSED": {
            "timestamp": "2024-02-27T12:34:56"
        },
        "BUILD_STARTED": {
            "timestamp": "2024-02-27T12:34:56"
        },
        "DEPLOY_FAILED": {
            "reason": "Any reason which prevents the build from deploying",
            "timestamp": "2024-02-27T12:34:56"
        },
        "DEPLOY_PASSED": {
            "url": "https://localhost:7234",
            "timestamp": "2024-02-27T12:34:56",
            "branded_access_url": "https://localhost:3000/BUILD_ID"
        }
    }
}
```