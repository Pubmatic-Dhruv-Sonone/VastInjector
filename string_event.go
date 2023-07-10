package main

import (
	"bytes"
	"sync"
)

type IVASTBuilder interface {
	Build() []byte
}

type VASTTrackingEventType int

const (
	VASTEventTypeUnknown VASTTrackingEventType = iota
	CreativeView
	Start
	FirstQuartile
	MidPointQuartile
	ThirdQuartile
	Complete
	AcceptInvitation
	Expand
	Collapse
	VASTEventsMax
)

func (t VASTTrackingEventType) String() string {
	switch t {
	case CreativeView:
		return "creativeView"
	case Start:
		return "start"
	case FirstQuartile:
		return "firstQuartile"
	case MidPointQuartile:
		return "midpoint"
	case ThirdQuartile:
		return "thirdQuartile"
	case Complete:
		return "complete"
	case AcceptInvitation:
		return "acceptInvitation"
	case Expand:
		return "expand"
	case Collapse:
		return "collapse"
	}
	return ""
}

type StringEvents struct {
	errors                 string
	impressions            string
	videoClicks            string
	nonLinearClickTracking string
	trackingEvents         string
	companionClickThrough  string
}

var bufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

type VASTEvents struct {
	Errors         []string
	Impressions    []string
	Clicks         []string
	TrackingEvents map[VASTTrackingEventType][]string //make it array of tuples <type,url>
}

func (sr *StringEvents) init(events *VASTEvents) {
	if events != nil {
		buf := bufferPool.Get().(*bytes.Buffer)

		//buf := bytes.Buffer{}
		buf.Reset()
		for _, value := range events.Errors {
			buf.WriteString("<Error><![CDATA[")
			buf.WriteString(value)
			buf.WriteString("]]></Error>")
		}
		sr.errors = buf.String()

		buf.Reset()
		for _, value := range events.Impressions {
			buf.WriteString(`<Impression><![CDATA[`)
			buf.WriteString(value)
			buf.WriteString(`]]></Impression>`)
		}
		sr.impressions = buf.String()

		buf.Reset()
		for _, value := range events.Clicks {
			buf.WriteString(`<ClickTracking><![CDATA[`)
			buf.WriteString(value)
			buf.WriteString(`]]></ClickTracking>`)
		}
		sr.videoClicks = buf.String()

		buf.Reset()
		for _, value := range events.Clicks {
			buf.WriteString(`<NonLinearClickTracking><![CDATA[`)
			buf.WriteString(value)
			buf.WriteString(`]]></NonLinearClickTracking>`)
		}

		buf.Reset()
		for _, value := range events.Clicks {
			buf.WriteString(`<CompanionClickThrough><![CDATA[`)
			buf.WriteString(value)
			buf.WriteString(`]]></CompanionClickThrough>`)
		}

		sr.companionClickThrough = buf.String()

		buf.Reset()

		//for testing purpose only
		for eventType := VASTEventTypeUnknown; eventType < VASTEventsMax; eventType++ {
			//for testing purpose only
			for _, value := range events.TrackingEvents[eventType] {
				//for eventType, urls := range events.TrackingEvents {
				//for _, value := range urls {
				buf.WriteString(`<Tracking event="`)
				buf.WriteString(eventType.String())
				buf.WriteString(`"><![CDATA[`)
				buf.WriteString(value)
				buf.WriteString(`]]></Tracking>`)
			}
		}
		sr.trackingEvents = buf.String()

		bufferPool.Put(buf)
	}
}

func NewStringEvents(events *VASTEvents) *StringEvents {
	se := &StringEvents{}
	se.init(events)
	return se
}
