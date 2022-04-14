package utils_test

import (
	"context"
	"github.com/realHoangHai/awesome/pkg/utils"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestCorrelationContext(t *testing.T) {
	expID := "123"
	// correlation id from x-correlation-id
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(utils.XCorrelationID, expID))
	id, ok := utils.CorrelationIDFromContext(ctx)
	if !ok || id != expID {
		t.Errorf("got correlation_id=%s, want correlation_id=%s", id, expID)
	}

	// correlation id from x-request-id
	ctx = metadata.NewIncomingContext(context.Background(), metadata.Pairs(utils.XRequestID, expID))
	id, ok = utils.CorrelationIDFromContext(ctx)
	if !ok || id != expID {
		t.Errorf("got correlation_id=%s, want correlation_id=%s", id, expID)
	}

	// generate new correlation id if not existed.
	ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{})
	id, ok = utils.CorrelationIDFromContext(ctx)
	if ok || id == "" {
		t.Errorf("got correlation_id=%s, want correlation_id=%s", id, expID)
	}
}
