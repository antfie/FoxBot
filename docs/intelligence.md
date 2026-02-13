# Intelligence: Local Feed Filtering

## Why Telegram?

FoxBot works with console, Slack, and Telegram — but **Telegram unlocks the full intelligence system**. When FoxBot sends you an RSS notification on Telegram, it attaches inline 👍/👎 buttons. Tap one and the built-in Naive Bayes classifier immediately learns from your feedback. Over time it understands which topics you care about and automatically suppresses the noise.

This all runs locally on your device. No data is sent to any cloud service, no API keys for third-party ML platforms, no external dependencies. Just a lightweight classifier stored in the same SQLite database FoxBot already uses.

Slack users still benefit from keyword matching and the `keyword_only` filter (see below), but without Telegram the classifier has no way to learn.

## Problem

FoxBot monitors RSS feeds and websites. Keyword matching alone produces noisy results — either too many notifications or missed articles that don't contain exact keywords. The goal is to learn what the user actually cares about and suppress irrelevant items automatically.

## Hardware Constraints

**Target device**: Raspberry Pi Zero 2 W Rev 1.0

| Spec | Value |
|------|-------|
| CPU | BCM2710A1, quad-core Cortex-A53 @ 1GHz |
| RAM | 512MB |
| GPU/NPU | None useful |
| Architecture | ARM64 |

This rules out any transformer-based model (even DistilBERT needs ~250MB and would be painfully slow). LLMs are completely out (TinyLlama 1.1B needs 2GB+).

## Approach: Naive Bayes Text Classifier

A Naive Bayes classifier is the best fit for this hardware:

- **Memory**: <5MB for the model
- **Inference**: <1ms per article
- **Training**: Incremental, no batch retraining needed
- **Dependencies**: Pure Go, no CGo or external binaries
- **Proven**: Same algorithm behind spam filters for 20+ years

### How It Works

1. Each RSS article title is tokenised into words (lowercase, split on non-alpha, drop words < 3 chars)
2. The classifier maintains per-word probability tables for "relevant" vs "irrelevant" per feed group
3. New articles get a relevance score (0-1) using log-space Bayes with Laplace smoothing
4. Score is combined with the existing keyword system to make a notify/suppress decision

### Training via User Feedback

When FoxBot sends a Telegram notification for an RSS item, it attaches two inline keyboard buttons: **👍** and **👎**. When the user taps one, the classifier updates its model immediately.

The classifier requires 30 labelled articles per feed group before it starts scoring. Until then, all articles are sent through for training.

Feedback is optional — you don't have to tap a button on every article. The classifier learns from whatever feedback you provide. If you only tap 👎 on things you don't want, it still learns. If you tap 👍 on things you enjoy, it gets better at surfacing similar items.

Duplicate presses of the same button are ignored. If you change your mind (e.g. tap 👍 then 👎), the old label is untrained and the new one is applied — only the last action counts.

### Per-Group Models

Separate classifiers are trained for each feed group (BBC, Security, etc.) since relevance criteria differ across topics. An untrained group has no effect on a trained one. If an article appears in multiple feed groups it is scored independently in each.

### Keywords as Hard Override

Keywords always punch through regardless of Bayes score. This handles the scenario where you've trained the model to suppress articles about a topic, but still want to know about exceptional events (e.g. "dies", "BREAKING", "ransomware"). The keyword list shifts from "topics I follow" to "events that always matter regardless of context."

## keyword_only Mode (Slack)

For Slack users who don't have Telegram's feedback buttons, the `keyword_only` setting provides a simpler filter:

- When `keyword_only: true` on a feed group, only keyword matches are sent to Slack
- Telegram still receives all items (with feedback buttons) for classifier training
- Console still receives everything

This is useful for high-volume feeds where you only want Slack pings about specific topics while Telegram handles the full feed with intelligent filtering.

## Notification Decision Flow

```
RSS Article
  |
  +-- Title keyword match? --> Always notify all channels (with feedback buttons)
  |
  +-- HTML body keyword match? --> Always notify all channels (with feedback buttons)
  |
  +-- No keyword match:
        +-- Bayes ready (>=30 labelled articles for this feed group)?
        |   +-- Score > 0.5 --> Notify all channels (with feedback buttons)
        |   +-- Score <= 0.5 --> Console only (suppressed)
        |
        +-- Bayes NOT ready --> Notify all channels (with feedback buttons, for training)

Slack filter (applied on top):
  keyword_only: true --> Slack only receives keyword matches
  keyword_only: false --> Slack receives everything that passes the above flow
```

## Polite RSS Fetching

FoxBot is a good citizen of the RSS ecosystem. It implements conditional HTTP requests to minimise bandwidth and server load:

