package controllers

import (
	"log"
	"gotrading/config"
	"gotrading/app/models"
	"gotrading/bitflyer"
)

func StreamIngestionData() {
	var tickerChannel = make(chan bitflyer.Ticker)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	go apiClient.GetRealTimeTicker(config.Config.ProductCode, tickerChannel)
	go func() {
		for ticker := range tickerChannel {
			log.Printf("action=StreamIngestionData, %v", ticker)
			for _, duration := range config.Config.Durations {
				isCreate := models.CreateCandleWithDuration(ticker, ticker.ProductCode, duration)
				if isCreate == true && duration == config.Config.TradeDuration {

				}
			}
		}
	}()
}
