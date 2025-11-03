package app

import (
	"fmt"

	awsEvent "github.com/aws/aws-lambda-go/events"
	"github.com/cruxstack/aws-securityhub-integration-slack-go/internal/events"
	"github.com/slack-go/slack"
)

type App struct {
	Config      *Config
	SlackClient *slack.Client
}

func New(cfg *Config) *App {
	return &App{
		Config:      cfg,
		SlackClient: slack.New(cfg.SlackToken),
	}
}

func (a *App) ParseEvent(e awsEvent.CloudWatchEvent) (events.SecurityHubEvent, error) {
	switch e.DetailType {
	case "Security Hub Findings - Imported":
		return events.NewSecurityHubFinding(e.Detail)
	default:
		return nil, fmt.Errorf("unknown cloudwatch event type: %s", e.DetailType)
	}
}

func (a *App) Process(evt awsEvent.CloudWatchEvent) error {
	e, err := a.ParseEvent(evt)
	if err != nil || !e.IsAlertable() {
		return err
	}
	m0, m1 := e.SlackMessage(a.Config.AwsConsoleURL, a.Config.AwsAccessPortalURL, a.Config.AwsAccessRoleName)
	_, _, err = a.SlackClient.PostMessage(a.Config.SlackChannel, m0, m1)
	return err
}
