package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"

	_ "github.com/lib/pq"

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

type Anal struct{}

func (a Anal) Print(searchTimes []float32) {
	fmt.Println("Median (100):", a.CalcMedian(searchTimes), "ms")
	fmt.Println("Avg (100):", a.CalcAvg(searchTimes), "ms")
	fmt.Println("Min (100):", a.CalcMin(searchTimes), "ms")
	fmt.Println("Max (100):", a.CalcMax(searchTimes), "ms")
}

func (a Anal) CalcAvg(n []float32) float32 {
	var sum float32 = 0

	for _, t := range n {
		sum += t
	}

	return sum / float32(len(n))
}

func (a Anal) CalcMin(n []float32) float32 {
	var min float32 = 1000

	for _, t := range n {
		if t < min {
			min = t
		}
	}

	return min
}

func (a Anal) CalcMax(n []float32) float32 {
	var max float32 = 0

	for _, t := range n {
		if t > max {
			max = t
		}
	}

	return max
}

func (a Anal) CalcMedian(n []float32) float32 {
	sort.Slice(n, func(i, j int) bool { return n[i] < n[j] })

	mNumber := len(n) / 2

	if len(n)%2 != 0 {
		return n[mNumber]
	}

	return float32(n[mNumber-1]+n[mNumber]) / 2.0
}

func main() {
	words := []string{"layer", "opposite", "waist", "become", "address", "adult", "upper", "twelve", "card", "prefer", "patient", "concerning", "welcome", "bread", "connect", "beyond", "law", "northern", "more", "gray", "west", "except", "OK", "negative", "nation", "program", "plenty", "wine", "information", "produce", "animal", "smart", "fear", "lock", "upper", "physical", "beautiful", "truck", "steady", "card", "walk", "rock", "bear", "grass", "hand", "odd", "proof", "decrease", "represent", "over", "quiet", "solve", "require", "important", "inform", "nose", "very", "crowd", "third", "request", "woman", "practical", "invite", "adjective", "wake", "soon", "itself", "relation", "fork", "food", "average", "change", "well", "each", "quality", "supply", "point", "dollar", "child", "pound", "balance", "suddenly", "cook", "notice", "traffic", "recognize", "drunk", "toilet", "always", "say", "reason", "under", "forget", "replace", "medical", "clothes", "breast", "straight", "duck", "admit"}

	dburl := "postgresql://roswell@127.0.0.1:5432/nextdb?sslmode=disable"

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
				if err := rows.Scan(&r.DisplayName, &r.Description); err != nil {
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
			a := Anal{}
			a.Print(searchTimes)
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
			if err := rows.Scan(&r.IncidentId, &r.DisplayName, &r.Description); err != nil {
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

		tplResults.Execute(w, map[string]interface{}{
			"Results": results,
			"Query":   q,
		})

	})

	log.Fatal(http.ListenAndServe(":1337", nil))
}
