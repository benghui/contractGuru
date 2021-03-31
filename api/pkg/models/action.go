package models

type Action struct {
	ActionName    string `json:"action_name"`
	FinanceFlag int `json:"finance_flag"`
	TransitionID int `json:"transition_id"`
	EndStateID int `json:"end_state_id"`
}
