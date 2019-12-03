package govdata

import (
	"context"
	"log"
	"time"
)

func listen(id string, last time.Time, revisions chan<- Revision) {
	for {
		resource, err := DefaultClient.ResourceShow(context.Background(), id)
		if err != nil {
			log.Println("ResourceShow:", err)
			<-time.After(30 * time.Second)
			continue
		}

		for i := len(resource.Revisions) - 1; i >= 0; i-- {
			if resource.Revisions[i].ResourceCreated.After(last) {
				revisions <- resource.Revisions[i]
			}
		}

		<-time.After(3 * time.Minute)
	}
}

// Subscribe starts listening revisions and dispatch them into a channel.
func Subscribe(id string, last time.Time) <-chan Revision {
	revisions := make(chan Revision)
	go listen(id, last, revisions)
	return revisions
}
