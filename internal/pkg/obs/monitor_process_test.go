package obs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func TestProcessViewsAreRegistered(t *testing.T) {
	// reset views
	view.Unregister(processViews...)

	var registeredViews = map[string]struct {
		mesure stats.Measure
		aggr   *view.Aggregation
		tags   []tag.Key
	}{
		"process/goroutine": {
			mesure: metricProcessGoroutine,
			aggr:   view.LastValue(),
			tags:   []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
		},
		"process/memory/total_alloc": {
			mesure: metricProcessMemoryTotalAlloc,
			aggr:   view.LastValue(),
			tags:   []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
		},
		"process/memory/sys": {
			mesure: metricProcessMemorySys,
			aggr:   view.LastValue(),
			tags:   []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
		},
		"process/memory/malloc": {
			mesure: metricProcessMemoryMalloc,
			aggr:   view.LastValue(),
			tags:   []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
		},
		"process/memory/free": {
			mesure: metricProcessMemoryFree,
			aggr:   view.LastValue(),
			tags:   []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
		},
		"process/memory/heap_alloc": {
			mesure: metricProcessMemoryHeapAlloc,
			aggr:   view.LastValue(),
			tags:   []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
		},
		"process/memory/heap_released": {
			mesure: metricProcessMemoryHeapReleased,
			aggr:   view.LastValue(),
			tags:   []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
		},
		"process/memory/heap_object": {
			mesure: metricProcessMemoryHeapObject,
			aggr:   view.LastValue(),
			tags:   []tag.Key{tagKeyProcessOS, tagKeyProcessArch, tagKeyProcessGoVersion},
		},
	}

	stop := MonitorProcess(time.Minute)
	defer stop()

	for name, expect := range registeredViews {
		var (
			name   = name
			expect = expect
		)
		t.Run(name, func(t *testing.T) {
			v := view.Find(name)
			require.NotNil(t, v)

			assert.Equal(t, expect.mesure, v.Measure)
			assert.Equal(t, expect.aggr.Type, v.Aggregation.Type)
			assert.ElementsMatch(t, expect.tags, v.TagKeys)
		})
	}
}

func TestProcessViewsAreFilled(t *testing.T) {
	// reset views
	view.Unregister(processViews...)

	stop := MonitorProcess(time.Millisecond * 700)
	defer stop()

	time.Sleep(time.Second) // make sure views get time to be filled

	for _, v := range processViews {
		var v = v
		t.Run(v.Name, func(t *testing.T) {
			rows, err := view.RetrieveData(v.Name)
			require.NoError(t, err)
			require.NotEmpty(t, rows)
		})
	}
}
