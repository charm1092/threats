package threat

import "testing"

func TestIsSeverityAllowed(t *testing.T)  {
	if !IsSeverityAllowed("high") {
        t.Error("уровень high должен быть допустимым")
    }

	if IsSeverityAllowed("bread") {
		t.Error("уровень banana должен быть отклонён")
	}
}