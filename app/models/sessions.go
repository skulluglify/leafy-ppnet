package models

import "skfw/papaya/pigeon/templates/basicAuth/models"

// requirement for basicAuth

type Sessions struct {
	*models.SessionModel
}

func (Sessions) TableName() string {

	return "sessions"
}