- **ETag / If-None-Match**: If a feed server returns an `ETag` header, FoxBot stores it and sends `If-None-Match` on the next request. If the feed hasn't changed, the server returns `304 Not Modified` with no body.
- **Last-Modified / If-Modified-Since**: Same principle using the `Last-Modified` header.
- **429 Too Many Requests**: If a server returns 429, FoxBot backs off and does not count it as a failure.
- **Failure tracking**: Consecutive failures per feed are counted. After 10 consecutive failures FoxBot sends a notification so you know a feed is broken. The counter resets on any successful fetch.

Cache headers and failure counters are stored in SQLite and survive restarts.

## Telegram Feedback: Polling, Not Webhooks

The Telegram Bot API offers two ways to receive user input:

1. **Webhooks** — Telegram pushes updates to a public HTTPS endpoint you host
2. **`getUpdates` polling** — You call `GET /bot{TOKEN}/getUpdates` periodically and Telegram returns queued updates

**We use option 2 (polling).** This requires no public endpoint, no TLS certificate, no port forwarding, and fits the existing architecture where FoxBot already polls on timers.

### How Polling Works

```
GET https://api.telegram.org/bot{TOKEN}/getUpdates?offset={LAST_UPDATE_ID+1}&timeout=0
```

- Returns a JSON array of `Update` objects (button taps, messages, etc.)
- The `offset` parameter tells Telegram to only return updates newer than the given ID
- With `timeout=0` this is a quick non-blocking check
- FoxBot stores the last seen update ID in SQLite to survive restarts

### Sending Notifications with Inline Buttons

RSS notifications are sent individually (not batched) with a `reply_markup` parameter containing inline keyboard buttons. Each button's `callback_data` contains a prefix (`r:` for relevant, `i:` for irrelevant) and a 10-character SHA256 hash of the article URL.

Non-RSS notifications (reminders, countdowns, site changes) continue through the existing batched Telegram queue without feedback buttons.

### Feedback Processor

A background goroutine polls `getUpdates` every 30 seconds:

```
Telegram Feedback Processor (every 30s)
  +-- GET /getUpdates?offset=N
  +-- For each CallbackQuery:
  |   +-- Parse callback_data -> (relevant/irrelevant, article_hash)
  |   +-- Look up article text from DB
  |   +-- Train classifier with (text, label)
  |   +-- POST /answerCallbackQuery
  |   +-- Update stored offset
  +-- Save offset to DB
```

## Database Schema

Migration `005.sql` adds the Bayes and Telegram state tables. Migration `006.sql` adds the HTTP cache table.

```sql
-- Word frequencies per class per feed group (the trained model)
CREATE TABLE bayes_model (
    feed_group TEXT NOT NULL,
    word       TEXT NOT NULL,
    relevant   INTEGER DEFAULT 0,
    irrelevant INTEGER DEFAULT 0,
    PRIMARY KEY (feed_group, word)
);

-- Article references for feedback lookup (cleaned up after 30 days)
CREATE TABLE bayes_article (
    hash       TEXT PRIMARY KEY,
    feed_group TEXT NOT NULL,
    title      TEXT NOT NULL,
    created    DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Total document counts per class per feed group
CREATE TABLE bayes_stats (
    feed_group TEXT PRIMARY KEY,
    relevant   INTEGER DEFAULT 0,
    irrelevant INTEGER DEFAULT 0
);

-- Key-value store for Telegram polling state
CREATE TABLE telegram_state (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Conditional HTTP request cache and failure tracking
CREATE TABLE http_cache (
    url            TEXT PRIMARY KEY,
    etag           TEXT NOT NULL DEFAULT '',
    last_modified  TEXT NOT NULL DEFAULT '',
    fail_count     INTEGER NOT NULL DEFAULT 0
);
```

## Package Structure

```
bayes/
  bayes.go       -- Classifier (Train, Score, IsReady) + Tokenize
  bayes_test.go  -- Unit tests
```

The classifier reads/writes through the existing `db` package. No external ML libraries.

## Alternatives Considered

| Approach | Viable? | Notes |
|----------|---------|-------|
| Naive Bayes | **Yes (chosen)** | Pure Go, <5MB, <1ms inference |
| TF-IDF + cosine similarity | Yes | Good complement, similar footprint |
| FastText | Marginal | Needs CGo or shelling out to C++ binary |
| Decision tree / random forest | Yes | More complex, marginal benefit over Bayes |
| DistilBERT / ONNX | No | ~250MB model, too slow on this CPU |
| Any LLM (TinyLlama, Phi-2, etc.) | No | Needs 2GB+ RAM minimum |

## Future Enhancements

- **TF-IDF scoring** as a second signal alongside Bayes
- **Keyword auto-discovery** from co-occurrence analysis of relevant articles
- **Confidence display** in notifications (e.g. "relevance: 87%")
- **Daily digest** for low-scoring items instead of full suppression
