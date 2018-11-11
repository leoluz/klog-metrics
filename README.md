# Log Metrics

This application consumes pod's logs and expose them as prometheus metrics.

## Environment variables

| Name                 | Description                                                    | Mandatory |
|----------------------|----------------------------------------------------------------|:---------:|
| APP_LOG_LEVEL        | ERROR < WARN < INFO < DEBUG (Default: INFO)                    | No        |
| APP_HTTP_PORT        | The application listening HTTP Port (default: 8080)            | No        |
| APP_ERROR_REGEX      | The regex applied to identify error logs                       | Yes       |
| APP_WITH_LABEL_KEY   | Only the pods containing this label key will have log metrics  | No        |
