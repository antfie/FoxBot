![GitHub License](https://img.shields.io/github/license/antfie/FoxBot)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/antfie/FoxBot)
[![Go Report Card](https://goreportcard.com/badge/github.com/antfie/FoxBot)](https://goreportcard.com/report/github.com/antfie/FoxBot)
![GitHub Release](https://img.shields.io/github/v/release/antfie/FoxBot)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/antfie/foxbot/total)
![Docker Image Size](https://img.shields.io/docker/image-size/antfie/foxbot/latest)
![Docker Pulls](https://img.shields.io/docker/pulls/antfie/foxbot)

# Introducing FoxBot

If a fox was a robot it would be fantastic. This is your very own personal robot, because you are fantastic. FoxBot is your AI-free* personal assistant, diligently doing your chores in the background and letting you know about things important to you. *AI has been used to build the robot, but is not used at runtime.

## What Can It Do?

* Poll your RSS feeds and notify you of relevant topics
* **Learn what you care about** — a built-in Naive Bayes classifier learns from your feedback to filter out noise and surface the articles that matter to you
* Detect changes to websites you care about
* Daily reminders throughout the day about things to be mindful of, like drinking water
* Countdown timers
* Deliver notifications via console, Slack, or Telegram

> **Tip:** For the best experience, use Telegram. RSS notifications come with inline 👍/👎 buttons that train the classifier to understand your preferences. Over time FoxBot learns which topics you care about and suppresses the rest — no cloud services, no data leaving your device. See [Intelligence](docs/intelligence.md) for details.

## How Do I Run It?

You can run this wherever you like. Just download the appropriate binary from [here](https://github.com/antfie/FoxBot/releases/latest).

### Using Docker

```bash
docker pull antfie/foxbot
docker run --rm -it -v "$(pwd):/app" antfie/foxbot
```

See the [Deployment Guide](docs/deployment.md) for systemd, Raspberry Pi, and other options.

## What Does It Look Like?

In the console you would see something like this:

![console.png](docs/images/console.png)

However FoxBot really shines when you use it as a Slack or Telegram bot:

![slack.png](docs/images/slack.png)

## How Do I Configure It?

A [config.yaml](https://github.com/antfie/FoxBot/blob/main/config.yaml) file will be generated on first run. See the [Configuration Guide](docs/configuration.md) for full details on all settings.

## Documentation

| Document | Description |
|----------|-------------|
| [Configuration Guide](docs/configuration.md) | All config options, keyword matching, feed groups |
| [Intelligence](docs/intelligence.md) | How the Naive Bayes classifier learns from your feedback |
| [Architecture](docs/architecture.md) | System design, data flow diagrams, package structure |
| [Deployment Guide](docs/deployment.md) | Docker, systemd, Raspberry Pi, building from source |

## Running Locally

```bash
git clone https://github.com/antfie/FoxBot.git
go run github.com/antfie/FoxBot
```

## How Can I Support This?

I welcome bug reports, fixes, features and donations to keep this going.

<p>
    <a href="https://www.buymeacoffee.com/antfie" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" "height="60" width="217"></a>
</p>

## Premium Features

If you need custom features, integrations or support we can help. Just email us at foxbot@antfie.com. We currently have the following premium feature availables:

- Monitoring share prices with buy/sell notifications

# Backlog

The following is a non-commital list of items we want to work through:

- Speed test functionality, ping, ICMP
- Consider daily summaries instead of regular updates
- Solar panel monitoring
- Weather
- Home automation?

# Credits

FoxBot was created by Anthony Fielding. Alert sounds by [Material Design](https://m2.material.io/design/sound/sound-resources.html) (Google), which are licenced under [CC-BY 4.0](https://creativecommons.org/licenses/by/4.0/legalcode).
