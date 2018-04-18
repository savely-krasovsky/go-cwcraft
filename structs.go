package main

type (
	item struct {
		ID        string         `json:"id"`
		Name      string         `json:"name"`
		Stats     stats          `json:"stats"`
		Type      string         `json:"type"`
		ManaCost  int            `json:"mana_cost,omitempty"`
		Composite bool           `json:"composite"`
		Recipe    map[string]int `json:"recipe,omitempty"`
	}

	stats struct {
		Attack  int `json:"attack,omitempty"`
		Defense int `json:"defense,omitempty"`
		Mana    int `json:"mana,omitempty"`
	}

	resource struct {
		ID        string         `json:"id"`
		Name      string         `json:"name"`
		ManaCost  int            `json:"mana_cost,omitempty"`
		Composite bool           `json:"composite"`
		Recipe    map[string]int `json:"recipe,omitempty"`
	}

	command struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Amount          int    `json:"amount"`
		CommandManaCost int    `json:"command_mana_cost"`
	}

	basic struct {
		Name   string `json:"name"`
		Amount int    `json:"amount"`
	}
)
