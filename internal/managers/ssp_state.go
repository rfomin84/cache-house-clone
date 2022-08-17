package managers

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type SupplySidePlatformState struct {
	SupplySidePlatforms []SupplySidePlatform
	Mutex               sync.RWMutex
	ClickadillaClient   ClickadillaClientInterface
	Logger              *logrus.Logger
}

func NewSupplySidePlatformState(clickadillaClient ClickadillaClientInterface, logger *logrus.Logger) *SupplySidePlatformState {

	return &SupplySidePlatformState{
		ClickadillaClient: clickadillaClient,
		Logger:            logger,
	}
}

func (sspState *SupplySidePlatformState) RunUpdate() {
	for {
		sspState.Update()
		time.Sleep(time.Minute * 2)
	}
}

func (sspState *SupplySidePlatformState) GetSupplySidePlatforms() []SupplySidePlatform {
	return sspState.SupplySidePlatforms
}

func (sspState *SupplySidePlatformState) Update() {
	sspState.Logger.Info("SspState: supply side platforms update started")

	newSupplySidePlatforms, err := sspState.ClickadillaClient.GetSupplySidePlatforms()

	if err != nil {
		sspState.Logger.Error(err.Error())
		return
	}

	sspState.Mutex.Lock()
	sspState.SupplySidePlatforms = newSupplySidePlatforms
	sspState.Mutex.Unlock()
	sspState.Logger.Info("SspState: supply side platforms update finished")
}
