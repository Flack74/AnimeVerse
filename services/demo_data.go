package services

import (
	"time"

	model "github.com/Flack74/mongoapi/models"
)

func SeedDemoUserData(userID string) error {
	demoAnimes := []model.Anime{
		{
			UserID:    userID,
			Name:      "Attack on Titan",
			Type:      "TV",
			Score:     9,
			Progress:  model.Progress{Watched: 87, Total: 87},
			Status:    "completed",
			Genre:     []string{"Action", "Drama", "Fantasy"},
			Notes:     "Epic story about humanity's fight against titans",
			Year:      2013,
			Season:    "Spring",
			ImageUrl:  "https://cdn.myanimelist.net/images/anime/10/47347.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID:    userID,
			Name:      "Demon Slayer",
			Type:      "TV",
			Score:     8,
			Progress:  model.Progress{Watched: 26, Total: 44},
			Status:    "watching",
			Genre:     []string{"Action", "Supernatural", "Historical"},
			Notes:     "Beautiful animation and compelling story",
			Year:      2019,
			Season:    "Spring",
			ImageUrl:  "https://cdn.myanimelist.net/images/anime/1286/99889.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID:    userID,
			Name:      "Your Name",
			Type:      "Movie",
			Score:     10,
			Progress:  model.Progress{Watched: 1, Total: 1},
			Status:    "completed",
			Genre:     []string{"Romance", "Drama", "Supernatural"},
			Notes:     "Masterpiece movie about body swapping",
			Year:      2016,
			Season:    "Fall",
			ImageUrl:  "https://cdn.myanimelist.net/images/anime/5/87048.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID:    userID,
			Name:      "Jujutsu Kaisen",
			Type:      "TV",
			Score:     0,
			Progress:  model.Progress{Watched: 0, Total: 24},
			Status:    "plan-to-watch",
			Genre:     []string{"Action", "School", "Supernatural"},
			Notes:     "Want to watch this popular series",
			Year:      2020,
			Season:    "Fall",
			ImageUrl:  "https://cdn.myanimelist.net/images/anime/1171/109222.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID:    userID,
			Name:      "Naruto",
			Type:      "TV",
			Score:     7,
			Progress:  model.Progress{Watched: 220, Total: 720},
			Status:    "on-hold",
			Genre:     []string{"Action", "Martial Arts", "Ninja"},
			Notes:     "Taking a break from fillers",
			Year:      2002,
			Season:    "Fall",
			ImageUrl:  "https://cdn.myanimelist.net/images/anime/13/17405.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	_, _, err := InsertMultipleAnimes(demoAnimes)
	return err
}