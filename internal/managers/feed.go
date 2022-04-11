package managers

type Feed struct {
	Id                 int                  `json:"id"`
	Geo                string               `json:"geo"`
	Formats            []string             `json:"formats"`
	IsDsp              bool                 `json:"is_dsp"`
	PaymentModel       string               `json:"payment_model"`
	Throttle           int                  `json:"throttle"`
	AccurateThrottle   int                  `json:"accurate_throttle"`
	AccountId          int                  `json:"account_id"`
	Labels             string               `json:"labels"`
	Discrepancy        float64              `json:"discrepancy"`
	AntiAdblock        bool                 `json:"anti_adblock"`
	AutoThrottle       bool                 `json:"auto_throttle"`
	TcidsBlacklist     string               `json:"tcids_blacklist"`
	TcidsWhitelist     string               `json:"tcids_whitelist"`
	Capping            int                  `json:"capping"`
	IpMismatch         bool                 `json:"ip_mismatch"`
	GeoMismatch        bool                 `json:"geo_mismatch"`
	UaMismatchFilter   bool                 `json:"ua_mismatch_filter"`
	IspMismatchFilter  bool                 `json:"isp_mismatch_filter"`
	TzMismatchFilter   bool                 `json:"tz_mismatch_filter"`
	CtrCapping         int                  `json:"ctr_capping"`
	ClickTtlHours      int                  `json:"click_ttl_hours"`
	SfFrom             int                  `json:"sf_from"`
	SfTo               int                  `json:"sf_to"`
	Ipv6               bool                 `json:"ipv6"`
	OsTypes            string               `json:"os_types"`
	CurrencyCoeff      float64              `json:"currency_coeff"`
	Sources            string               `json:"sources"`
	GeoSiteId          []string             `json:"geo_site_id"`
	UniquenessSettings []UniquenessSettings `json:"uniqueness_settings"`
}

type UniquenessSettings struct {
	Format         string `json:"format"`
	Clicks         int    `json:"clicks"`
	TimeSettings   int    `json:"time_settings"`
	UniqueType     string `json:"unique_type"`
	UniqueCoverage string `json:"unique_coverage"`
}
