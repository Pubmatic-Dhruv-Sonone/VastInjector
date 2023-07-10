package main

import (
	"fmt"
)

func main() {

	events := &VASTEvents{
		Errors:      getURLs(sampleURL, "/track", "/error", 1),
		Impressions: getURLs(sampleURL, "/track", "/imp", 1),
		Clicks:      getURLs(sampleURL, "/track", "/click", 1),
		TrackingEvents: map[VASTTrackingEventType][]string{
			CreativeView:     getURLs(sampleURL, "/track", "/creativeView", 1),
			Start:            getURLs(sampleURL, "/track", "/start", 1),
			FirstQuartile:    getURLs(sampleURL, "/track", "/firstQuartile", 1),
			MidPointQuartile: getURLs(sampleURL, "/track", "/midPointQuartile", 1),
			ThirdQuartile:    getURLs(sampleURL, "/track", "/thirdQuartile", 1),
			Complete:         getURLs(sampleURL, "/track", "/complete", 1),
		},
	}

	builder := NewStringIndexVASTBuilder(str, events)
	response := builder.Build()
	fmt.Println(string(response))
}
