package pokeapi

import (
	"encoding/json"
	"testing"
)

func TestPokemonUnmarshalJSON(t *testing.T) {
	raw := `{
		"name": "pikachu",
		"base_experience": 112,
		"height": 4,
		"weight": 60,
		"stats": [
			{"base_stat": 35, "stat": {"name": "hp"}},
			{"base_stat": 55, "stat": {"name": "attack"}},
			{"base_stat": 90, "stat": {"name": "speed"}}
		],
		"types": [
			{"type": {"name": "electric"}}
		]
	}`

	var p Pokemon
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if p.Name != "pikachu" {
		t.Errorf("Name = %q, want %q", p.Name, "pikachu")
	}
	if p.BaseExperience != 112 {
		t.Errorf("BaseExperience = %d, want 112", p.BaseExperience)
	}
	if p.Height != 4 {
		t.Errorf("Height = %d, want 4", p.Height)
	}
	if p.Weight != 60 {
		t.Errorf("Weight = %d, want 60", p.Weight)
	}

	if len(p.Stats) != 3 {
		t.Fatalf("len(Stats) = %d, want 3", len(p.Stats))
	}
	wantStats := []Stat{
		{Name: "hp", BaseStat: 35},
		{Name: "attack", BaseStat: 55},
		{Name: "speed", BaseStat: 90},
	}
	for i, want := range wantStats {
		if p.Stats[i] != want {
			t.Errorf("Stats[%d] = %+v, want %+v", i, p.Stats[i], want)
		}
	}

	if len(p.Types) != 1 || p.Types[0] != "electric" {
		t.Errorf("Types = %v, want [electric]", p.Types)
	}
}

func TestPokemonUnmarshalJSONMultipleTypes(t *testing.T) {
	raw := `{
		"name": "charizard",
		"base_experience": 267,
		"height": 17,
		"weight": 905,
		"stats": [],
		"types": [
			{"type": {"name": "fire"}},
			{"type": {"name": "flying"}}
		]
	}`

	var p Pokemon
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(p.Types) != 2 {
		t.Fatalf("len(Types) = %d, want 2", len(p.Types))
	}
	if p.Types[0] != "fire" || p.Types[1] != "flying" {
		t.Errorf("Types = %v, want [fire flying]", p.Types)
	}
	if len(p.Stats) != 0 {
		t.Errorf("len(Stats) = %d, want 0", len(p.Stats))
	}
}

func TestPokemonUnmarshalJSONEmptyArrays(t *testing.T) {
	raw := `{
		"name": "missingno",
		"base_experience": 0,
		"height": 0,
		"weight": 0,
		"stats": [],
		"types": []
	}`

	var p Pokemon
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(p.Stats) != 0 {
		t.Errorf("len(Stats) = %d, want 0", len(p.Stats))
	}
	if len(p.Types) != 0 {
		t.Errorf("len(Types) = %d, want 0", len(p.Types))
	}
}
