package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "crud/lib/database"
	utils "crud/utils"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Database connection setup
var db, _ = sql.Open("sqlite3", "./lib/database/db.sqlite")

func Routes() {
	http.HandleFunc("/api/msg/", RequestType)

	// log.Println("/ Listening on port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func RequestType(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		HandlerGet(w, r)
	case http.MethodPost:
		HandlerPost(w, r)
	case http.MethodDelete:
		HandlerDelete(w, r)
	case http.MethodPatch:
		HandlerUpdate(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler function (GET)
func HandlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := strings.TrimPrefix(r.URL.Path, "/api/msg/")

		if id == "" {
			ip := r.RemoteAddr
			fmt.Println("User IP:", ip, "Requested All Messages")

			rows, err := db.Query("SELECT id, created, username, msg FROM messages")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var responses []utils.Data

			for rows.Next() {
				var id, created, username, msg string
				err := rows.Scan(&id, &created, &username, &msg)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				response := utils.Data{
					Id:       id,
					Created:  created,
					Username: username,
					Msg:      msg,
				}

				responses = append(responses, response)
			}

			responseJSON, err := json.Marshal(responses)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseJSON)
		} else {
			ip := r.RemoteAddr
			fmt.Println("User IP:", ip, "Requested Message with ID:", id)

			row := db.QueryRow("SELECT created, username, msg FROM messages WHERE id = ?", id)

			var created, username, msg string

			err := row.Scan(&created, &username, &msg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			response := utils.Data{
				Id:       id,
				Created:  created,
				Username: username,
				Msg:      msg,
			}

			responseJSON, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseJSON)
		}
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

// Handler function (POST)
func HandlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data utils.Data
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if data.Username == "" || data.Msg == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		id := uuid.New().String()
		data.Id = id

		fmt.Println("User", data.Username, "Sent Message:", data.Msg)

		currentTime := time.Now().Format("2006-01-02 15:04:05")

		_, err = db.Exec("INSERT INTO messages (id, created, username, msg) VALUES (?, ?, ?, ?)", data.Id, currentTime, data.Username, data.Msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

// Handler function (DELETE)
func HandlerDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		var data utils.Data
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Println(w, err.Error(), http.StatusBadRequest)
			return
		}

		if data.Id == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		fmt.Println("Message with id:", data.Id, "got deleted")

		_, err = db.Exec("DELETE FROM messages WHERE id = ?", data.Id)
		if err != nil {
			log.Println(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func HandlerUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PATCH" {
		var data utils.Data
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if data.Id == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		fmt.Println(data.Id, "got uppdated with message:", data.Msg)

		_, err = db.Exec("UPDATE messages SET msg=? WHERE id=?", data.Msg, data.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if data.Id == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}
}
