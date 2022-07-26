// Package csv contains parsing methods
package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mkokoulin/exchanges-history-app/internal/models"
)

const (
	DateCol           int = 0
	CryptoamountCol   int = 2
	FiatamountCol     int = 3
	FeeCol            int = 4
	CryptocurrencyCol int = 7
	PaymethodCol      int = 8
	TypeCol           int = 11
	StatusCol         int = 12
)

func parser(data [][]string, filters ...func(eh models.ExchangesHistory) bool) ([]models.ExchangesHistory, error) {
	var exchangesHistory []models.ExchangesHistory

	for i, line := range data {
		if i > 0 {
			var rec models.ExchangesHistory
			for col, field := range line {
				trimmed := strings.TrimSpace(field)
				switch col {
				case DateCol:
					t, err := time.Parse("02-01-2006 15:04:05", trimmed)
					if err != nil {
						fmt.Println(err)
					}

					rec.Date = t
				case CryptoamountCol:
					fl, err := strconv.ParseFloat(trimmed, 32)
					if err != nil {
						return nil, err
					}
					rec.Cryptoamount = fl
				case FiatamountCol:
					fl, err := strconv.ParseFloat(trimmed, 32)
					if err != nil {
						return nil, err
					}
					rec.Fiatamount = fl
				case FeeCol:
					fl, err := strconv.ParseFloat(trimmed, 32)
					if err != nil {
						return nil, err
					}
					rec.Fee = fl
				case CryptocurrencyCol:
					rec.Cryptocurrency = trimmed
				case PaymethodCol:
					rec.Paymethod = trimmed
				case TypeCol:
					rec.Type = trimmed
				case StatusCol:
					rec.Status = trimmed
				default:
					log.Printf("The column %d does not fit the conditions of the function", col)
				}
			}

			if combineFilters(rec, filters...) {
				exchangesHistory = append(exchangesHistory, rec)
			}
		}
	}

	return exchangesHistory, nil
}

func Reader(reader io.Reader) ([]models.ExchangesHistory, error) {
	csvReader := csv.NewReader(reader)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	exchangesHistory, err := parser(data, filterType, filterStatus)
	if err != nil {
		return nil, err
	}

	return exchangesHistory, nil
}
