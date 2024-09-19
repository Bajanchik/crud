package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "password"
	dbName   = "crud"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	if createUser(db, user.Name, user.Email) != nil {
		http.Error(w, "Failed to crete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User created successfully")
	log.Println("User created successfully")
}

func createUser(db *sql.DB, name, email string) error {
	query := "INSERT INTO users (name, email) VALUES (?, ?)"
	_, err := db.Exec(query, name, email)
	if err != nil {
		return err
	}
	return nil
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Convertation error", http.StatusNotFound)
	}

	user, err := GetUser(db, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.Println("GET user successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUser(db *sql.DB, id int) (*User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	row := db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Convertation error", http.StatusNotFound)
	}

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	if UpdateUser(db, userID, user.Name, user.Email) != nil {
		http.Error(w, "User not found", http.StatusNotFound)
	}
	log.Println("User updated successfully")
	fmt.Fprintln(w, "User updated successfully")
}

func UpdateUser(db *sql.DB, id int, name, email string) error {
	query := "UPDATE users SET name = ?, email = ? WHERE id = ?"
	_, err := db.Exec(query, name, email, id)
	if err != nil {
		return err
	}
	return nil
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Convertation error", http.StatusNotFound)
	}

	if DeleteUser(db, userID) != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	log.Println("User deleted successfully")
	fmt.Fprintln(w, "User deleted successfully")

}

func DeleteUser(db *sql.DB, id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", helloHandler)
	r.HandleFunc("/user", createUserHandler).Methods("POST")
	r.HandleFunc("/user{id}", getUserHandler).Methods("GET")
	r.HandleFunc("/user{id}", updateUserHandler).Methods("PUT")
	r.HandleFunc("/user{id}", deleteUserHandler).Methods("DELETE")

	log.Println("Server listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
