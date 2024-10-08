package types

type Config struct {
	CheckForNewVersions bool
	DBPath              string
	Output              Output
	Reminders           *Reminders
	Countdown           *Countdown
	RSS                 *RSS
	SiteChanges         *SiteChange
}

type Output struct {
	Console bool
	Slack   *Slack
}

type Slack struct {
	Token     string
	ChannelId string
	Duration  *TimeDuration
}
