package main

import (
	"bytes"
	"sort"
	"strings"
)

const (
	creativesStartTag         = "<Creatives>"
	trackingEventsTagStartTag = "<TrackingEvents>"
	trackingEventsTagEndTag   = "</TrackingEvents>"
	videoClicksStartTag       = "<VideoClicks>"
	videoClicksEndTag         = "</VideoClicks>"
	nonLinearStartTag         = "<NonLinear>"
	nonLinearEndTag           = "</NonLinear>"
	linearEndTag              = "</Linear>"
	nonLinearAdsEndTag        = "</NonLinearAds>"
	wrapperEndTag             = "</Wrapper>"
	wrapperStartTag           = "<Wrapper>"
	inLineEndTag              = "</InLine>"
	adSystemEndTag            = "</AdSystem>"
	creativeEndTag            = "</Creative>"
	companionStartTag         = "<Companion>"
	companionEndTag           = "</Companion>"
	impressionEndTag          = "</Impression>"
	companionAdsEndTag        = "</CompanionAds>"
	adElementEndTag           = "</Ad>"
)

const (
	emtpyVast = `<VAST version=\"3.0\"><Ad><Wrapper>
	<AdSystem>prebid.org wrapper</AdSystem>
	<VASTAdTagURI><![CDATA[%v]]></VASTAdTagURI>
	<Creatives></Creatives>
	</Wrapper></Ad></VAST>`
)

type StringIndexVASTBuilder struct {
	in     string
	events *StringEvents
}

func NewStringIndexVASTBuilder(in string, events *VASTEvents) IVASTBuilder {
	return &StringIndexVASTBuilder{
		in:     in,
		events: NewStringEvents(events),
	}
}

