package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/antfie/FoxBot/types"
	"github.com/antfie/FoxBot/utils"
	"gopkg.in/yaml.v3"
)

type yamlTimeCheck struct {
	Frequency string `yaml:"frequency"`
	From      string `yaml:"from"`
	To        string `yaml:"to"`
}

type yamlConfig struct {
	CheckForNewVersions bool   `yaml:"check_for_new_versions"`
	DBPath              string `yaml:"db_path"`
	LogPath             string `yaml:"log_path"`
	Output              struct {
		Console bool `yaml:"console"`
		Slack   *struct {
			Token     string `yaml:"token"`
			ChannelId string `yaml:"channel_id"`
			From      string `yaml:"from"`
			To        string `yaml:"to"`
		} `yaml:"slack"`
		Telegram *struct {
			Token  string `yaml:"token"`
			ChatID string `yaml:"chat_id"`
			From   string `yaml:"from"`
			To     string `yaml:"to"`
		} `yaml:"telegram"`
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
			KeywordOnly         bool     `yaml:"keyword_only"`
			ImportantKeywords   []string `yaml:"important_keywords"`
			IgnoreURLSignatures []string `yaml:"ignore_url_signatures"`
			HTML                struct {
				ContentTags         []string `yaml:"tags"`
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
	configFile := filepath.Clean("config.yaml")

	if len(os.Args) == 2 {
		configFile = filepath.Clean(os.Args[1])

		_, err := os.Stat(configFile) //#nosec G703 -- config path is from CLI arg

		if err != nil {
			log.Fatalf("Could not open config file %q.", configFile) //#nosec G706 -- path is cleaned
		}
	}

	_, err := os.Stat(configFile) //#nosec G703 -- config path is from CLI arg

	if err != nil {
		log.Print("No config file found. Creating a new config file...")
		err := os.WriteFile(configFile, defaultConfigData, 0600) //#nosec G703 -- config path is from CLI arg

		if err != nil {
			log.Fatal(err)
		}
	}

	return parseConfigFile(configFile)
}

func parseConfigFile(configFilePath string) *types.Config {
	yamlFile, err := os.ReadFile(filepath.Clean(configFilePath)) //#nosec G703 -- config path is from CLI arg

	if err != nil {
		log.Panic(err)
	}

	config := &yamlConfig{}

	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		log.Panic(err)
	}

	return &types.Config{
		CheckForNewVersions: config.CheckForNewVersions,
		DBPath:              config.DBPath,
		LogPath:             config.LogPath,
		Output: types.Output{
			Console:  config.Output.Console,
			Slack:    parseSlack(config),
			Telegram: parseTelegram(config),
		},
		Reminders:   parseReminders(config),
		Countdown:   parseCountdown(config),
		RSS:         parseRSS(config),
		SiteChanges: parseSiteChanges(config),
	}
}

func parseSlack(config *yamlConfig) *types.Slack {
	if config.Output.Slack == nil {
		return nil
	}

	return &types.Slack{
		Token:     config.Output.Slack.Token,
		ChannelId: config.Output.Slack.ChannelId,
		Duration:  parseDuration(config.Output.Slack.From, config.Output.Slack.To),
	}
}

func parseTelegram(config *yamlConfig) *types.Telegram {
	if config.Output.Telegram == nil {
		return nil
	}

	return &types.Telegram{
		Token:    config.Output.Telegram.Token,
		ChatID:   config.Output.Telegram.ChatID,
		Duration: parseDuration(config.Output.Telegram.From, config.Output.Telegram.To),
	}
}

func parseReminders(config *yamlConfig) *types.Reminders {
	if config.Reminders == nil {
		return nil
	}

	return &types.Reminders{
		Check:     parseTimeCheck(config.Reminders.Check),
		Reminders: config.Reminders.Remidners,
	}
}

func parseCountdown(config *yamlConfig) *types.Countdown {
	if config.Countdown == nil {
		return nil
	}

	countdownTimers := make([]types.CountdownTimer, len(config.Countdown.Timers))

	for i, x := range config.Countdown.Timers {
		countdownTimers[i] = types.CountdownTimer{
			Name: x.Name,
			Date: utils.ParseDateFromString(x.Date),
		}
	}

	return &types.Countdown{
		Check:  parseTimeCheck(config.Countdown.Check),
		Timers: countdownTimers,
	}
}

func parseRSS(config *yamlConfig) *types.RSS {
	if config.RSS == nil {
		return nil
	}

	var feeds []types.RSSFeed

	for _, rssGroup := range config.RSS.Feeds {
		for _, rssFeed := range rssGroup.Feeds {
			feeds = append(feeds, types.RSSFeed{
				Group:                   rssGroup.Group,
				KeywordOnly:             rssGroup.KeywordOnly,
				ImportantKeywords:       utils.MergeStringArrays(rssGroup.ImportantKeywords, config.RSS.ImportantKeywords),
				IgnoreURLSignatures:     rssGroup.IgnoreURLSignatures,
				Name:                    rssFeed.Name,
				URL:                     rssFeed.URL,
				HTMLContentTags:         rssGroup.HTML.ContentTags,
				HTMLImportantKeywords:   rssGroup.HTML.ImportantKeywords,
				HTMLIgnoreURLSignatures: rssGroup.HTML.IgnoreURLSignatures,
			})
		}
	}

	return &types.RSS{
		Check: parseTimeCheck(config.RSS.Check),
		Feeds: feeds,
	}
}

func parseSiteChanges(config *yamlConfig) *types.SiteChange {
	if config.SiteChanges == nil {
		return nil
	}

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

	return &types.SiteChange{
		Check: parseTimeCheck(config.SiteChanges.Check),
		Sites: sites,
	}
}

func parseTimeCheck(check yamlTimeCheck) types.TimeFrequencyAndDuration {
	return types.TimeFrequencyAndDuration{
		Frequency: utils.ParseDurationFromString(check.Frequency),
		Duration:  parseDuration(check.From, check.To),
	}
}

func parseDuration(from, to string) *types.TimeDuration {
	// Both from and to need to be set
	if len(from) < 1 || len(to) < 1 {
		return nil
	}

	return &types.TimeDuration{
		From: utils.ParseTimeFromString(from),
		To:   utils.ParseTimeFromString(to),
	}
}
