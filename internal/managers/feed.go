package managers

type Feed struct {
	Id                   int                  `json:"id"`
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
	SiteNames            string               `json:"site_names"`
	SspNames             string               `json:"ssp_names"`
	GeoSiteId            []GeoSiteId          `json:"geo_siteid"`
	UniquenessSettings   []UniquenessSettings `json:"uniqueness_settings"`
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
