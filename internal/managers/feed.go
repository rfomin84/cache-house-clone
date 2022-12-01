package managers

type Feed struct {
	Id                   int                  `json:"id"`
	Name                 string               `json:"name"`
	CreatedAt            string               `json:"created_at"`
	Geo                  string               `json:"geo"`
	Formats              []string             `json:"formats"`
	IsDsp                bool                 `json:"is_dsp"`
	SspBlacklistIncluded bool                 `json:"ssp_blacklist_included"`
	PaymentModel         string               `json:"payment_model"`
	Throttle             int                  `json:"throttle"`
	AccurateThrottle     int                  `json:"accurate_throttle"`
	AccountId            int                  `json:"account_id"`
	AccountName          string               `json:"account_name"`
	Labels               string               `json:"labels"`
	LabelsIds            string               `json:"labels_ids"`
	Discrepancy          float64              `json:"discrepancy"`
	AntiAdblock          bool                 `json:"anti_adblock"`
	AutoThrottle         bool                 `json:"auto_throttle"`
	TcidsBlacklist       string               `json:"tcids_blacklist"`
	TcidsWhitelist       string               `json:"tcids_whitelist"`
	Capping              int                  `json:"capping"`
	IpMismatch           bool                 `json:"ip_mismatch"`
	GeoMismatch          bool                 `json:"geo_mismatch"`
	UaMismatchFilter     bool                 `json:"ua_mismatch_filter"`
	IspMismatchFilter    bool                 `json:"isp_mismatch_filter"`
	TzMismatchFilter     bool                 `json:"tz_mismatch_filter"`
	CtrCapping           int                  `json:"ctr_capping"`
	ClickTtlHours        int                  `json:"click_ttl_hours"`
	SfFrom               int                  `json:"sf_from"`
	SfTo                 int                  `json:"sf_to"`
	MinScore             int                  `json:"min_score"`
	MaxScore             int                  `json:"max_score"`
	Ipv6                 bool                 `json:"ipv6"`
	OsTypes              string               `json:"os_types"`
	CurrencyCoeff        float64              `json:"currency_coeff"`
	Sources              string               `json:"sources"`
	SspIds               string               `json:"ssp_ids"`
	GeoSiteId            []GeoSiteId          `json:"geo_siteid"`
	UniquenessSettings   []UniquenessSettings `json:"uniqueness_settings"`
	RtbCategoryIds       string               `json:"rtb_category_ids"`
	ClickTtlMin          int                  `json:"click_ttl_min"`
	CacheTtlMin          *int                 `json:"cache_ttl_min"`
	KeywordsBlacklist    string               `json:"keywords_blacklist"`
	KeywordsWhitelist    string               `json:"keywords_whitelist"`
	TrackOnlyViewed      bool                 `json:"track_only_viewed"`
	GoogleBotsFilter     bool                 `json:"google_bots_filter"`
	Browsers             []Browsers           `json:"browsers"`
	LanguageFilter       LanguageFilter       `json:"language_filter"`
	ClickDelay           *float64             `json:"click_delay"`
}

type FeedTargers struct {
	Id             int            `json:"id"`
	Geo            string         `json:"geo"`
	Formats        []string       `json:"formats"`
	OsTypes        string         `json:"os_types"`
	Sources        string         `json:"sources"`
	Browsers       []Browsers     `json:"browsers"`
	LanguageFilter LanguageFilter `json:"language_filter"`
}

type FeedSupplySidePlatforms struct {
	Id                   int    `json:"id"`
	SspIds               string `json:"ssp_ids"`
	SspBlacklistIncluded bool   `json:"ssp_blacklist_included"`
}

type FeedLabels struct {
	Id        int    `json:"id"`
	Labels    string `json:"labels"`
	LabelsIds string `json:"labels_ids"`
}

type FeedRtbCategories struct {
	Id             int    `json:"id"`
	RtbCategoryIds string `json:"rtb_category_ids"`
}

type UniquenessSettings struct {
	Format         string `json:"format"`
	Clicks         int    `json:"clicks"`
	TimeSettings   int    `json:"time_settings"`
	UniqueType     string `json:"unique_type"`
	UniqueCoverage string `json:"unique_coverage"`
}

type GeoSiteId struct {
	Accept  bool     `json:"accept"`
	Spot    string   `json:"spot"`
	Country []string `json:"country"`
}

type Browsers struct {
	Accept         bool   `json:"accept"`
	Browser        string `json:"browser"`
	BrowserVersion int    `json:"browser_version"`
}

type LanguageFilter struct {
	Accept    bool     `json:"accept"`
	Languages []string `json:"languages"`
}

type FeedsNetworks struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	NetworkId   int    `json:"network_id"`
	NetworkName string `json:"network_name"`
}
