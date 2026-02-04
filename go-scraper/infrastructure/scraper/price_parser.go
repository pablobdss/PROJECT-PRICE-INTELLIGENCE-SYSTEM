package scraper

import (
	"regexp"
	"strconv"
	"strings"
)


var currencyMap = map[string]struct {
	Code      string
	Separator string
}{
	"£":  {Code: "GBP", Separator: "."},
	"$":  {Code: "USD", Separator: "."},
	"R$": {Code: "BRL", Separator: ","},
	"€":  {Code: "EUR", Separator: ","},
}

var numberRegex = regexp.MustCompile(`[0-9.,]+`)


func parsePrice(raw string) (float64, string) {
	raw = strings.TrimSpace(raw)

	detectedCurrency := "BRL"
	decimalSeparator := ","


	for symbol, config := range currencyMap {
		if strings.Contains(raw, symbol) {
			detectedCurrency = config.Code
			decimalSeparator = config.Separator
			break
		}
	}

	
	cleanNum := numberRegex.FindString(raw)

	
	if decimalSeparator == "," {
		cleanNum = strings.ReplaceAll(cleanNum, ".", "")  
		cleanNum = strings.ReplaceAll(cleanNum, ",", ".") 
	}

	if decimalSeparator == "." {
		cleanNum = strings.ReplaceAll(cleanNum, ",", "") 
	} 
	
	val, err := strconv.ParseFloat(cleanNum, 64)
	
	if err != nil {return 0, detectedCurrency}

	return val, detectedCurrency
}
