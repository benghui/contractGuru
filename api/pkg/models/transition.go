package models

type RequestTransition struct {
	TransitionID   int `json:"transition_id"`
	EndStateID int `json:"current_state_id"`
}
