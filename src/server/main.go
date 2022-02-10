package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"

	"fts_pg/src/server/analytics"
	"fts_pg/src/server/sqlstore"
	"fts_pg/src/server/templates"
)

var db *sql.DB

type result struct {
	IncidentId  string
	DisplayName string
	Description string
	Snippet     template.HTML
}

func main() {
	words := []string{"layer", "opposite", "waist", "become", "address", "adult", "upper", "twelve", "card", "prefer", "patient", "concerning", "welcome", "bread", "connect", "beyond", "law", "northern", "more", "gray", "west", "except", "OK", "negative", "nation", "program", "plenty", "wine", "information", "produce", "animal", "smart", "fear", "lock", "upper", "physical", "beautiful", "truck", "steady", "card", "walk", "rock", "bear", "grass", "hand", "odd", "proof", "decrease", "represent", "over", "quiet", "solve", "require", "important", "inform", "nose", "very", "crowd", "third", "request", "woman", "practical", "invite", "adjective", "wake", "soon", "itself", "relation", "fork", "food", "average", "change", "well", "each", "quality", "supply", "point", "dollar", "child", "pound", "balance", "suddenly", "cook", "notice", "traffic", "recognize", "drunk", "toilet", "always", "say", "reason", "under", "forget", "replace", "medical", "clothes", "breast", "straight", "duck", "admit"}

	dburl := "postgresql://postgres@127.0.0.1:5432/nextdb?sslmode=disable"

	var err error
	if db, err = sql.Open("postgres", dburl); err != nil {
		log.Fatal(err)
	}

	tplHome := template.Must(template.New(".").Parse(templates.TplStrHome))
	tplResults := template.Must(template.New(".").Parse(templates.TplStrResults))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.FormValue("q")
		fmt.Println(q)
		if q == "" {
			rows, err := db.Query(sqlstore.ListRules)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer rows.Close()
			results := make([]result, 0, 10)
			for rows.Next() {
				var r result
				var snip string
				if err := rows.Scan(&r.IncidentId, &r.DisplayName, &r.Description); err != nil {
					http.Error(w, err.Error(), 404)
					return
				}
				r.Snippet = template.HTML(strings.Replace(snip, "\n", "<br>", -1))
				results = append(results, r)
			}

			tplHome.Execute(w, map[string]interface{}{
				"Results": results,
			})
			return
		}
		if q == "calculate" {
			searchTimes := make([]float32, 0, 100)
			for i, word := range words {
				println(i, "/100")

				rows, err := db.Query(sqlstore.CalcSearchIncidents, word)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				defer rows.Close()
				var explain []map[string]interface{}
				for rows.Next() {
					var explainRaw []uint8
					if err := rows.Scan(&explainRaw); err != nil {
						http.Error(w, err.Error(), 404)
						return
					}

					if err = json.Unmarshal(explainRaw, &explain); err != nil {
						http.Error(w, err.Error(), 500)
						return
					}
				}
				executionTime := float32(explain[0]["Execution Time"].(float64))
				searchTimes = append(searchTimes, executionTime)
			}
			a := analytics.CreateAnal()
			a.Show(searchTimes)
			return
		}

		if len(q) > 100 {
			q = q[:100]
		}
		rows, err := db.Query(sqlstore.SearchIncidents, q)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()
		results := make([]result, 0, 10)
		for rows.Next() {
			var r result
			var snip string
			if err := rows.Scan(&r.IncidentId); err != nil {
				http.Error(w, err.Error(), 404)
				return
			}
			r.Snippet = template.HTML(strings.Replace(snip, "\n", "<br>", -1))
			results = append(results, r)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), 404)
			return
		}

		err = tplResults.Execute(w, map[string]interface{}{
			"Results": results,
			"Query":   q,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	const PORT = "1337"

	fmt.Printf("Server starting on %v port\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
