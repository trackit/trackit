package tagging

import (
	"context"

	"github.com/trackit/jsonlog"
)

// UpdateMostUsedTagsForAccount updates most used tags in MySQL for the specified AWS account
func UpdateMostUsedTagsForAccount(ctx context.Context, account int, awsAccount string) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	logger.Debug("UpdateMostUsedTagsForAccount", nil)

	return nil
}
