package main

import (
	"time"
)

type Task struct {
	Name         string         `yaml:"-"`
	Notification bool           `yaml:"notification"`
	Watch        []*WatchConfig `yaml:"watch"`
	Script       string         `yaml:"script"`

	events chan *Event
	jitome *Jitome
}

func (task *Task) Wait() {
	for {
		event := <-task.events
		buffer := []*Event{event}
		bufferedFilesMap := map[string]int{event.Ev.Name: 1}

		timer := time.NewTimer(300 * time.Millisecond)

	outer:
		for {
			select {
			case nextEvent := <-task.events:
				// ignore a event is caused by same file path.
				if event.Ev.Name != nextEvent.Ev.Name {
					if _, exists := bufferedFilesMap[nextEvent.Ev.Name]; !exists {
						buffer = append(buffer, nextEvent)
						bufferedFilesMap[nextEvent.Ev.Name] = 1
					}
				}
			case <-timer.C:
				for _, be := range buffer {
					task.jitome.events <- be
				}
				break outer
			}
		}
	}
}
