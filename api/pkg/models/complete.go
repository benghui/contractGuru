package models

import "time"

type Completed struct {
	CompletedID             int       `json:"completed_id"`
	RequestID               int       `json:"request_id"`
	CounterpartyName        string    `json:"counterparty_name"`
	CounterpartyInformation string    `json:"counterparty_information"`
	ContractValue           float64   `json:"contract_value"`
	ContractType            string    `json:"contract_type"`
	Region                  string    `json:"region"`
	EffectiveDate           time.Time `json:"effective_date"`
	TerminationDate         time.Time `json:"termination_date"`
	RenewalDate             time.Time `json:"renewal_date"`
	Purpose                 string    `json:"purpose"`
	Background              string    `json:"background"`
}
