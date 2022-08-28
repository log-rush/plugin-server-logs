package pluginServerLogs

import (
	"fmt"

	"github.com/log-rush/distribution-server/domain"
	"github.com/log-rush/distribution-server/pkg/app"
	"github.com/log-rush/distribution-server/pkg/devkit"
	logging "github.com/log-rush/go-client"
)

type Config struct {
	StreamName string
	Id         string
	Key        string
	BatchSize  int
}

type ServerLogsPlugin struct {
	config Config
	Plugin app.Plugin
	stream logging.Stream
}

func NewServerLogsPlugin(config Config) ServerLogsPlugin {
	plugin := ServerLogsPlugin{
		config: config,
	}

	p := devkit.NewPlugin(
		"server-logs",
		nil,
		nil,
		func(context *app.Context) domain.Logger {
			batchSize := config.BatchSize
			if batchSize < 20 {
				batchSize = 20
			}
			plugin.stream = logging.NewLogStream(logging.ClientOptions{
				DataSourceUrl: fmt.Sprintf("http://%s:%d/", context.Config.Host, context.Config.Port),
				BatchSize:     batchSize,
			}, config.StreamName, config.Id, config.Key)

			return devkit.NewLogger(plugin.HandleLog)
		},
	)
	plugin.Plugin = p

	return plugin
}

func (p *ServerLogsPlugin) HandleLog(level devkit.LogLevel, template string, args ...interface{}) {
	log := fmt.Sprintf("[server] [%s] %s", level, fmt.Sprintf(template, args...))
	p.stream.Log(log)
}
