package tasks

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antfie/FoxBot/types"
	"github.com/antfie/FoxBot/utils"
	"github.com/mmcdole/gofeed"
)

const daysNewsConsideredOld = 30

var rssMutex sync.Mutex
var rssOnce sync.Once

func (c *Context) RSS() {
	rssOnce.Do(func() {
		// Delete any old news
		c.DB.Exec(fmt.Sprintf("DELETE FROM rss WHERE created < date('now', '-%d day')", daysNewsConsideredOld))
		c.DB.BayesCleanupOldArticles()
	})

	if c.Config.RSS.Check.Duration != nil && !utils.IsWithinDuration(time.Now(), *c.Config.RSS.Check.Duration) {
		return
	}

	for _, feed := range c.Config.RSS.Feeds {
		go c.processRSSFeed(feed)
	}
}

func (c *Context) processRSSFeed(feed types.RSSFeed) {
	parsedFeed, err := c.fetchFeed(feed.URL)

	if err != nil {
		return
	}

	if parsedFeed == nil {
		// 304 Not Modified
		return
	}

	for _, item := range parsedFeed.Items {
		if isIgnored(feed, item, c) {
			continue
		}

		formattedName := feed.Name
		if len(feed.Group) > 0 {
			formattedName = fmt.Sprintf("%s:%s", feed.Group, feed.Name)
		}

		foundKeyword := utils.StringContainsWordIgnoreCase(item.Title, feed.ImportantKeywords)
		formattedTitle := item.Title

		if len(foundKeyword) > 0 {
			formattedTitle = strings.ReplaceAll(item.Title, foundKeyword, fmt.Sprintf("*%s*", foundKeyword))
		}

		formattedLink := item.Link

		message := fmt.Sprintf("[%s]: %s - <%s>", formattedName, formattedTitle, formattedLink)

		// No title keyword found so look at the contents of the link
		if len(foundKeyword) == 0 {
			foundKeyword = c.processContents(feed, formattedLink)

			if len(foundKeyword) > 0 {
				message = fmt.Sprintf("[%s]: %s *%s* - <%s>", formattedName, item.Title, foundKeyword, formattedLink)
			}
		}

		if len(foundKeyword) > 0 {
			// Keyword match - always notify all channels
			c.notifyRSS(fmt.Sprintf("📰 🚨 %s", message), feed.Group, item.Link, true)
		} else if c.Bayes != nil && c.Bayes.IsReady(feed.Group) {
			// Bayes has enough data - let it decide
			score := c.Bayes.Score(feed.Group, item.Title)
			if score > 0.5 {
				c.notifyRSS(fmt.Sprintf("📰 %s", message), feed.Group, item.Link, false)
			} else {
				utils.NotifyConsole(fmt.Sprintf("📰 %s", message))
			}
		} else {
			// Bayes not ready - send everything for training
			c.notifyRSS(fmt.Sprintf("📰 %s", message), feed.Group, item.Link, false)
		}
	}
}

