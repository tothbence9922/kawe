# kawe
Kubernetes Application Watcher Entity - KAWE

## About this project...
 * This project is part of my thesis for the master's degree at Budapest University of Technology and Economics - Faculty of Electrical Engineering and Informatics (VIK) - Department of Automation and Applied Informatics.

## Developer notes

To run the application locally using the .env file, you can "source" the file via

`set -o allexport && source .env && set +o allexport` command (on mac).
`export KAWE_HTTP_PORT=8001 && export KAWE_PROMETHEUS_PORT=8002 && go run main.go` command (on windows, git bash).
