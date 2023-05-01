package models

import "skfw/papaya/pigeon/templates/basicAuth/models"

// requirement for basicAuth

type Session struct {
	*models.SessionModel
}

func (Session) TableName() string {

	return "sessions"
}
