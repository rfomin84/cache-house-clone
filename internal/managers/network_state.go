package managers

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type NetworkState struct {
	Networks          []Network
	Mutex             sync.RWMutex
	ClickadillaClient ClickadillaClientInterface
	Logger            *logrus.Logger
}

func NewNetworkState(clientClient ClickadillaClientInterface, logger *logrus.Logger) *NetworkState {
	return &NetworkState{
		ClickadillaClient: clientClient,
		Logger:            logger,
	}
}

func (networkState *NetworkState) RunUpdate() {
	for {
		networkState.Update()
		time.Sleep(time.Minute * 2)
	}
}

func (networkState *NetworkState) Update() {
	networkState.Logger.Info("NetworkState: networks update started")

	newNetworks, err := networkState.ClickadillaClient.GetNetworks()

	if err != nil {
		networkState.Logger.Error(err.Error())
		return
	}

	networkState.Mutex.Lock()
	networkState.Networks = newNetworks
	networkState.Mutex.Unlock()
	networkState.Logger.Info("NetworkState: networks update finished")
}

func (networkState *NetworkState) GetNetworks() []Network {
	return networkState.Networks
}