func (c *Context) processContents(feed types.RSSFeed, url string) string {
	if len(feed.HTMLContentTags) < 1 {
		return ""
	}

	for _, x := range feed.HTMLIgnoreURLSignatures {
		if strings.Contains(url, x) {
			return ""
		}
	}

	itemResponse := utils.HttpRequest("GET", url, nil, nil)

	if itemResponse == nil {
		c.NotifyBad(fmt.Sprintf("RSS: Could not query %s", url))
		return ""
	}

	defer itemResponse.Body.Close()

	if itemResponse.StatusCode != http.StatusOK {
		c.NotifyBad(fmt.Sprintf("RSS: Article (body) returned status of %s for %s", itemResponse.Status, url))
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(itemResponse.Body)

	if err != nil {
		c.NotifyBad(fmt.Sprintf("RSS: HTML parsing issue for %s", url))
		return ""
	}

	contents := doc.Find(strings.Join(feed.HTMLContentTags, ", ")).Text()

	if len(contents) < 1 {
		log.Printf("RSS: Could not find HTML contents %s", url)
		return ""
	}

	return utils.StringContainsWordIgnoreCase(contents, feed.HTMLImportantKeywords)
}

func articleHash(link string) string {
	h := sha256.Sum256([]byte(link))
	return hex.EncodeToString(h[:5]) // 10 hex chars
}

func (c *Context) notifyRSS(message, feedGroup, link string, isGood bool) {
	if c.Config.Output.Console {
		if isGood {
			utils.NotifyConsoleGood(message)
		} else {
			utils.NotifyConsole(message)
		}
	}

	if c.Slack != nil {
		c.DB.QueueSlackNotification(message)
	}

	if c.Telegram != nil {
		hash := articleHash(link)
		c.DB.BayesSaveArticle(hash, feedGroup, message)
		c.Telegram.SendWithFeedback(message, hash)
	}
}

func isIgnored(feed types.RSSFeed, item *gofeed.Item, c *Context) bool {
	if item.PublishedParsed != nil && item.PublishedParsed.Add(time.Hour*24*daysNewsConsideredOld).Before(time.Now()) {
		return true
	}

	for _, x := range feed.IgnoreURLSignatures {
		if strings.Contains(item.Link, x) {
			return true
		}
	}

	rssMutex.Lock()
	defer rssMutex.Unlock()

	if c.DB.IsRSSLinkInDB(item.Link) {
		return true
	}

	return false
}

const feedFailureThreshold = 10

func (c *Context) fetchFeed(feedURL string) (*gofeed.Feed, error) {
	etag, lastModified, _ := c.DB.GetHTTPCache(feedURL)

	headers := map[string]string{}

	if len(etag) > 0 {
		headers["If-None-Match"] = etag
	}

	if len(lastModified) > 0 {
		headers["If-Modified-Since"] = lastModified
	}

	response := utils.HttpRequest("GET", feedURL, headers, nil)

	if response == nil {
		failCount := c.DB.IncrementHTTPCacheFailCount(feedURL)

		if failCount == feedFailureThreshold {
			c.NotifyBad(fmt.Sprintf("RSS feed has failed %d times: %s", failCount, feedURL))
		}

		log.Printf("Could not fetch feed: %s", feedURL)
		return nil, fmt.Errorf("could not fetch feed: %s", feedURL)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	if response.StatusCode == http.StatusNotModified {
		c.DB.SetHTTPCache(feedURL, etag, lastModified)
		return nil, nil
	}

	if response.StatusCode == http.StatusTooManyRequests {
		log.Printf("RSS feed returned 429 Too Many Requests: %s", feedURL)
		return nil, fmt.Errorf("rate limited: %s", feedURL)
	}

	if response.StatusCode != http.StatusOK {
		failCount := c.DB.IncrementHTTPCacheFailCount(feedURL)

		if failCount == feedFailureThreshold {
			c.NotifyBad(fmt.Sprintf("RSS feed has failed %d times (last status: %s): %s", failCount, response.Status, feedURL))
		}

		log.Printf("RSS feed returned status %s: %s", response.Status, feedURL)
		return nil, fmt.Errorf("bad status %s: %s", response.Status, feedURL)
	}

	// Success - store cache headers and reset fail counter
	newEtag := response.Header.Get("ETag")
	newLastModified := response.Header.Get("Last-Modified")
	c.DB.SetHTTPCache(feedURL, newEtag, newLastModified)

	fp := gofeed.NewParser()
	parsedFeed, err := fp.Parse(response.Body)

	if err != nil {
		failCount := c.DB.IncrementHTTPCacheFailCount(feedURL)

		if failCount == feedFailureThreshold {
			c.NotifyBad(fmt.Sprintf("RSS feed has failed %d times (parse error): %s", failCount, feedURL))
		}

		log.Printf("Could not parse feed: %s", feedURL)
		return nil, fmt.Errorf("could not parse feed: %s", feedURL)
	}

	return parsedFeed, nil
}
