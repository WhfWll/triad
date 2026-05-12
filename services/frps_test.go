// Package services
// @Author bcy2007  2024/12/13 10:25
package services

import (
	"context"
	"testing"
)

func TestFindRunningFrpsService(t *testing.T) {
	srv := FrpsService{}
	ctx := context.Background()
	status, err := srv.FindRunningFrpsService(ctx)
	if err != nil {
		t.Log(err)
	}
	t.Log(status)
}
