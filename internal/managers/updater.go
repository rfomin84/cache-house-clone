package managers

import (
	"context"
	"time"
)

type UpdateInstance interface {
	Update()
}

type Updater struct {
	instance UpdateInstance
	interval time.Duration
}

func NewUpdater(instance UpdateInstance, interval time.Duration) *Updater {
	return &Updater{
		instance: instance,
		interval: interval,
	}
}

func (u *Updater) StartPeriodicUpdate(ctx context.Context) {
	go func() {
		timer := time.NewTimer(u.interval)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				u.instance.Update()
			}
		}
	}()
}
