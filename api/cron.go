package api

import (
	"context"

	"github.com/sunnyegg/go-so/cron"
)

func (server *Server) registerCron() {
	cronClient := cron.NewCron()
	cronClient.AddFunc("@hourly", cron.ValidateToken(context.Background(), server.store, server.config))
	cronClient.AddFunc("0 0 */2 * * *", cron.DeleteExpiredSession(context.Background(), server.store, server.config))

	cronClient.Start()
}
