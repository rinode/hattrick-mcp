package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lucianoq/hattrick/chpp"
	"github.com/lucianoq/hattrick/chpp/id"
)

func TestPlain(t *testing.T) {
	p := &chpp.Player{
		ID:          id.Player(123),
		FirstName:   "Jan",
		Form:        chpp.SkillLevelExcellent,
		Specialty:   chpp.SpecialtyTechnical,
		Salary:      chpp.Money(25000),
		ArrivalDate: chpp.FromTime(time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)),
	}
	b, err := json.Marshal(plain([]*chpp.Player{p, nil}))
	if err != nil {
		t.Fatal(err)
	}
	var got []map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	want := map[string]any{
		"ID":        float64(123),
		"FirstName": "Jan",
		"Form":      "excellent (8)",
		"Specialty": "Technical (1)",
		"Salary":    float64(25000),
	}
	for k, w := range want {
		if got[0][k] != w {
			t.Errorf("%s = %v, want %v", k, got[0][k], w)
		}
	}
	if s, _ := got[0]["ArrivalDate"].(string); s == "" {
		t.Errorf("ArrivalDate = %v, want a date string", got[0]["ArrivalDate"])
	}
	if _, ok := got[0]["LastName"]; ok {
		t.Error("empty LastName not dropped")
	}
	if got[1] != nil {
		t.Errorf("nil player = %v, want null", got[1])
	}
}
