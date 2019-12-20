/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2019
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package query

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestFindMeasurementIdsInLink(t *testing.T) {
	cases := []struct {
		name   string
		link   string
		wanted map[int]int
	}{
		{name: "one without index 1", link: "stablenet.de/?id=33", wanted: map[int]int{0: 33}},
		{name: "one without index 2", link: "stablenet.de/?chart=555&id=34", wanted: map[int]int{0: 34}},
		{name: "several measurement ids 1", link: "stablenet.de/?chart=555&0id=34&0value1=1000&0value1=2000&1id=56&1value0=1001", wanted: map[int]int{0: 34, 1: 56}},
		{name: "several measurement ids 2", link: "stablenet.de/?chart=555&1id=34&1value1=1000&0value1=2000&0id=56&1value0=1001", wanted: map[int]int{1: 34, 0: 56}},
		{name: "several measurement ids mixed", link: "stablenet.de/?chart=555&id=34&0value1=1000&0value1=2000&1id=56&1value0=1001", wanted: map[int]int{0: 34, 1: 56}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := findMeasurementIdsInLink(tt.link)
			testify.Equal(t, len(tt.wanted), len(got), "length is different")
			for key, value := range tt.wanted {
				gotValue, ok := got[key]
				if !testify.True(t, ok, "key %d not available in got map", key) {
					continue
				}
				testify.Equal(t, value, gotValue, "for key %d the expected value %d differs from the got one %d", key, value, gotValue)
			}
		})
	}
}
