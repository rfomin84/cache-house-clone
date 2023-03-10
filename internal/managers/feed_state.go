package managers

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type FeedState struct {
	Feeds                []Feed
	AllFeeds             []AllFeeds
	FeedsNetworks        []FeedsNetworks
	FeedsAccountManagers []FeedsAccountManagers
	Mutex                sync.RWMutex
	ClickadillaClient    ClickadillaClientInterface
	Logger               *logrus.Logger
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
	fs.Mutex.RLock()
	defer fs.Mutex.RUnlock()
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

func (fs *FeedState) GetAllFeeds(billingTypes []string, feedType FeedType) []AllFeeds {
	fs.Mutex.RLock()
	defer fs.Mutex.RUnlock()
	feeds := make([]AllFeeds, 0, 500)

	for _, feed := range fs.AllFeeds {
		var currentFeedType FeedType
		if feed.IsDsp == 1 {
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

func (fs *FeedState) GetFeedsNetworks() []FeedsNetworks {
	fs.Mutex.RLock()
	defer fs.Mutex.RUnlock()
	return fs.FeedsNetworks
}

func (fs *FeedState) GetFeedsAccountManagers() []FeedsAccountManagers {
	fs.Mutex.RLock()
	defer fs.Mutex.RUnlock()
	return fs.FeedsAccountManagers
}

func (fs *FeedState) RunUpdate() {
	for {
		fs.Update()
		time.Sleep(time.Minute * 2)
	}
}

func (fs *FeedState) RunUpdateAllFeeds() {
	for {
		fs.UpdateAllFeeds()
		time.Sleep(time.Minute * 2)
	}
}

func (fs *FeedState) RunUpdateFeedsNetworks() {
	for {
		fs.UpdateFeedNetworks()
		time.Sleep(time.Minute * 10)
	}
}

func (fs *FeedState) UpdateFeedNetworks() {
	fs.Logger.Info("FeedState: feeds-networks update started")
	newFeedsNetworks, err := fs.ClickadillaClient.GetFeedsNetworks()
	if err != nil {
		fs.Logger.Error(err.Error())
		return
	}
	fs.Mutex.Lock()
	fs.FeedsNetworks = newFeedsNetworks
	fs.Mutex.Unlock()

	fs.Logger.Info("FeedState: feeds-networks update finished")
}

func (fs *FeedState) RunUpdateFeedsAccountManagers() {
	for {
		fs.UpdateFeedsAccountManagers()
		time.Sleep(time.Minute * 60)
	}
}

func (fs *FeedState) UpdateFeedsAccountManagers() {
	fs.Logger.Info("FeedState: feeds-account-managers update started")
	newFeedsAccountManagers, err := fs.ClickadillaClient.GetFeedsAccountManagers()
	if err != nil {
		fs.Logger.Error(err.Error())
		return
	}
	fs.Mutex.Lock()
	fs.FeedsAccountManagers = newFeedsAccountManagers
	fs.Mutex.Unlock()

	fs.Logger.Info("FeedState: feeds-account-managers update finished")
}

func (fs *FeedState) Update() {

	fs.Logger.Info("FeedState: feeds update started")

	newFeeds, err := fs.ClickadillaClient.GetFeeds()
	if err != nil {
		fs.Logger.Error(err.Error())
		return
	}

	var wg sync.WaitGroup
	newFeedsTargetsMap := make(map[int]FeedTargers)
	newFeedsSupplySidePlatformsMap := make(map[int]FeedSupplySidePlatforms)
	newFeedsLabelsMap := make(map[int]FeedLabels)
	newFeedsRtbCategoriesMap := make(map[int]FeedRtbCategories)

	wg.Add(1)
	go func() {
		defer wg.Done()

		newFeedsTargets, err := fs.ClickadillaClient.GetFeedsTargets()
		if err != nil {
			fs.Logger.Error(err.Error())
			return
		}
		for _, feedTarget := range newFeedsTargets {
			newFeedsTargetsMap[feedTarget.Id] = feedTarget
		}
	}()

	wg.Add(1)
	go func() {
		wg.Done()

		newFeedsSupplySidePlatforms, err := fs.ClickadillaClient.GetFeedsSupplySidePlatforms()
		if err != nil {
			fs.Logger.Error(err.Error())
			return
		}
		for _, feedSupplySidePlatform := range newFeedsSupplySidePlatforms {
			newFeedsSupplySidePlatformsMap[feedSupplySidePlatform.Id] = feedSupplySidePlatform
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		newFeedsLabels, err := fs.ClickadillaClient.GetFeedsLabels()
		if err != nil {
			fs.Logger.Error(err.Error())
			return
		}
		for _, feedLabel := range newFeedsLabels {
			newFeedsLabelsMap[feedLabel.Id] = feedLabel
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		newFeedsRtbCategories, err := fs.ClickadillaClient.GetFeedsRtbCategories()
		if err != nil {
			fs.Logger.Error(err.Error())
			return
		}
		for _, feedRtbCategories := range newFeedsRtbCategories {
			newFeedsRtbCategoriesMap[feedRtbCategories.Id] = feedRtbCategories
		}
	}()

	wg.Wait()

	if len(newFeeds) == 0 {
		return
	}

	if len(newFeedsTargetsMap) == 0 {
		return
	}

	if len(newFeedsSupplySidePlatformsMap) == 0 {
		return
	}

	if len(newFeedsLabelsMap) == 0 {
		return
	}

	if len(newFeedsRtbCategoriesMap) == 0 {
		return
	}

	for i, feed := range newFeeds {
		newFeeds[i].Geo = newFeedsTargetsMap[feed.Id].Geo
		newFeeds[i].Formats = newFeedsTargetsMap[feed.Id].Formats
		newFeeds[i].Sources = newFeedsTargetsMap[feed.Id].Sources
		newFeeds[i].OsTypes = newFeedsTargetsMap[feed.Id].OsTypes
		newFeeds[i].SspBlacklistIncluded = newFeedsSupplySidePlatformsMap[feed.Id].SspBlacklistIncluded
		newFeeds[i].SspIds = newFeedsSupplySidePlatformsMap[feed.Id].SspIds
		newFeeds[i].Labels = newFeedsLabelsMap[feed.Id].Labels
		newFeeds[i].LabelsIds = newFeedsLabelsMap[feed.Id].LabelsIds
		newFeeds[i].RtbCategoryIds = newFeedsRtbCategoriesMap[feed.Id].RtbCategoryIds
		newFeeds[i].BrowsersWhitelist = newFeedsTargetsMap[feed.Id].BrowserWhitelist
		newFeeds[i].BrowsersBlacklist = newFeedsTargetsMap[feed.Id].BrowserBlacklist
		newFeeds[i].LanguageWhitelist = newFeedsTargetsMap[feed.Id].LanguageWhitelist
		newFeeds[i].LanguageBlacklist = newFeedsTargetsMap[feed.Id].LanguageBlacklist
		newFeeds[i].DomainWhitelist = newFeedsTargetsMap[feed.Id].DomainWhitelist
		newFeeds[i].DomainBlacklist = newFeedsTargetsMap[feed.Id].DomainBlacklist
	}

	fs.Mutex.Lock()
	fs.Feeds = newFeeds
	fs.Mutex.Unlock()
	fs.Logger.Info("FeedState: feeds update finished")
}

func (fs *FeedState) UpdateAllFeeds() {
	fs.Logger.Info("FeedState: all feeds update started")

	newAllFeeds, err := fs.ClickadillaClient.GetAllFeeds()

	if err != nil {
		fs.Logger.Error(err.Error())
		return
	}

	fs.Mutex.Lock()
	fs.AllFeeds = newAllFeeds
	fs.Mutex.Unlock()
	fs.Logger.Info("FeedState: all feeds update finished")
}
