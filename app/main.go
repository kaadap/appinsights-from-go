package main

import (
	"fmt"
	"os"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func main() {

	logger := MakeAppInsightsLogger()

	appinsights.NewDiagnosticsMessageListener(func(msg string) error {
		fmt.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
		return nil
	})

	logger.logMetricsToAppInsightsPeriodically()
}

type AppInsightsLogger struct {
	client appinsights.TelemetryClient
}

func MakeAppInsightsLogger() *AppInsightsLogger {

	ikey := os.Getenv("APP_INSIGHTS_INSTRUMENTATION_KEY")
	telemetryConfig := appinsights.NewTelemetryConfiguration(ikey)

	endpoint := os.Getenv("APP_INSIGHTS_INGESTION_ENDPOINT")
	telemetryConfig.EndpointUrl = endpoint

	// Configure how many items can be sent in one call to the data collector:
	telemetryConfig.MaxBatchSize = 8192

	// Configure the maximum delay before sending queued telemetry:
	telemetryConfig.MaxBatchInterval = 2 * time.Second

	logger := &AppInsightsLogger{
		client: appinsights.NewTelemetryClientFromConfig(telemetryConfig),
	}

	return logger
}

func (logger *AppInsightsLogger) logMetrics(ch chan struct{}) {

	logger.client.TrackMetric("Workers.Desired", 15)
	logger.client.TrackMetric("Workers.Available", 8)
	logger.client.TrackMetric("Workers.InUse", 342)
	logger.client.TrackMetric("Workers.Initializing", 9)

	logger.client.TrackTrace("Test trace message", appinsights.Information)

	logger.client.TrackEvent("Test event successfully sent!")
	logger.client.Channel().Flush()
	time.Sleep(1 * time.Second)
	ch <- struct{}{}
}

func (logger *AppInsightsLogger) logMetricsToAppInsightsPeriodically() {
	wait := make(chan struct{})
	for {
		go logger.logMetrics(wait)
		<-wait
	}
}
