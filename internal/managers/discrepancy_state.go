package managers

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/sirupsen/logrus"
	"log"
	"sync"
	"time"
)

type DiscrepResponse struct {
	Date        string  `json:"date"`
	FeedId      int     `json:"feed_id"`
	Discrepancy float64 `json:"discrepancy"`
	Finalized   bool    `json:"finalized"`
}

type DiscrepancyState struct {
	StartDate         time.Time
	EndDate           time.Time
	Discrepancies     []Discrepancies
	Mutex             sync.RWMutex
	ClickadillaClient ClickadillaClientInterface
	FeedState         *FeedState
	Logger            *logrus.Logger
}

func NewDiscrepancyState(clientClient ClickadillaClientInterface, logger *logrus.Logger, feedState *FeedState) *DiscrepancyState {
	return &DiscrepancyState{
		ClickadillaClient: clientClient,
		Logger:            logger,
		FeedState:         feedState,
	}
}

func (discrepancyState *DiscrepancyState) RunUpdate() {
	for {
		discrepancyState.Update()
		time.Sleep(time.Minute * 30)
	}
}

func (discrepancyState *DiscrepancyState) Update() {
	startDate := carbon.Now().SubMonths(6)

	result := make([]Discrepancies, 0)
	wg := sync.WaitGroup{}

	for i := 0; i < 6; i++ {
		wg.Add(1)
		start := startDate.AddMonths(i)
		end := start.AddMonths(1)

		go func(start, end time.Time) {
			defer func() {
				wg.Done()
			}()
			discrep, err := discrepancyState.ClickadillaClient.GetDiscrepancies(start, end)
			if err != nil {
				discrepancyState.Logger.Error(err.Error())
				return
			}
			result = append(result, discrep...)
		}(start.Carbon2Time(), end.Carbon2Time())
	}

	wg.Wait()
	discrepancyState.Discrepancies = result
	fmt.Println("execute", len(discrepancyState.Discrepancies))
}

func (discrepancyState *DiscrepancyState) GetDiscrepancies(startDate, endDate time.Time, billingTypes []string, isDsp FeedType) []DiscrepResponse {

	allFeeds := discrepancyState.FeedState.GetFeeds(billingTypes, isDsp)

	log.Println("feeds ", len(allFeeds))

	result := make([]DiscrepResponse, 0)
	groupByDiscreps := discrepancyState.groupByDate()

	log.Println("groupByDiscreps count", len(groupByDiscreps))
	_, ok := groupByDiscreps["2022-03-18"]
	fmt.Println("OK1", ok)
	_, ok = groupByDiscreps["2022-08-18"]

	start := carbon.Time2Carbon(startDate)
	end := carbon.Time2Carbon(endDate)

	for day := start; day.Lte(end); day = day.AddDay() {
		fmt.Println(day.ToDateString())
		for _, feed := range allFeeds {
			discrepValue := 1.0
			finalized := false
			discrep, ok := groupByDiscreps[day.ToDateString()][feed.Id]
			if ok {
				discrepValue = discrep.Discrepancy
				finalized = true
			}

			result = append(result, DiscrepResponse{
				Date:        day.ToDateString(),
				FeedId:      feed.Id,
				Discrepancy: discrepValue,
				Finalized:   finalized,
			})
		}
	}
	fmt.Println("len result ", len(result))
	return result
}

func (discrepancyState *DiscrepancyState) groupByDate() map[string]map[int]Discrepancies {

	result := make(map[string]map[int]Discrepancies, 0)

	for _, discrep := range discrepancyState.Discrepancies {
		date := carbon.Time2Carbon(discrep.Date).ToDateString()
		mm, ok := result[date]
		if !ok {
			mm = make(map[int]Discrepancies, 0)
		}

		mm[discrep.FeedId] = discrep
		result[date] = mm
	}

	return result
}
