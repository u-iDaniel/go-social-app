package mailer

import "embed"

const (
	FromName            = "GoSocial"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

// The following line embeds the templates directory into the binary using compiler directives (https://gobyexample.com/embed-directive)
//
//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}
