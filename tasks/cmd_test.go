package tasks

import (
	"context"
	"testing"
)

func TestSyncCommand(t *testing.T) {
	SyncCommand(context.TODO(), 1, "python -u ../haha.py", false)
}