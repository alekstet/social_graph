package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	cnf "github.com/alekstet/social_graph/conf"

	_ "github.com/lib/pq"
)

type Store struct {
	conf cnf.Conf
	db   *sql.DB
}

func New(config *cnf.Conf) *Store {
	return &Store{
		conf: *config,
	}
}

type Resp struct {
	Matrix [][]int `json:"matrix"`
	Info   `json:"info"`
}

type Info struct {
	Max int     `json:"max"`
	Min int     `json:"min"`
	Avg float32 `json:"avg"`
}

func (s *Store) S(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		journal_soc := [][3]int{}
		var min, max, max_to, max_from, max_node int

		rows, err := s.db.Query("SELECT p_from, p_to, COUNT(*) FROM public.social GROUP BY p_from, p_to")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var soc [3]int
			err := rows.Scan(&soc[0], &soc[1], &soc[2])
			if err != nil {
				w.WriteHeader(500)
				return
			}
			journal_soc = append(journal_soc, soc)
		}

		err = s.db.QueryRow("SELECT MAX(p_to), MAX(p_from) FROM public.social").Scan(&max_to, &max_from)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		if max_to > max_from {
			max_node = max_to
		} else {
			max_node = max_from
		}

		err = s.db.QueryRow(
			`SELECT MIN(cnt), MAX(cnt) FROM 
			(SELECT p_from, p_to, COUNT(p_from) cnt 
			FROM public.social GROUP BY p_from, p_to) cnt`).Scan(&min, &max)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		graf := make([][]int, max_node)
		for i := range graf {
			graf[i] = make([]int, max_node)
		}

		for _, j := range journal_soc {
			graf[j[0]-1][j[1]-1] = j[2]
			graf[j[1]-1][j[0]-1] = j[2]
		}

		avg := (float32(min) + float32(max)) / 2
		inf := Info{max, min, avg}
		resp := Resp{graf, inf}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(500)
			return
		} else if len(graf) == 0 {
			w.WriteHeader(500)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResp)
		}

	case "PUT":
		from, err := strconv.Atoi(r.URL.Query().Get("from"))
		if err != nil {
			w.WriteHeader(400)
			return
		}
		to, err := strconv.Atoi(r.URL.Query().Get("to"))
		if err != nil {
			w.WriteHeader(400)
			return
		}
		if from == to {
			w.WriteHeader(400)
		} else if from <= 0 || to <= 0 {
			w.WriteHeader(400)
		} else {
			_, err := s.db.Exec("INSERT INTO public.social (p_from, p_to) VALUES ($1, $2)", from, to)
			if err != nil {
				w.WriteHeader(500)
				return
			}
		}
	}
}

func (s *Store) DB() error {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		s.conf.Host, s.conf.PortBase, s.conf.User, s.conf.Password, s.conf.DBName)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}

	result, err := db.Query("CREATE TABLE IF NOT EXISTS public.social (p_from INTEGER, p_to INTEGER)")
	if err != nil || result == nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func main() {
	conf, err := cnf.Cnf()
	if err != nil {
		log.Fatalf("error with config file: %s", err)
	}
	s := New(conf)
	err = s.DB()
	if err != nil {
		log.Fatalf("error with db: %s", err)
	}
	defer s.db.Close()
	http.HandleFunc("/social", s.S)

	err = http.ListenAndServe(conf.PortApp, nil)
	if err != nil {
		log.Fatalf("error with serve: %s", err)
	}
}
