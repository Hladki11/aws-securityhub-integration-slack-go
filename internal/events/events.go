package events

import (
	"github.com/slack-go/slack"
)

type SecurityHubEvent interface {
	IsAlertable() bool
	SlackMessage(consoleURL, accessPortalURL, accessRoleName string) (slack.MsgOption, slack.MsgOption)
}
