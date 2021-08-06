package games

type Game struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Developer string   `json:"developer"`
	Rating    string   `json:"rating"`
	Genres    []string `json:"genres"`
}
