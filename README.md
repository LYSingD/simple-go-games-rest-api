# Simple rest api with Golang

A simple REST API for self-practice. 
Inspired by kubucation YouTube channel ([link](https://www.youtube.com/watch?v=2v11Ym6Ct9Q&t=573s&ab_channel=kubucation))


## Tech Stack
* Golang
* Curl command

  
## Data type
There is only one object called `game` contain the following default data:
```
type Game struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Developer string   `json:"developer"`
	Rating    string   `json:"rating"`
	Genres    []string `json:"genres"`
}

[
    {
	ID:        "816d17bb-c943-4b8c-ba8a-54e0429985c7",
	Name:      "Ghost of Tsushima",
	Developer: "Sucker Punch",
	Rating:    "M",
	Genres:    []string{"General", "Action Adventure", "Open-World"},
    },
    {
	ID:        "ee9015a7-6219-4620-af98-cf78601c6446",
	Name:      "Monster Hunter World: Iceborne",
	Developer: "Capcom",
	Rating:    "T",
	Genres:    []string{"Action"},
    },
    {
	ID:        "c2b337b0-839f-40b8-b43d-921bfd2812a8",
	Name:      "Watch Dog",
	Developer: "Ubisoft",
	Rating:    "M",
	Genres:    []string{"Action", "Adventure"},
    }
]
```

  
## Run Locally

Clone the project

```bash
  git clone https://github.com/LYSingD/go-games-rest-api.git
```

Go to the project directory

```bash
  cd [path]/go-games-rest-api/
```

Configure the desired port in `/server.go`, line 14.
```go
log.Fatal(http.ListenAndServe(":[port]", nil))
```

Run the server

```bash
  go run server.go
```

Use `curl` command, like
```bash
  curl localhost:8080/games
```

  
## Features

- `GET /games` returns a list with all games
- `POST /games/` adds a new game
- `GET /games/:id` return a JSON object of a specific game 
- `PUT /games/:id` updates a specific game
- `DELETE /games/:id` removes a specific game

  
