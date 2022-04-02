package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/alekstet/social_graph/conf"
	"github.com/julienschmidt/httprouter"

	_ "github.com/lib/pq"
)

type Store struct {
	conf   conf.Conf
	router *httprouter.Router
	db     *sql.DB
}

func New(config *conf.Conf) *Store {
	return &Store{
		conf:   *config,
		router: httprouter.New(),
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

func (s *Store) GetFromDB() ([]byte, error) {
	journal_soc := [][3]int{}
	var min, max, max_to, max_from, max_node int

	rows, err := s.db.Query("SELECT p_from, p_to, COUNT(*) FROM public.social GROUP BY p_from, p_to")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var soc [3]int
		err := rows.Scan(&soc[0], &soc[1], &soc[2])
		if err != nil {
			return nil, err
		}
		journal_soc = append(journal_soc, soc)
	}

	err = s.db.QueryRow("SELECT MAX(p_to), MAX(p_from) FROM public.social").Scan(&max_to, &max_from)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	resp := Solve(max_node, min, max, journal_soc)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	} else if len(resp.Matrix) == 0 {
		return nil, err
	} else {
		return jsonResp, nil
	}
}

func Solve(max_node, min, max int, journal_soc [][3]int) Resp {
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
	return resp
}

func (s *Store) GetMatrix(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := s.GetFromDB()
	if err != nil {
		w.WriteHeader(500)
		return
	} else {
		w.Write(res)
	}
}

func (s *Store) PutData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		return
	} else if from <= 0 || to <= 0 {
		w.WriteHeader(400)
		return
	} else {
		_, err := s.db.Exec("INSERT INTO public.social (p_from, p_to) VALUES ($1, $2)", from, to)
		if err != nil {
			w.WriteHeader(500)
			return
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
	conf, err := conf.Cnf()
	if err != nil {
		log.Fatalf("error with config file: %s", err)
	}
	s := New(conf)
	err = s.DB()
	if err != nil {
		log.Fatalf("error with db: %s", err)
	}
	defer s.db.Close()

	s.router.GET("/social", s.GetMatrix)
	s.router.PUT("/social", s.PutData)

	err = http.ListenAndServe(conf.PortApp, s.router)
	if err != nil {
		log.Fatalf("error with serve: %s", err)
	}
}
