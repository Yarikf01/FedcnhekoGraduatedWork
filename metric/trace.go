package metric

import (
	"context"
	"strconv"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opentelemetry.io/otel"

	"github.com/Yarikf01/graduatedwork/api/utils"
)

const (
	tracer           = "recon.com/trace"
	apiRequestTag    = "api_request"
	reverseGeoTag    = "reverse_geo_coding"
	textGeoSearchTag = "text_geo_search"
)

func RunWithSpan(ctx context.Context, spanName string, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracer).Start(ctx, spanName)
	defer span.End()
	return f(ctx)
}

func RecordCountMetric(ctx context.Context, value string) {
	ctx, err := tag.New(ctx, []tag.Mutator{tag.Insert(typeTag, value)}...)
	if err != nil {
		log.FromContext(ctx).With("error", err).Error("failed to create count metric tag for:", value)
	}
	stats.Record(ctx, countMeasure.M(0))
}

func RecordSumMetric(ctx context.Context, value string, num int) {
	ctx, err := tag.New(ctx, []tag.Mutator{tag.Insert(typeTag, value)}...)
	if err != nil {
		log.FromContext(ctx).With("error", err).Error("failed to create sum metric tag for:", value)
	}
	stats.Record(ctx, sumMeasure.M(int64(num)))
}

func RecordApiRequest(ctx context.Context, latency time.Duration, path string, status int) {
	ctx, err := contextWithTags(ctx, apiRequestTag, path, strconv.Itoa(status))
	if err != nil {
		log.FromContext(ctx).With("error", err).Error("failed to create api request metric tag for:", path)
	}
	stats.Record(ctx, countMeasure.M(0))
	stats.Record(ctx, latencyMeasure.M(latency.Milliseconds()))
}

func RecordReverseGeoCoding(ctx context.Context, latency time.Duration, source string, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	ctx, err := contextWithTags(ctx, reverseGeoTag, source, status)
	if err != nil {
		log.FromContext(ctx).With("error", err).Error("failed to create reverse geo metric tag for:", source)
	}
	stats.Record(ctx, countMeasure.M(0))
	stats.Record(ctx, latencyMeasure.M(latency.Milliseconds()))
}

func RecordTextGeoSearch(ctx context.Context, latency time.Duration, source string, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	ctx, err := contextWithTags(ctx, textGeoSearchTag, source, status)
	if err != nil {
		log.FromContext(ctx).With("error", err).Error("failed to create text geo search metric tag for:", source)
	}
	stats.Record(ctx, countMeasure.M(0))
	stats.Record(ctx, latencyMeasure.M(latency.Milliseconds()))
}

// helpers

func contextWithTags(ctx context.Context, mType, source, status string) (context.Context, error) {
	return tag.New(ctx, []tag.Mutator{
		tag.Insert(typeTag, mType),
		tag.Insert(sourceTag, source),
		tag.Insert(statusTag, status),
	}...)
}
