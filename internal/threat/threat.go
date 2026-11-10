package threat

import "strings"

type Threat struct {
	MD5         string `json:"md5"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

func NewThreat(md5, name, description, severity string) *Threat {
	return &Threat{
		MD5:         md5,
		Name:        name,
		Description: description,
		Severity:    severity,
	}
}

const allowedMD5 = "0123456789abcdefABCDEF"
var allowedSeverity = []string{"safe", "low", "medium", "high"}

func IsMD5Allowed(md5 string) bool {
	if len(md5) != 32 || md5 == "" {
		return false
	}

	for _, c := range md5 {
		if !strings.ContainsRune(allowedMD5, c) {
			return false
		}
	}

	return true
}

func IsSeverityAllowed(severity string) bool {
	for _, v := range allowedSeverity {
		if severity == v {
			return true
		}
	} 

	return false
}