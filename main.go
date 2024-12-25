package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type AlbumsResponseData struct {
	Albums []Album `json:"albums"`
}

type AlbumsResponse struct {
	Data AlbumsResponseData `json:"data"`
	Ok   bool               `json:"ok"`
}

type AlbumResponseData struct {
	Album Album `json:"album"`
}

type AlbumResponse struct {
	Data AlbumResponseData `json:"data"`
	Ok   bool              `json:"ok"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func main() {

	http.HandleFunc("/albums/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		if len(parts) != 3 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Not Found"})
			return
		}

		if r.Method != http.MethodGet && r.Method != http.MethodPut && r.Method != http.MethodDelete {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Method Not Allowed"})
			return
		}

		content, err := os.ReadFile("albums.json")

		if err != nil {
			log.Fatal("Error when opening albums json file: ", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
			return
		}

		var albums []Album
		err = json.Unmarshal(content, &albums)

		if err != nil {
			log.Fatal("Error parsing albums json payload: ", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
			return
		}

		var id string = parts[2]

		if r.Method == http.MethodGet {

			var album Album
			for _, a := range albums {
				if a.ID == id {
					album = a
					break
				}
			}

			if album.ID == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Album Not Found"})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(AlbumResponse{Ok: true, Data: AlbumResponseData{Album: album}})

		} else if r.Method == http.MethodPut {

			var album Album
			err := json.NewDecoder(r.Body).Decode(&album)

			if err != nil {
				log.Fatal("Invalid payload for album", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Bad Request"})
				return
			}

			var albumIndex int = -1
			for i, a := range albums {
				if a.ID == id {
					albumIndex = i
					break
				}
			}

			if albumIndex == -1 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Album Not Found"})
				return
			}

			if id != album.ID {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Album Not Found"})
				return
			}

			albums[albumIndex] = album

			file, err := os.OpenFile("albums.json", os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				log.Fatal("Error when opening albums json file: ", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
				return
			}

			err = json.NewEncoder(file).Encode(albums)

			if err != nil {
				log.Fatal("Error writing to albums json file: ", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(AlbumResponse{Ok: true, Data: AlbumResponseData{Album: album}})
		} else if r.Method == http.MethodDelete {

			var albumIndex int = -1
			for i, a := range albums {
				if a.ID == id {
					albumIndex = i
					break
				}
			}

			if albumIndex == -1 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Album Not Found"})
				return
			}

			albums = append(albums[:albumIndex], albums[albumIndex+1:]...)

			file, err := os.OpenFile("albums.json", os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				log.Fatal("Error when opening albums json file: ", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
				return
			}

			err = json.NewEncoder(file).Encode(albums)

			if err != nil {
				log.Fatal("Error writing to albums json file: ", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(AlbumsResponse{Ok: true, Data: AlbumsResponseData{Albums: albums}})
		}

	})

	http.HandleFunc("/albums", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Method Not Allowed"})
			return
		}

		if r.Method == http.MethodGet {
			content, err := os.ReadFile("albums.json")

			if err != nil {
				log.Fatal("Error when opening albums json file: ", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
				return
			}

			var albums []Album
			err = json.Unmarshal(content, &albums)

			if err != nil {
				log.Fatal("Error parsing albums json payload: ", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(AlbumsResponse{Ok: true, Data: AlbumsResponseData{Albums: albums}})

		} else if r.Method == http.MethodPost {

			var album Album
			err := json.NewDecoder(r.Body).Decode(&album)

			if err != nil {
				log.Fatal("Invalid payload for album", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Bad Request"})
				return
			}

			content, err := os.ReadFile("albums.json")

			if err != nil {
				log.Fatal("Error when opening albums json file: ", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
				return
			}

			var albums []Album
			err = json.Unmarshal(content, &albums)

			if err != nil {
				log.Fatal("Error parsing albums json payload: ", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
				return
			}

			for _, a := range albums {
				if a.ID == album.ID {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusConflict)
					json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Conflict - Album already exists"})
					return
				}
			}

			albums = append(albums, album)

			file, err := os.OpenFile("albums.json", os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				log.Fatal("Error when opening albums json file: ", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
				return
			}

			err = json.NewEncoder(file).Encode(albums)

			if err != nil {
				log.Fatal("Error writing to albums json file: ", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Ok: false, Message: "Internal Server Error"})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(AlbumsResponse{Ok: true, Data: AlbumsResponseData{Albums: albums}})

		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")

	})

	fmt.Println("Server is running on port http://locahost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
