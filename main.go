package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type Album struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Year   int    `json:"year"`
	Artist string `json:"artist"`
	Image  []byte `json:"image,omitempty"`
}

func initDB() {
	var err error
	// Modify with your RDS MySQL database connection information
	DB_DSN := "admin:gjl990110@tcp(demo-mysql.coje1ucxmgno.us-west-2.rds.amazonaws.com:3306)/album_store"
	db, err = sql.Open("mysql", DB_DSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL!")

	// Create the albums table if it doesn't exist
	createTable()
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS albums (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		year INT NOT NULL,
		artist VARCHAR(255) NOT NULL,
		image LONGBLOB
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create albums table: %v", err)
	} else {
		fmt.Println("Albums table created or already exists.")
	}
}

// Handle POST request to save album
func uploadAlbum(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // Limit upload size

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid image upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read image data
	imageData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}

	// Read JSON data
	title := r.FormValue("title")
	year := r.FormValue("year")
	artist := r.FormValue("artist")

	// Store in DB
	_, err = db.Exec("INSERT INTO albums (title, year, artist, image) VALUES (?, ?, ?, ?)", title, year, artist, imageData)
	if err != nil {
		http.Error(w, "Failed to insert album", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Album uploaded successfully")
}

// Handle GET request to retrieve album by ID
func getAlbum(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	row := db.QueryRow("SELECT id, title, year, artist, image FROM albums WHERE id = ?", id)
	var album Album
	err := row.Scan(&album.ID, &album.Title, &album.Year, &album.Artist, &album.Image)
	if err != nil {
		http.Error(w, "Album not found", http.StatusNotFound)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

func main() {
	initDB()
	defer db.Close() // Ensure database connection is closed at the end

	router := mux.NewRouter()
	router.HandleFunc("/album", uploadAlbum).Methods("POST")
	router.HandleFunc("/album/{id}", getAlbum).Methods("GET")

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
