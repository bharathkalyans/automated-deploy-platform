CREATE TABLE data (
  build_id VARCHAR(255) PRIMARY KEY,
  project_github_url VARCHAR(255) NOT NULL,
  events JSONB NOT NULL
);


INSERT INTO data (build_id, project_github_url, events)
VALUES (
    'a9793c42-be5a-41f1-8ed3-13b505bfd719',
    'https://github.com/bharathkalyans/github-profile-viewer-react-app',
    '{
        "BUILD_QUEUED": {"timestamp": "2024-02-27T12:34:56"},
        "BUILD_STARTED": {"timestamp": "2024-02-27T12:34:56"},
        "BUILD_FAILED": {
            "timestamp": "2024-02-27T12:34:56",
            "reason": "GITHUB_PRIVATE"
        },
        "BUILD_PASSSED": {"timestamp": "2024-02-27T12:34:56"},
        "DEPLOY_PASSED": {
            "timestamp": "2024-02-27T12:34:56",
            "branded_access_url": "https://localhost:3000/BUILD_ID",
            "url": "https://localhost:7234"
        },
        "DEPLOY_FAILED": {
            "timestamp": "2024-02-27T12:34:56",
            "reason": "Any reason which prevents the build from deploying"
        }
    }'
);


INSERT INTO data (build_id, project_github_url, events)
VALUES (
    'b1793c42-be5a-41f1-8ed3-13b505bfd719',
    'https://github.com/bharathkalyans/github-profile-viewer-react-app',
    '{
        "BUILD_QUEUED": {"timestamp": "2024-02-27T12:34:56"},
        "BUILD_STARTED": {"timestamp": "2024-02-27T12:34:56"},
        "BUILD_FAILED": {
            "timestamp": "2024-02-27T12:34:56",
            "reason": "GITHUB_PRIVATE"
        },
        "BUILD_PASSSED": {"timestamp": "2024-02-27T12:34:56"},
        "DEPLOY_PASSED": {
            "timestamp": "2024-02-27T12:34:56",
            "branded_access_url": "https://localhost:3000/BUILD_ID",
            "url": "https://localhost:7234"
        },
        "DEPLOY_FAILED": {
            "timestamp": "2024-02-27T12:34:56",
            "reason": "Any reason which prevents the build from deploying"
        }
    }'
);