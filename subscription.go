package govdata

import (
	"context"
	"log"
	"time"
)

func listen(resourceID string, last time.Time, revisions chan<- Revision) {
	for {
		resource, err := DefaultClient.ResourceShow(context.Background(), resourceID)
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

		last = time.Now()
		<-time.After(10 * time.Minute)
	}
}

// Subscribe starts listening revisions and dispatch them into a channel.
func Subscribe(resourceID string, last time.Time) <-chan Revision {
	revisions := make(chan Revision)
	go listen(resourceID, last, revisions)
	return revisions
}

func listenPackage(packageID string, lastModified map[string]time.Time, events chan<- Resource) {
	for {
		pkg, err := DefaultClient.PackageShow(context.Background(), packageID)
		if err != nil {
			log.Println("PackageShow:", err)
			<-time.After(30 * time.Second)
			continue
		}

		for i := 0; i < len(pkg.Resources); i++ {
			rid := pkg.Resources[i].ID

			modified, ok := lastModified[rid]
			if ok && !pkg.Resources[i].LastModified.After(modified) {
				continue
			}

			resource, err := DefaultClient.ResourceShow(context.Background(), rid)
			if err != nil {
				log.Println("ResourceShow:", err)
				<-time.After(30 * time.Second)
				i--
				continue
			}

			// Notify about new resource.
			if !ok {
				events <- *resource
				lastModified[rid] = resource.LastModified.Time
				continue
			}

			// Resource may not contain revisions.
			if len(resource.Revisions) == 0 {
				continue
			}

			// Notify about latest changes in the resource.
			events <- *resource
			lastModified[rid] = pkg.Resources[i].LastModified.Time
		}

		<-time.After(60 * time.Minute)
	}
}

func SubscribePackage(packageID string, lastModified map[string]time.Time) <-chan Resource {
	events := make(chan Resource)
	go listenPackage(packageID, lastModified, events)
	return events
}
