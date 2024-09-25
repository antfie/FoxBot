package config

import (
	"github.com/antfie/FoxBot/types"
	"github.com/antfie/FoxBot/utils"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
)

type yamlTimeCheck struct {
	Frequency string `yaml:"frequency"`
	From      string `yaml:"from"`
	To        string `yaml:"to"`
}

type yamlConfig struct {
	DBPath string `yaml:"db_path"`
	Output struct {
		Console bool `yaml:"console"`
		Slack   *struct {
			Token     string `yaml:"token"`
			ChannelId string `yaml:"channel_id"`
		} `yaml:"slack"`
	} `yaml:"output"`
	Reminders *struct {
		Check     yamlTimeCheck `yaml:"check"`
		Remidners []string      `yaml:"reminders"`
	} `yaml:"reminders"`
	Countdown *struct {
		Check  yamlTimeCheck `yaml:"check"`
		Timers []struct {
			Name string `yaml:"name"`
			Date string `yaml:"date"`
		} `yaml:"timers"`
	} `yaml:"countdown"`
	RSS *struct {
		Check             yamlTimeCheck `yaml:"check"`
		ImportantKeywords []string      `yaml:"important_keywords"`
		Feeds             []struct {
			Group               string   `yaml:"group"`
			ImportantKeywords   []string `yaml:"important_keywords"`
			IgnoreURLSignatures []string `yaml:"ignore_url_signatures"`
			HTML                struct {
				Tag                 string   `yaml:"tag"`
				ImportantKeywords   []string `yaml:"important_keywords"`
				IgnoreURLSignatures []string `yaml:"ignore_url_signatures"`
			} `yaml:"html"`
			Feeds []struct {
				Name string `yaml:"name"`
				URL  string `yaml:"url"`
			} `yaml:"feeds"`
		} `yaml:"feeds"`
	} `yaml:"rss"`
	SiteChanges *struct {
		Check yamlTimeCheck `yaml:"check"`
		Sites []struct {
			URL                        string   `yaml:"url"`
			ConnectionSuccessSignature string   `yaml:"connection_success_signature"`
			KeywordsToFind             []string `yaml:"keywords_to_find"`
			PhrasesThatMightChange     []string `yaml:"phrases_that_might_change"`
			Hash                       string   `yaml:"hash"`
		} `yaml:"sites"`
	} `yaml:"site_changes"`
}

func Load(defaultConfigData []byte) *types.Config {
	configFile := "config.yaml"

	if len(os.Args) == 2 {
		configFile = os.Args[1]

		_, err := os.Stat(configFile)

		if err != nil {
			log.Fatalf("Could not open config file \"%s\".", configFile)
		}
	}

	_, err := os.Stat(configFile)

	if err != nil {
		log.Print("No config file found. Creating a new config file...")
		err := os.WriteFile(configFile, defaultConfigData, 0600)

		if err != nil {
			log.Fatal(err)
		}
	}

	return parseConfigFile(configFile)
}

func parseConfigFile(configFilePath string) *types.Config {
	yamlFile, err := os.ReadFile(path.Clean(configFilePath))

	if err != nil {
		log.Panic(err)
	}

	config := &yamlConfig{}

	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		log.Panic(err)
	}

	return &types.Config{
		DBPath: config.DBPath,
		Output: types.Output{
			Console: config.Output.Console,
			Slack:   parseSlack(config),
		},
		Reminders:   parseReminders(config),
		Countdown:   parseCountdown(config),
		RSS:         parseRSS(config),
		SiteChanges: parseSiteChanges(config),
	}
}

func parseSlack(config *yamlConfig) *types.Slack {
	var slack *types.Slack

	if config.Output.Slack != nil {
		slack = &types.Slack{
			Token:     config.Output.Slack.Token,
			ChannelId: config.Output.Slack.ChannelId,
		}
	}
	return slack
}

func parseReminders(config *yamlConfig) *types.Reminders {
	var reminders *types.Reminders

	if config.Reminders != nil {
		if len(config.Reminders.Remidners) > 0 {
			reminders = &types.Reminders{
				Check:     parseTimeCheck(config.Reminders.Check),
				Reminders: config.Reminders.Remidners,
			}
		}
	}

	return reminders
}

func parseCountdown(config *yamlConfig) *types.Countdown {
	var countdown *types.Countdown

	if config.Countdown != nil {
		countdownTimers := make([]types.CountdownTimer, len(config.Countdown.Timers))
		for i, x := range config.Countdown.Timers {
			countdownTimers[i] = types.CountdownTimer{
				Name: x.Name,
				Date: utils.ParseDateFromString(x.Date),
			}
		}

		countdown = &types.Countdown{
			Check:  parseTimeCheck(config.Countdown.Check),
			Timers: countdownTimers,
		}
	}
	return countdown
}

func parseRSS(config *yamlConfig) *types.RSS {
	var rss *types.RSS

	if config.RSS != nil {
		var feeds []types.RSSFeed

		for _, rssGroup := range config.RSS.Feeds {
			for _, rssFeed := range rssGroup.Feeds {
				feeds = append(feeds, types.RSSFeed{
					Group:                   rssGroup.Group,
					ImportantKeywords:       utils.MergeStringArrays(rssGroup.ImportantKeywords, config.RSS.ImportantKeywords),
					IgnoreURLSignatures:     rssGroup.IgnoreURLSignatures,
					Name:                    rssFeed.Name,
					URL:                     rssFeed.URL,
					HTMLTag:                 rssGroup.HTML.Tag,
					HTMLImportantKeywords:   rssGroup.HTML.ImportantKeywords,
					HTMLIgnoreURLSignatures: rssGroup.HTML.IgnoreURLSignatures,
				})
			}
		}

		rss = &types.RSS{
			Check: parseTimeCheck(config.RSS.Check),
			Feeds: feeds,
		}
	}

	return rss
}

func parseSiteChanges(config *yamlConfig) *types.SiteChange {
	var siteChange *types.SiteChange

	if config.SiteChanges != nil {
		sites := make([]types.SiteChangeSite, len(config.SiteChanges.Sites))

		for i, x := range config.SiteChanges.Sites {
			sites[i] = types.SiteChangeSite{
				URL:                        x.URL,
				ConnectionSuccessSignature: x.ConnectionSuccessSignature,
				KeywordsToFind:             x.KeywordsToFind,
				PhrasesThatMightChange:     x.PhrasesThatMightChange,
				Hash:                       x.Hash,
			}
		}

		siteChange = &types.SiteChange{
			Check: parseTimeCheck(config.SiteChanges.Check),
			Sites: sites,
		}
	}

	return siteChange
}

func parseTimeCheck(check yamlTimeCheck) types.TimeFrequencyAndDuration {
	duration := types.TimeFrequencyAndDuration{
		Frequency: utils.ParseDuarionFromString(check.Frequency),
	}

	if len(check.From) > 0 && len(check.To) > 0 {
		duration.Duration = &types.TimeDuration{
			From: utils.ParseTimeFromString(check.From),
			To:   utils.ParseTimeFromString(check.To),
		}
	}

	return duration
}
