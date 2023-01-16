package managers

import (
	"context"
	"sync"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const UPDATE_PERIOD = 30

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
	Duration          time.Duration
}

func NewDiscrepancyState(clientClient ClickadillaClientInterface, logger *logrus.Logger, feedState *FeedState) *DiscrepancyState {
	return &DiscrepancyState{
		ClickadillaClient: clientClient,
		Logger:            logger,
		FeedState:         feedState,
		Duration:          time.Minute * UPDATE_PERIOD,
	}
}

func (discrepancyState *DiscrepancyState) RunUpdate() {
	for {
		discrepancyState.Update()
		time.Sleep(discrepancyState.Duration)
	}
}

func (discrepancyState *DiscrepancyState) Update(ctx context.Context) { // context need
	discrepancyState.Logger.Info("DiscrepancyState: Discrepancies update started")

	const SIX_MONTHS = 6
	startDate := carbon.Now().SubMonths(SIX_MONTHS)

	mutex := sync.Mutex{}
	result := make([]Discrepancies, 0)

	g, _ := errgroup.WithContext(ctx)
	for i := 0; i < SIX_MONTHS; i++ {
		i = i
		g.Go(func() error {
			start := startDate.AddMonths(i)
			end := start.AddMonths(1)

			discrep, err := discrepancyState.ClickadillaClient.GetDiscrepancies(start.Carbon2Time(), end.Carbon2Time())
			if err != nil {
				discrepancyState.Logger.Error(err.Error())
				return err
			}

			mutex.Lock()
			result = append(result, discrep...)
			mutex.Unlock()

			return nil
		})
	}

	// backoff
	if err := g.Wait(); err == nil {
		discrepancyState.Mutex.Lock()
		discrepancyState.Discrepancies = result
		discrepancyState.Mutex.Unlock()

		discrepancyState.Duration = time.Minute * UPDATE_PERIOD
	} else {
		discrepancyState.Duration = time.Minute * 1
	}

	discrepancyState.Logger.Info("DiscrepancyState: Discrepancies update finished")
}

func (discrepancyState *DiscrepancyState) GetDiscrepancies(startDate, endDate time.Time, billingTypes []string, isDsp FeedType) []DiscrepResponse {
	allFeeds := discrepancyState.FeedState.GetAllFeeds(billingTypes, isDsp)

	result := make([]DiscrepResponse, 0)
	groupByDiscreps := discrepancyState.groupByDate()

	start := carbon.Time2Carbon(startDate)
	end := carbon.Time2Carbon(endDate)

	for day := start; day.Lte(end); day = day.AddDay() {
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

	return result
}

func (discrepancyState *DiscrepancyState) groupByDate() map[string]map[int]Discrepancies {
	discrepancyState.Mutex.RLock()
	defer discrepancyState.Mutex.RUnlock()

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
