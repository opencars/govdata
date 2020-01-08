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
		<-time.After(3 * time.Minute)
	}
}

// Subscribe starts listening revisions and dispatch them into a channel.
func Subscribe(resourceID string, last time.Time) <-chan Revision {
	revisions := make(chan Revision)
	go listen(resourceID, last, revisions)
	return revisions
}

func listenPackage(packageID string, last time.Time, revisions chan<- Revision) {
	for {
		pkg, err := DefaultClient.PackageShow(context.Background(), packageID)
		if err != nil {
			log.Println("PackageShow:", err)
			<-time.After(30 * time.Second)
			continue
		}

		for i := 0; i < len(pkg.Resources); i++ {
			resource, err := DefaultClient.ResourceShow(context.Background(), pkg.Resources[i].ID)
			if err != nil {
				log.Println("ResourceShow:", err)
				<-time.After(30 * time.Second)
				i--
				continue
			}

			for i := len(resource.Revisions) - 1; i >= 0; i-- {
				if resource.Revisions[i].ResourceCreated.After(last) {
					revisions <- resource.Revisions[i]
				}
			}

			last = time.Now()
			<-time.After(3 * time.Minute)
		}
	}
}

func SubscribePackage(packageID string, last time.Time) <-chan Revision {
	revisions := make(chan Revision)
	go listenPackage(packageID, last, revisions)
	return revisions
}
