# Configuration

FoxBot uses a `config.yaml` file. On first run, a default config is generated automatically. You can also pass a custom path as a CLI argument:

```bash
./foxbot /path/to/my-config.yaml
```

## Full Config Reference

### Top-Level Settings

```yaml
# Check GitHub for new releases on startup
check_for_new_versions: true

# Path to the SQLite database (created automatically)
db_path: data.db

# Write logs to a file in addition to stderr (optional)
# log_path: foxbot.log
```

### Output

Configure where notifications are delivered. You can enable multiple outputs simultaneously.

```yaml
output:
  # Print to console with colour and audio alerts
  console: true

  # Slack bot integration
  slack:
    token: xoxb-your-token
    channel_id: C01234567
    from: 08:00   # optional: only deliver between these times
    to: 18:00

  # Telegram bot integration
  telegram:
    token: "123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
    chat_id: "123456789"
    from: 08:00
    to: 21:00
```

When `from`/`to` are set, messages are queued in SQLite and held until the time window opens. Remove both to deliver immediately at any hour.

**Setting up Slack:** Create a [Slack app](https://api.slack.com/docs/apps) and install it to your workspace to get a bot token.

**Setting up Telegram:** Message [@BotFather](https://core.telegram.org/bots#botfather) to create a bot and get a token. Then message your bot and visit `https://api.telegram.org/bot<token>/getUpdates` to find your `chat_id`.

> **Recommended:** Telegram is the best output for RSS feeds. Each RSS notification includes inline feedback buttons that train the built-in Naive Bayes classifier to learn what you care about. Over time, irrelevant articles are automatically suppressed. See [Intelligence](intelligence.md) for details.

### Reminders

Cycle through a shuffled list of motivational reminders.

```yaml
reminders:
  check:
    frequency: hourly
    from: 08:00
    to: 17:00
  reminders:
    - What are you grateful for?
    - Take time away from the keyboard
```

### Countdown Timers

Get notified when the remaining time changes (e.g. "3 months" becomes "2 months").

```yaml
countdown:
  check:
    frequency: hourly
    from: 08:00
    to: 17:00
  timers:
    - name: New Year
      date: 01/01/2026   # DD/MM/YYYY format
```

### RSS Feeds

```yaml
rss:
  check:
    frequency: half_hourly

  # Global keywords — merged into every feed group's title matching
  important_keywords:
    - FoxBot

  feeds:
    - group: BBC                    # optional group label
      keyword_only: true            # only alert on keyword matches (see below)
      important_keywords:           # group-level title keywords
        - BREAKING
      ignore_url_signatures:        # skip items with these URL patterns
        - /sport/
      html:                         # scan article body for keywords
        tags:                       # CSS selectors to extract content from
          - main
        important_keywords:         # body-specific keywords
          - breaking news story
        ignore_url_signatures:      # skip body scan for these URLs
          - /play/
      feeds:
        - name: Main
          url: https://feeds.bbci.co.uk/news/rss.xml

    - feeds:                        # feeds without a group name
        - name: xkcd
          url: https://xkcd.com/rss.xml
```

#### Keyword Matching

Keywords use **word-boundary** matching (regex `\b`), case-insensitive. This means:
- `hack` matches "hack" but **not** "hacker" or "hacking"
- Add variants explicitly: `hack`, `hacker`, `hacking`

Keywords are checked in two stages:
1. **Title** — using merged global + group `important_keywords`
2. **HTML body** — using `html.important_keywords` (only if `html.tags` is configured)

If either stage finds a match, the notification is marked with a `🚨` alert.

#### keyword_only Mode (Slack)

The `keyword_only` setting controls Slack notification filtering:

- When `keyword_only: true` on a feed group, only keyword matches are sent to **Slack**
- **Telegram** still receives all items (with feedback buttons for classifier training)
- **Console** still receives everything

This is useful for high-volume feeds where you only want Slack pings about specific topics while Telegram handles the full feed with intelligent filtering. Default is `false`.

#### ignore_url_signatures

Skip RSS items (or body scanning) when the URL contains a given substring. Useful for filtering out sport, video, or other irrelevant sections.

### Site Changes

Monitor websites for content changes.

```yaml
site_changes:
  check:
    frequency: half_hourly
  sites:
    - url: https://example.com/page
      # Verify the page loaded correctly
      connection_success_signature: "Expected Text"
      # Alert if these phrases disappear
      phrases_that_might_change:
        - "Opening soon in Spring 2026"
      # Alert on keyword appearances
      keywords_to_find:
        - sale
      # Alert if the page content hash changes
      hash: "abc123"
```

Detection methods (all optional, can combine):

| Method | What it does |
|--------|-------------|
| `connection_success_signature` | Verifies the page loaded correctly before checking anything else |
| `phrases_that_might_change` | Alerts when a known phrase disappears (substring match) |
| `keywords_to_find` | Alerts when a keyword appears (word-boundary match) |
| `hash` | Alerts when the BLAKE2b hash of the page body changes from the configured value |

### Frequency Options

Used in all `check.frequency` fields:

| Value | Interval |
|-------|----------|
| `half_hourly` | 30 minutes |
| `hourly` | 1 hour |

### Time Windows

All `from`/`to` fields use 24-hour `HH:MM` format. When set, the task only runs (or messages only deliver) within that window. Omit both to run at any time.
