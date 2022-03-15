package managers

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type FeedState struct {
	Feeds             []Feed
	Mutex             sync.RWMutex
	ClickadillaClient ClickadillaClientInterface
	Logger            *logrus.Logger
}

type FeedType int

const (
	Reseller FeedType = iota
	Dsp
	All
)

func NewFeedState(clientClient ClickadillaClientInterface, logger *logrus.Logger) *FeedState {
	return &FeedState{
		ClickadillaClient: clientClient,
		Logger:            logger,
	}
}

func (fs *FeedState) GetFeeds(billingTypes []string, feedType FeedType) []Feed {
	feeds := make([]Feed, 0, 500)

	for _, feed := range fs.Feeds {
		var currentFeedType FeedType
		if feed.IsDsp {
			currentFeedType = Dsp
		} else {
			currentFeedType = Reseller
		}

		if feedType != All && currentFeedType != feedType {
			continue
		}

		if len(billingTypes) == 0 {
			feeds = append(feeds, feed)
			continue
		}

		for _, billingType := range billingTypes {
			for _, feedBillingType := range feed.Formats {
				if billingType == feedBillingType {
					feeds = append(feeds, feed)
				}
			}
		}
	}

	return feeds
}

func (fs *FeedState) RunUpdate() {
	for {
		fs.Update()
		time.Sleep(time.Minute * 2)
	}
}

func (fs *FeedState) Update() {

	fs.Logger.Info("FeedState: feeds update started")

	newFeeds, err := fs.ClickadillaClient.GetFeeds()

	if err != nil {
		fs.Logger.Error(err.Error())
		return
	}

	fs.Mutex.Lock()
	fs.Feeds = newFeeds
	fs.Mutex.Unlock()
	fs.Logger.Info("FeedState: feeds update finished")
}
