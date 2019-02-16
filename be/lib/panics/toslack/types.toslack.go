package toslack

type Config struct {
	Env         string
	IPAddress   string
	SlackURL    string
	WithMention bool
}

type panicsToSlack struct {
	env         string
	ipAddress   string
	slackURL    string
	withMention bool
}
