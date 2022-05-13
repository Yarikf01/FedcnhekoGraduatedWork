package metric

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	typeTag   = tag.MustNewKey("metric_type")
	sourceTag = tag.MustNewKey("metric_source")
	statusTag = tag.MustNewKey("metric_status")

	countMeasure = stats.Int64("count_metric", "count", stats.UnitDimensionless)
	countView    = &view.View{
		Name:        countMeasure.Name(),
		Measure:     countMeasure,
		Description: countMeasure.Description(),
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{typeTag, sourceTag, statusTag},
	}

	latencyMeasure = stats.Int64("latency_metric", "latency", stats.UnitMilliseconds)
	latencyDist    = []float64{1, 10, 50, 100, 250, 500, 1000, 10000}
	latencyView    = &view.View{
		Name:        latencyMeasure.Name(),
		Measure:     latencyMeasure,
		Description: latencyMeasure.Description(),
		Aggregation: view.Distribution(latencyDist...),
		TagKeys:     []tag.Key{typeTag, sourceTag, statusTag},
	}

	sumMeasure = stats.Int64("stats_metric", "sum", stats.UnitDimensionless)
	sumView    = &view.View{
		Name:        sumMeasure.Name(),
		Measure:     sumMeasure,
		Description: sumMeasure.Description(),
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{typeTag, sourceTag, statusTag},
	}
)
