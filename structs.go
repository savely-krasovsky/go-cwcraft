package main

type (
	equipmentItem struct {
		ID        string         `json:"id"`
		Name      string         `json:"name"`
		Stats     stats          `json:"stats"`
		Type      string         `json:"type"`
		ManaCost  int            `json:"mana_cost,omitempty"`
		Composite bool           `json:"composite"`
		Recipe    map[string]int `json:"recipe,omitempty"`
	}

	alchemyItem struct {
		ID        string         `json:"id"`
		Name      string         `json:"name"`
		Effect    string         `json:"effect"`
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
		Name       string `json:"name"`
		Amount     int    `json:"amount"`
		UserAmount int    `json:"user_amount,omitempty"`
	}

	login struct {
		Status string `json:"status"`
		ID     int    `json:"id" form:"id" query:"id"`
		Code   string `json:"code" form:"code" query:"code"`
	}

	user struct {
		ID    string         `json:"_key"`
		Token string         `json:"token,omitempty"`
		Stock map[string]int `json:"stock,omitempty"`
	}
)
