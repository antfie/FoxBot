![GitHub License](https://img.shields.io/github/license/antfie/FoxBot)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/antfie/FoxBot)
[![Go Report Card](https://goreportcard.com/badge/github.com/antfie/FoxBot)](https://goreportcard.com/report/github.com/antfie/FoxBot)
![GitHub Release](https://img.shields.io/github/v/release/antfie/FoxBot)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/antfie/foxbot/total)
![Docker Image Size](https://img.shields.io/docker/image-size/antfie/foxbot/latest)
![Docker Pulls](https://img.shields.io/docker/pulls/antfie/foxbot)

# Introducing FoxBot

If a fox was a robot it would be fantastic. This is your very own personal robot, because you are fantastic. FoxBot is your AI-free personal assistant, diligently doing your chores in the background and letting you know about things important to you.

## What Can It Do?

* Poll your RSS feeds and notify you of relvent topics
* Detect changes to websites you care about
* Daily reminders throughout the day about things to be mindful of, like drinking water
* Countdown timers

## How Do I Run It?

You can run this wherever you like. Just download the appropriate binary from [here](https://github.com/antfie/FoxBot/releases/latest).

### Using Docker

```bash
docker pull antfie/foxbot
docker run --rm -it -v "$(pwd):/app" antfie/foxbot
```

## What Does It Look Like?

In the console you would see something like this:

![console.png](docs/images/console.png)

However FoxBot really shines when you use it as a Slack bot:

![slack.png](docs/images/slack.png)

## How Do I Configure It?

There is a [config.yaml](https://github.com/antfie/FoxBot/blob/main/config.yaml) file which will be generated on first run.

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

- Document how to deploy to Raspberry Pi
- Document the configuration file
- Speed test functionality, ping, ICMP
- Reduce noises and notifications at night when sleeping
- Consider daily summaries instead of regular updates
- Solar panel monitoring
- Weather
- Home automation?

# Credits

FotBot was created by Anthony Fielding. Alert sounds by [Material Design](https://m2.material.io/design/sound/sound-resources.html) (Google), which are licenced under [CC-BY 4.0](https://creativecommons.org/licenses/by/4.0/legalcode).
