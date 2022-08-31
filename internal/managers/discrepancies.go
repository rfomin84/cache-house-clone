package managers

import (
	"encoding/json"
	"time"
)

type Discrepancies struct {
	Date        time.Time `json:"date"`
	FeedId      int       `json:"feed_id"`
	Discrepancy float64   `json:"discrepancy"`
}

func (d *Discrepancies) UnmarshalJSON(dataBytes []byte) error {
	type DiscrepanciesAlias Discrepancies

	aliasValue := &struct {
		*DiscrepanciesAlias
		Date string `json:"date"`
	}{
		DiscrepanciesAlias: (*DiscrepanciesAlias)(d),
	}

	if err := json.Unmarshal(dataBytes, aliasValue); err != nil {
		return err
	}

	d.Date, _ = time.Parse("2006-01-02", aliasValue.Date)

	return nil
}

//func (d Discrepancies) MarshalJSON() ([]byte, error) {
//	log.Println("MarshalJSON")
//	type DiscrepanciesAlias Discrepancies
//
//	aliasValue := struct {
//		DiscrepanciesAlias
//		CreatedAt string `json:"date"`
//	}{
//		DiscrepanciesAlias: DiscrepanciesAlias(d),
//		CreatedAt:          d.Date.Format("2006-01-02"),
//	}
//
//	return json.Marshal(aliasValue)
//}
