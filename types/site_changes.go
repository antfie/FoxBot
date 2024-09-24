package types

type SiteChange struct {
	Check TimeFrequencyAndDuration
	Sites []SiteChangeSite
}

type SiteChangeSite struct {
	URL                        string
	ConnectionSuccessSignature string
	KeywordsToFind             []string
	PhrasesThatMightChange     []string
	Hash                       string
}
