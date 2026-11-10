package http

import (
	"encoding/json"
	"errors"
	"kaspersky/internal/threat"
	"time"
)

type ThreatDTO struct {
	MD5         string `json:"md5"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type UpdateThreatDTO struct {
	NewName        string `json:"name"`
	NewDescription string `json:"description"`
	NewSeverity    string `json:"severity"`
}

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func (t *ThreatDTO) ValidateForCreate() error {
	if t.Name == "" {
		return errors.New("Имя угрозы не задано")
	}

	if t.Severity == "" || !threat.IsSeverityAllowed(t.Severity) {
		return errors.New("Уровень опасности не задан или недопустим")
	}

	if !threat.IsMD5Allowed(t.MD5) {
		return errors.New("Недопустимое значение md5")
	}

	return nil
}

func (e *ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		panic(e)
	}

	return string(b)
}
