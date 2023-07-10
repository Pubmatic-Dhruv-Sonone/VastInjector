package main

import "testing"

func BenchmarkXxx(b *testing.B) {

	for i := 0; i < b.N; i++ {
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

		//	txt := strings.ReplaceAll(text[0], `<VAST version="3.0">`, `<VAST version="4.2" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns="http://www.iab.com/VAST">`)
		builder := NewStringIndexVASTBuilder(str, events)
		//builder := NewETreeVASTBuilder(text[0], events)
		builder.Build()
	}
}
