package types

type RSS struct {
	Check TimeFrequencyAndDuration
	Feeds []RSSFeed
}

type RSSFeed struct {
	Group                   string
	KeywordOnly             bool
	ImportantKeywords       []string
	IgnoreURLSignatures     []string
	Name                    string
	URL                     string
	HTMLContentTags         []string
	HTMLImportantKeywords   []string
	HTMLIgnoreURLSignatures []string
}
