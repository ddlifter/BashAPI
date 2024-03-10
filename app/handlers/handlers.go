package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	container "github.com/ddlifter/BashAPI/app/services"
	database "github.com/ddlifter/BashAPI/db"
	"github.com/gorilla/mux"
)

func RunCommand(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	var command string
	err := db.QueryRow("SELECT Script FROM Commands WHERE id = $1", id).Scan(&command)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = db.Exec("UPDATE Commands SET Status = $1 WHERE id = $2", "running", id)
	if err != nil {
		log.Fatal(err)
	}

	resultChan := make(chan string)

	go func() {
		output, err := container.CreateAndRunDockerContainer(command)
		if err != nil {
			log.Fatal(err)
		}
		resultChan <- output

		_, err = db.Exec("UPDATE Commands SET Status = $1 WHERE id = $2", output, id)
		if err != nil {
			log.Fatal(err)
		}
	}()

	output := <-resultChan

	fmt.Fprintf(w, "Command output: %s", output)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func StopCommand(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	var status string
	err := db.QueryRow("SELECT Status FROM Commands WHERE id = $1", id).Scan(&status)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if status != "running" {
		w.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(w, "Command with ID %s is not running, cannot be stopped", id)
		return
	}

	_, err = db.Exec("UPDATE Commands SET Status = $1 WHERE id = $2", "stopped", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while stopping command with ID %s", id)
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Command with ID %s has been successfully stopped", id)
}

func GetCommands(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM Commands")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	commands := map[int]database.Command{}
	for rows.Next() {
		var u database.Command
		if err := rows.Scan(&u.ID, &u.Name, &u.Status, &u.Script); err != nil {
			log.Fatal(err)
		}
		commands[u.ID] = u
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(commands)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func GetCommand(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	var u database.Command
	err := db.QueryRow("SELECT * FROM Commands WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Status, &u.Script)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func AddCommand(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	var u database.Command
	json.NewDecoder(r.Body).Decode(&u)

	res, err := db.Exec("INSERT INTO Commands (Name, Status, Script) VALUES ($1, $2, $3)", u.Name, "waiting", u.Script)
	if err != nil {
		log.Fatal(err)
	}

	id, _ := res.LastInsertId()
	u.ID = int(id)

	json.NewEncoder(w).Encode(u)
}

func DeleteCommand(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	var u database.Command
	err := db.QueryRow("SELECT * FROM Commands WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Status, &u.Script)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		_, err := db.Exec("DELETE FROM Commands WHERE id = $1", id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode("Command deleted")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteCommands(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	rows, err := db.Query("SELECT ID FROM Commands")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var u database.Command
		if err := rows.Scan(&u.ID); err != nil {
			log.Fatal(err)
		}
		_, err := db.Exec("DELETE FROM Commands WHERE id = $1", u.ID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