func (builder *StringIndexVASTBuilder) Build() []byte {

	type pair struct {
		pos int
		tag string
	}

	var (
		offset  int
		pairs   = make([]pair, 0, 15)
		adBlock string
	)
	// Iterate over all <Ad> tags
	for {
		adIndex := strings.Index(builder.in[offset:], adElementEndTag)
		if adIndex == -1 {
			//last block
			adBlock = builder.in[offset:]
			break
		} else {
			adBlock = builder.in[offset : offset+adIndex+len(adElementEndTag)]
		}

		//For each <Ad> tag search for <Wrapper> or <Inline> tag and add impression
		// and error tags.
		if wrapperEndIndex := strings.LastIndex(adBlock, wrapperEndTag); wrapperEndIndex != -1 {
			if impressionEndIndex := strings.Index(adBlock, impressionEndTag); impressionEndIndex != -1 {
				pairs = append(pairs, pair{pos: offset + impressionEndIndex + len(impressionEndTag), tag: builder.events.impressions + builder.events.errors})
			} else {
				pairs = append(pairs, pair{pos: offset + wrapperEndIndex, tag: builder.events.impressions + builder.events.errors})
			}
		} else if inLineEndIndex := strings.Index(adBlock, inLineEndTag); inLineEndIndex != -1 {
			if impressionEndIndex := strings.Index(adBlock, impressionEndTag); impressionEndIndex != -1 {
				pairs = append(pairs, pair{pos: offset + impressionEndIndex + len(impressionEndTag), tag: builder.events.impressions + builder.events.errors})
			} else {
				pairs = append(pairs, pair{pos: offset + inLineEndIndex, tag: builder.events.impressions + builder.events.errors})
			}
		}

		// For each <Ad> tag iterate over <Creative> tags, add tracking and click events for Linear
		// NonLinear and Companion Ads
		creativeStartIndex := strings.Index(adBlock, creativesStartTag)
		creativeOffSet := offset + creativeStartIndex
		block := ""
		for {

			creativeIndex := strings.Index(builder.in[creativeOffSet:], creativeEndTag)
			if creativeIndex == -1 {
				//last block
				block = builder.in[creativeOffSet:]
				break
			} else {
				block = builder.in[creativeOffSet : creativeOffSet+creativeIndex+len(creativeEndTag)]
			}

			// Search for <Linear> tag and add video click and tracking events
			if linearIndex := strings.Index(block, linearEndTag); linearIndex != -1 {

				// Adding video clicks events
				if videoClickEndIndex := strings.Index(block, videoClicksEndTag); videoClickEndIndex != -1 {
					pairs = append(pairs, pair{pos: creativeOffSet + videoClickEndIndex, tag: builder.events.videoClicks})
				} else {
					pairs = append(pairs, pair{pos: creativeOffSet + linearIndex, tag: videoClicksStartTag + builder.events.videoClicks + videoClicksEndTag})
				}

				// Adding tracking events
				if trackingEventsEndIndex := strings.Index(block, trackingEventsTagEndTag); trackingEventsEndIndex != -1 {
					pairs = append(pairs, pair{pos: creativeOffSet + trackingEventsEndIndex, tag: builder.events.trackingEvents})
				} else {
					pairs = append(pairs, pair{pos: creativeOffSet + linearIndex, tag: trackingEventsTagStartTag + builder.events.trackingEvents + trackingEventsTagEndTag})
				}
			}

			// Search for <NonLinear> tag and add  non linear clicktracking and tracking events
			if nonLinearAdsIndex := strings.Index(block, nonLinearAdsEndTag); nonLinearAdsIndex != -1 {
				if nonLinearIndex := strings.Index(block, nonLinearEndTag); nonLinearIndex == -1 {
					pairs = append(pairs, pair{pos: creativeOffSet + nonLinearAdsIndex, tag: nonLinearStartTag + builder.events.nonLinearClickTracking + nonLinearEndTag})
				} else {
					localOffset := creativeOffSet
					isLastNonLinear := false
					nonLinearAds := block
					for !isLastNonLinear {
						if nonLinearIndex := strings.Index(nonLinearAds, nonLinearEndTag); nonLinearIndex == -1 {
							isLastNonLinear = true
						} else {
							pairs = append(pairs, pair{pos: localOffset + nonLinearIndex, tag: builder.events.nonLinearClickTracking})
							localOffset = localOffset + nonLinearIndex + len(nonLinearEndTag)
							nonLinearAds = nonLinearAds[nonLinearIndex+len(nonLinearEndTag):]
						}
					}
				}

				// Adding tracking events
				if trackingEventsEndIndex := strings.Index(block, trackingEventsTagEndTag); trackingEventsEndIndex != -1 {
					pairs = append(pairs, pair{pos: creativeOffSet + trackingEventsEndIndex, tag: builder.events.trackingEvents})
				} else {
					pairs = append(pairs, pair{pos: creativeOffSet + nonLinearAdsIndex, tag: trackingEventsTagStartTag + builder.events.trackingEvents + trackingEventsTagEndTag})
				}
			}

			// search for <Companion> Ads and add companion tracking events
			if companionAdsIndex := strings.Index(block, companionAdsEndTag); companionAdsIndex != -1 {
				if companionIndex := strings.Index(block, companionEndTag); companionIndex == -1 {
					pairs = append(pairs, pair{pos: creativeOffSet + companionAdsIndex, tag: companionStartTag + builder.events.companionClickThrough + companionEndTag})
				} else {
					localOffset := creativeOffSet
					isLastCompanion := false
					companionAds := block
					for !isLastCompanion {
						if companionIndex := strings.Index(companionAds, companionEndTag); companionIndex == -1 {
							isLastCompanion = true
						} else {
							pairs = append(pairs, pair{pos: localOffset + companionIndex, tag: builder.events.companionClickThrough})
							localOffset = localOffset + companionIndex + len(companionEndTag)
							companionAds = companionAds[companionIndex+len(companionEndTag):]
						}
					}
				}
			}
			creativeOffSet = creativeOffSet + len(block)
		}
		offset = offset + len(adBlock)
	}

	//sort all events position
	sort.SliceStable(pairs[:], func(i, j int) bool {
		return pairs[i].pos < pairs[j].pos
	})

	//regenerate vast xml
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)
	offset = 0
	for i := range pairs {
		if offset != pairs[i].pos {
			buf.WriteString(builder.in[offset:pairs[i].pos])
			offset = pairs[i].pos
		}
		buf.WriteString(pairs[i].tag)
	}
	buf.WriteString(builder.in[offset:])

	return []byte(buf.String())
}

// goos: darwin
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
// BenchmarkXxx-12    	   95138	     11938 ns/op	   12290 B/op	      51 allocs/op
// PASS
// ok  	_/Users/test/Test	2.130s

// Working
// 1) Divide the vastXML into multiple <Ad> tags
// 2) For each <Ad> element, search for </Inline> or </Wrapper> tag and add impression
// and error tag
// 3) For each Wrapper or Inline tag, search for all the creatives
// 4) For each creative, Look for Linear, NonLinear and CompanionAds and add
// tracking , click events
