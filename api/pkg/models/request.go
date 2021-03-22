package models

import (
	"time"
)
type Request struct {
	RequestID    int       `json:"request_id"`
	RequesterID  int       `json:"requester_id"`
	BuID         int       `json:"bu_id"`
	CurrentState int       `json:"state_id"`
	RequestName  string    `json:"request_name"`
	RequestDate  time.Time `json:"request_date"`
	FinanceFlag  int       `json:"finance_flag"`
}
