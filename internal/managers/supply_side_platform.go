package managers

type ManagerId int
type AccountId int

type SupplySidePlatform struct {
	ID                int        `json:"id"`
	Suspend           bool       `json:"suspend"`
	UUID              string     `json:"uuid"`
	Name              string     `json:"name"`
	AccountName       string     `json:"account_name"`
	Timeout           int        `json:"timeout"`
	Throttle          int        `json:"throttle"`
	TrafficController bool       `json:"traffic_controller"`
	ControllerRatio   int        `json:"controller_ratio"`
	ContentType       string     `json:"content_type"`
	ManagerID         *ManagerId `json:"manager_id"`
	AccountID         *AccountId `json:"account_id"`
	PaymentModel      string     `json:"payment_model"`
	Margin            float64    `json:"margin"`
	CurrentMargin     float64    `json:"current_margin"`
	AdFormat          string     `json:"ad_format"`
	CampaignIDs       string     `json:"campaign_ids"`
}
