# Erxes automation tool

For those who have a passion for contribution vim supported

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/65029a0136564080877ca651ed452091)](https://app.codacy.com/manual/munkhorgil/erxes-automation-contributors?utm_source=github.com&utm_medium=referral&utm_content=munkhorgil/erxes-automation-contributors&utm_campaign=Badge_Grade_Dashboard)


## Requirements
[Install Golang](https://golang.org/doc/install)
  Install Packages
```go
	go get "github.com/gizak/termui/v3"
	go get "github.com/andybrewer/mack"
```
  Set path to your erxes project
```go
const (
	Erxes             = "path/to/erxes"
	ErxesAPI          = "path/to/erxes-api"
	ErxesIntegrations = "path/to/erxes-integrations"
)
```

```go
# Replace with your daily report website
 _, err := mack.Tell(browser, `open location "https://trello.com/"`)
```

```go
 cd path/to/erxes-automation-contributors/src
 go run main.go
```
