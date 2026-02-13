package tasks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/antfie/FoxBot/crypto"
	"github.com/antfie/FoxBot/types"
	"github.com/antfie/FoxBot/utils"
)

const maxProcessedHashes = 1000

var (
	processedHashes   = make(map[string]struct{})
	processedHashesMu sync.Mutex
)

func (c *Context) SiteChanges() {
	if c.Config.SiteChanges.Check.Duration != nil && !utils.IsWithinDuration(time.Now(), *c.Config.SiteChanges.Check.Duration) {
		return
	}

	for _, site := range c.Config.SiteChanges.Sites {
		go c.checkDifference(site)
	}
}

func (c *Context) checkDifference(site types.SiteChangeSite) {
	response := utils.HttpRequest("GET", site.URL, nil, nil)

	if response == nil {
		c.NotifyBad(fmt.Sprintf("checkDifference: Could not query API %s", site.URL))
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.NotifyBad(fmt.Sprintf("checkDifference: API returned status of %s for %s", response.Status, site.URL))
		return
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Print(err)
		return
	}

	if len(site.ConnectionSuccessSignature) > 0 {
		if !strings.Contains(string(body), site.ConnectionSuccessSignature) {
			c.NotifyGood(fmt.Sprintf("Could not find success signature in response for URL: %s", site.URL))
			return
		}
	}

	lowerCaseBody := strings.ToLower(string(body))

	if len(site.KeywordsToFind) > 0 {
		for _, keyword := range site.KeywordsToFind {
			if strings.Contains(lowerCaseBody, strings.ToLower(keyword)) {
				c.NotifyGood(fmt.Sprintf("Keyword \"%s\" found for URL: %s", keyword, site.URL))
			}
		}
	}

	if len(site.PhrasesThatMightChange) > 0 {
		for _, phrase := range site.PhrasesThatMightChange {
			if !strings.Contains(lowerCaseBody, strings.ToLower(phrase)) {
				c.NotifyGood(fmt.Sprintf("Phrase \"%s\" not found for URL: %s", phrase, site.URL))
			}
		}
	}

	c.detectHashChanges(site, body)
}

func (c *Context) detectHashChanges(site types.SiteChangeSite, body []byte) {
	if len(site.Hash) < 1 {
		return
	}

	hash, err := crypto.HashDataToString(body)

	if err != nil {
		log.Print(err)
		return
	}

	if hash == site.Hash {
		return
	}

	processedHashesMu.Lock()
	defer processedHashesMu.Unlock()

	if _, seen := processedHashes[hash]; seen {
		return
	}

	c.NotifyGood(fmt.Sprintf("Body is different for URL: %s: %s", site.URL, hash))

	err = os.WriteFile(path.Clean(path.Join("data", hash)), body, 0600)

	if err != nil {
		log.Print(err)
	}

	if len(processedHashes) >= maxProcessedHashes {
		processedHashes = make(map[string]struct{})
	}

	processedHashes[hash] = struct{}{}
}
