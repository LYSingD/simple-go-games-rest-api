package games

type Game struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Developer string   `json:"developer"`
	Rating    string   `json:"rating"`
	Genres    []string `json:"genres"`
}

var GameList = map[string]Game{
	"816d17bb-c943-4b8c-ba8a-54e0429985c7": {
		ID:        "816d17bb-c943-4b8c-ba8a-54e0429985c7",
		Name:      "Ghost of Tsushima",
		Developer: "Sucker Punch",
		Rating:    "M",
		Genres:    []string{"General", "Action Adventure", "Open-World"},
	},
	"ee9015a7-6219-4620-af98-cf78601c6446": {
		Name:      "Monster Hunter World: Iceborne",
		Developer: "Capcom",
		Rating:    "T",
		Genres:    []string{"Action"},
	},
	"c2b337b0-839f-40b8-b43d-921bfd2812a8": {
		Name:      "Watch Dog",
		Developer: "Ubisoft",
		Rating:    "M",
		Genres:    []string{"Action", "Adventure"},
	},
}
