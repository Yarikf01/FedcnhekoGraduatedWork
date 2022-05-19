package log_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/Yarikf01/graduatedwork/api/utils"
)

func TestLogWithContext(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, log.L, log.FromContext(ctx))
	newLog := log.L.With("field-1", 12)
	ctx = log.WithLogger(ctx, newLog)

	assert.Equal(t, newLog, log.FromContext(ctx))
}
