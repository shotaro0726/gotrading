package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"regexp"

	"gotrading/config"
	"gotrading/models"
	"github.com/labstack/gommon/log"
)

var templates = template.Must(template.ParseFiles("app/views/chart.html"))

func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	// limit := 100
	// duration := "1s"
	// durationTime := config.Config.Duration[duration]
	// df, _ := models.GetAllCandle(config.Config.ProductCode, durationTime, limit)
	err := templates.ExecuteTemplate(w, "chart.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type JsonError struct {
	Error string `json:"error"`
	Code int `json:"code"`
}

func ApiError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JsonError{Error: errMessage, Code: code})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonError)
}

var apiValidPath = regexp.MustCompile("^/api/candle/$")

func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := apiValidPath.FindStringSubmatch(r.URL.Path)
		if len(m) == 0 {
			ApiError(w, "Not Found", http.StatusNotFound)
		}
		fn(w,r)
	}
}

func apiCandleHandler(w http.ResponseWriter, r *http.Request) {
	productCode := r.URL.Query().Get("product_code")
	if productCode	== "" {
		ApiError(w, "No product_code_param", http.StatusBadRequest)
		return 
	}
	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if strLimit == "" || err != nil || limit < 0 || limit > 1000 {
		limit = 1000
	}

	duration := r.URL.Query().Get("duration")
	if duration == "" {
		duration = "1m"
	} 
	durationTime := config.Config.Durations[duration]
	df, _ := models.GetAllCandle(productCode, durationTime, limit)
	sma := r.URL.Query().Get("sma")

	if sma != "" {
		strSmaPeriod1 := r.URL.Query().Get("smaPeriod1")
		strSmaPeriod2 := r.URL.Query().Get("smaPeriod2")
		strSmaPeriod3 := r.URL.Query().Get("smaPeriod3")

		period1, err := strconv.Atoi(strSmaPeriod1)
		if strSmaPeriod1 == "" || err != nil || period1 < 0 {
			period1 = 7
		}

		period2, err := strconv.Atoi(strSmaPeriod1)
		if strSmaPeriod2 == "" || err != nil || period2 < 0 {
			period2 = 14
		}

		period3, err := strconv.Atoi(strSmaPeriod1)
		if strSmaPeriod3 == "" || err != nil || period3 < 0 {
			period3 = 50
		}
		df.AddSma(period1)
		df.AddSma(period2)
		df.AddSma(period3)
	}

	bbands := r.URL.Query().Get("bbands")
	if bbands != "" {
		strN := r.URL.Query().Get("bbandsN")
		strK := r.URL.Query().Get("bbandsK")
		n, err := strconv.Atoi(strN)
		if strN == "" || err != nil || n < 0 {
			n = 20
		}
		k, err := strconv.Atoi(strK)
		if strK == "" || err != nil || k < 0 {
			k = 2
		}
		df.AddBBands(n, float64(k))
	}

	js, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func StartWebServer() error {
	http.HandleFunc("api/candle/", apiMakeHandler(apiCandleHandler))
	http.HandleFunc("/chart/", viewChartHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}