package models

import (
	"html"
	"strings"
	"time"
)

type Request struct {
	RequestID      int       `json:"request_id"`
	RequesterID    int       `json:"requester_id"`
	BuID           int       `json:"bu_id"`
	CurrentStateID int       `json:"current_state_id"`
	RequestName    string    `json:"request_name"`
	CreatedAt      time.Time `json:"created_at"`
	FinanceFlag    int       `json:"finance_flag"`
}

func (r *Request) Prepare() {
	r.RequestName = html.EscapeString(strings.TrimSpace(r.RequestName))
}
