package main

import (
    "database/sql"
    "encoding/csv"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func loadUsersFromCSV(db *sql.DB, filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return err
    }

    // Start from the second record to skip the header
    for _, record := range records[1:] {
        _, err := db.Exec("INSERT OR IGNORE INTO users (username, password) VALUES (?, ?)", record[0], record[1])
        if err != nil {
            return err
        }
    }
    return nil
}

func main() {
    // Initialize the SQLite database
    db, err := sql.Open("sqlite3", "./user.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create users table
    createTableSQL := `CREATE TABLE IF NOT EXISTS users (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "username" TEXT NOT NULL UNIQUE,
        "password" TEXT NOT NULL
    );`
    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatal(err)
    }

    // Load users from CSV file
    err = loadUsersFromCSV(db, "users.csv")
    if err != nil {
        log.Fatal(err)
    }

    // Define the handler
    http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var user User
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil {
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }

        var storedPassword string
        err = db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&storedPassword)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if user.Password != storedPassword {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Double-check user/password verification
        err = db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&storedPassword)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if user.Password != storedPassword {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "Success")
    })

    // Start the HTTP server
    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

