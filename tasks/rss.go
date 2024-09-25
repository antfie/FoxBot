package tasks

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/antfie/FoxBot/types"
	"github.com/antfie/FoxBot/utils"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type rssItem struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	PublishedDate string `xml:"pubDate"`
}

type rssStructure struct {
	RSS struct {
		Title string    `xml:"title"`
		Item  []rssItem `xml:"item"`
	} `xml:"channel"`
}

const daysNewsConsideredOld = 30

var rssMutex sync.Mutex
var rssFirstRun = true

func (c *Context) RSS() {
	if rssFirstRun {
		// Delete any old news
		c.DB.Query(fmt.Sprintf("DELETE FROM rss WHERE created > date('now', '+%d day')", daysNewsConsideredOld))
		rssFirstRun = false
	}

	if c.Config.RSS.Check.Duration != nil && !utils.IsWithinDuration(time.Now(), *c.Config.RSS.Check.Duration) {
		return
	}

	for _, feed := range c.Config.RSS.Feeds {
		go c.processRSSFeed(feed)
	}
}

func (c *Context) processRSSFeed(feed types.RSSFeed) {
	response := utils.HttpRequest("GET", feed.URL, nil, nil)

	if response == nil {
		utils.NotifyBad(fmt.Sprintf("RSS: Could not query API  %s", feed.URL))
		return
	}

	if response.StatusCode != http.StatusOK {
		utils.NotifyBad(fmt.Sprintf("RSS: API returned status of %s for %s", response.Status, feed.URL))
		return
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Panic(err)
	}

	err = response.Body.Close()

	if err != nil {
		log.Panic(err)
	}

	data := &rssStructure{}
	err = xml.Unmarshal(body, data)

	if err != nil {
		log.Panic(err)
	}

	for _, item := range data.RSS.Item {
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

		message := fmt.Sprintf("[%s]: %s - <%s>", formattedName, formattedTitle, item.Link)

		if len(foundKeyword) == 0 {
			foundKeyword = processContents(feed, item.Link)

			if len(foundKeyword) > 0 {
				message = fmt.Sprintf("[%s]: %s *%s* - <%s>", formattedName, item.Title, foundKeyword, item.Link)
			}
		}

		if len(foundKeyword) > 0 {
			c.Notify(fmt.Sprintf("📰 🚨 %s", message))
		} else {
			c.Notify(fmt.Sprintf("📰 %s", message))
		}
	}
}

func processContents(feed types.RSSFeed, url string) string {
	if len(feed.HTMLTag) < 1 {
		return ""
	}

	for _, x := range feed.HTMLIgnoreURLSignatures {
		if strings.Contains(url, x) {
			return ""
		}
	}

	itemResponse := utils.HttpRequest("GET", url, nil, nil)

	if itemResponse == nil {
		utils.NotifyBad(fmt.Sprintf("RSS: Could not query  %s", url))
		return ""
	}

	if itemResponse.StatusCode != http.StatusOK {
		utils.NotifyBad(fmt.Sprintf("RSS: Article (body) returned status of %s for %s", itemResponse.Status, url))
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(itemResponse.Body)

	if err != nil {
		utils.NotifyBad(fmt.Sprintf("RSS: HTML parsing issue for  %s", url))
		return ""
	}

	err = itemResponse.Body.Close()

	if err != nil {
		log.Panic(err)
	}

	contents := doc.Find(feed.HTMLTag).Text()

	if len(contents) < 1 {
		log.Printf("RSS: Could not find HTML contents %s", url)
		return ""
	}

	return utils.StringContainsWordIgnoreCase(contents, feed.HTMLImportantKeywords)
}

func isIgnored(feed types.RSSFeed, item rssItem, c *Context) bool {
	timestamp := utils.ParseRSSTimestampFromString(item.PublishedDate)

	if timestamp.Add(time.Hour * 24 * daysNewsConsideredOld).Before(time.Now()) {
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
