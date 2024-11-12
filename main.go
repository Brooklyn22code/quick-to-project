package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func main() {
	db, err := sql.Open("sqlite", "./project.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS project (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        path TEXT NOT NULL
    )`)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("command> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "create":
			var project Project
			fmt.Print("project name: ")
			project.Name, _ = reader.ReadString('\n')
			project.Name = strings.TrimSpace(project.Name)

			fmt.Print("project path: ")
			project.Path, _ = reader.ReadString('\n')
			project.Path = strings.TrimSpace(project.Path)

			result, err := db.Exec("INSERT INTO project (name, path) VALUES (?, ?)", project.Name, project.Path)
			if err != nil {
				fmt.Println("Error inserting project:", err)
				continue
			}
			id, _ := result.LastInsertId()
			fmt.Printf("Project created with ID: %d\n", id)

		case "open":
			fmt.Print("project ID: ")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID format")
				continue
			}

			var project Project
			row := db.QueryRow("SELECT id, name, path FROM project WHERE id =?", id)
			if err := row.Scan(&project.Id, &project.Name, &project.Path); err == sql.ErrNoRows {
				fmt.Println("Project not found")
				continue
			} else if err != nil {
				fmt.Println("Error fetching project:", err)
				continue
			}

			project.Path = strings.Trim(project.Path, `"'`)

			fmt.Println("Changing to project path:", project.Path)
			if err = os.Chdir(project.Path); err != nil {
				fmt.Println("Error changing directory:", err)
				continue
			}

			cmd := exec.Command("code", ".")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Println("Error opening project with editor:", err)
			}

		case "list":
			rows, err := db.Query("SELECT id, name, path FROM project")
			if err != nil {
				fmt.Println("Error fetching projects:", err)
				continue
			}
			defer rows.Close()
			fmt.Println("ID\tName\tPath")
			for rows.Next() {
				var project Project
				if err := rows.Scan(&project.Id, &project.Name, &project.Path); err != nil {
					fmt.Println("Error scanning row:", err)
					continue
				}
				fmt.Printf("%d\t%s\t%s\n", project.Id, project.Name, project.Path)
			}

		case "delete":
			fmt.Print("project ID: ")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID format")
				continue
			}

			result, err := db.Exec("DELETE FROM project WHERE id =?", id)
			if err != nil {
				fmt.Println("Error deleting project:", err)
				continue
			}
			rowsAffected, _ := result.RowsAffected()
			fmt.Printf("%d project(s) deleted\n", rowsAffected)

		case "exit":
			return

		default:
			fmt.Println("Invalid command")
		}
	}
}
