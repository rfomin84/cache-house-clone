package managers

import (
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
	"time"
)

type ClickadillaClientMock struct {
	mock.Mock
}

func (c ClickadillaClientMock) GetFeeds() ([]Feed, error) {
	args := c.Called()
	return args.Get(0).([]Feed), args.Error(1)
}

func (c ClickadillaClientMock) GetSupplySidePlatforms() ([]SupplySidePlatform, error) {
	args := c.Called()
	return args.Get(0).([]SupplySidePlatform), args.Error(1)
}

func (c ClickadillaClientMock) GetNetworks() ([]Network, error) {
	args := c.Called()
	return args.Get(0).([]Network), args.Error(1)
}

func (c ClickadillaClientMock) GetDiscrepancies(startDate, endDate time.Time) ([]Discrepancies, error) {
	args := c.Called()
	return args.Get(0).([]Discrepancies), args.Error(1)
}

func TestUpdate(t *testing.T) {
	feed := Feed{
		Id:      1,
		Geo:     "UA",
		Formats: []string{"push", "inpage"},
		IsDsp:   true,
	}
	feeds := []Feed{feed}
	logger, hook := test.NewNullLogger()
	clickadillaClientMock := ClickadillaClientMock{}
	clickadillaClientMock.On("GetFeeds").Once().Return(feeds, nil)
	feedState := &FeedState{
		ClickadillaClient: clickadillaClientMock,
		Mutex:             sync.RWMutex{},
		Feeds:             []Feed{},
		Logger:            logger,
	}

	feedState.Update()

	assert.Equal(t, feeds, feedState.Feeds)
	assert.Equal(t, 2, len(hook.Entries))
}
