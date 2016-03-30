package main

import "time"

type Target struct {
	Name         string         `yaml:"-"`
	Notification bool           `yaml:"notification"`
	Watch        []*WatchConfig `yaml:"watch"`
	Script       string         `yaml:"script"`

	events chan *Event
	jitome *Jitome
}

func (target *Target) Wait() {
	for {
		event := <-target.events

		buffer := []*Event{event}
		bufferedFilesMap := map[string]int{event.Ev.Name: 1}

		timer := time.NewTimer(300 * time.Millisecond)

	outer:
		for {
			select {
			case nextEvent := <-target.events:
				// ignore a event is caused by same file path.
				if event.Ev.Name != nextEvent.Ev.Name {
					if _, exists := bufferedFilesMap[nextEvent.Ev.Name]; !exists {
						buffer = append(buffer, nextEvent)
						bufferedFilesMap[nextEvent.Ev.Name] = 1
					}
				}
			case <-timer.C:
				for _, be := range buffer {
					target.jitome.Events <- be
				}
				break outer
			}
		}
	}
}
