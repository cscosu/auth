package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (r *Router) processVote(w http.ResponseWriter, req *http.Request) {
	userId, hasUserId := getUserIDFromContext(req.Context())

	if hasUserId {
		candidateId, err := strconv.Atoi(req.FormValue("vote"))
		if err != nil {
			log.Println("Failed to get vote:", err)
			http.Error(w, "Failed to get vote", http.StatusBadRequest)
			return
		}

		row := r.db.QueryRow(`
			UPDATE users SET last_seen_timestamp = strftime('%s', 'now') WHERE buck_id = ?1
			RETURNING name_num
		`, userId)
		var nameNum string
		err = row.Scan(&nameNum)
		if err != nil {
			log.Println("Failed to get user:", err, userId)
			http.Redirect(w, req, "/signout", http.StatusTemporaryRedirect)
			return
		}

		var electionId int
		err = r.db.QueryRow("SELECT election_id FROM elections WHERE timestamp = 1 LIMIT 1").Scan(&electionId)
		if err != nil {
			log.Println("Failed to get election:", err)
			http.Error(w, "Failed to get election", http.StatusInternalServerError)
			return
		}

		_, err = r.db.Exec("INSERT INTO votes (election_id, user_id) VALUES (?, ?)", electionId, userId)
		if err != nil {
			log.Println("Failed to insert vote:", err)
			http.Error(w, "Failed to insert vote", http.StatusInternalServerError)
			return
		}

		_, err = r.db.Exec("UPDATE candidates SET votes = votes + 1 WHERE candidate_id = ?", candidateId)
		if err != nil {
			log.Println("Failed to insert vote:", err)
			http.Error(w, "Failed to insert vote", http.StatusInternalServerError)
			return
		}

		err = Templates.ExecuteTemplate(w, "voting-form.html.tpl", map[string]any{
			"nameNum":  nameNum,
			"hasVoted": true,
		})
		if err != nil {
			log.Println("Failed to render template:", err)
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}

}

func (r *Router) vote(w http.ResponseWriter, req *http.Request) {
	userId, hasUserId := getUserIDFromContext(req.Context())

	if hasUserId {
		var nameNum string
		err := r.db.QueryRow(`
			UPDATE users SET last_seen_timestamp = strftime('%s', 'now') WHERE buck_id = ?1
			RETURNING name_num
		`, userId).Scan(&nameNum)
		if err != nil {
			log.Println("Failed to get user:", err, userId)
			http.Redirect(w, req, "/signout", http.StatusTemporaryRedirect)
			return
		}

		var electionId int
		var electionName string
		err = r.db.QueryRow(`
			SELECT election_id, name FROM elections WHERE timestamp = 1 LIMIT 1
		`).Scan(&electionId, &electionName)
		if err == sql.ErrNoRows {
			err = Templates.ExecuteTemplate(w, "vote.html.tpl", map[string]any{
				"nameNum": nameNum,
			})
			if err != nil {
				log.Println("Failed to render template:", err)
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil {
			log.Println("Failed to get election:", err)
			http.Error(w, "Failed to get election", http.StatusInternalServerError)
			return
		}

		var hasVoted bool
		err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM votes WHERE election_id = ? AND user_id = ?)", electionId, userId).Scan(&hasVoted)
		if err != nil {
			log.Println("Failed to get if user has voted:", err)
			http.Error(w, "Failed to get if user has voted", http.StatusInternalServerError)
			return
		}

		if hasVoted {
			err = Templates.ExecuteTemplate(w, "vote.html.tpl", map[string]any{
				"nameNum":      nameNum,
				"hasVoted":     true,
				"electionName": electionName,
			})
			if err != nil {
				log.Println("Failed to render template:", err)
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			return
		}

		type Candidate struct {
			Id   int
			Name string
		}

		rows, err := r.db.Query("SELECT candidate_id, name FROM candidates WHERE election_id = ?", electionId)
		if err != nil {
			log.Println("Failed to get candidates:", err)
			http.Error(w, "Failed to get candidates", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var candidates []Candidate

		for rows.Next() {
			var candidate Candidate
			err = rows.Scan(&candidate.Id, &candidate.Name)
			if err != nil {
				log.Println("Failed to get candidate:", err)
				http.Error(w, "Failed to get candidate", http.StatusInternalServerError)
				return
			}
			candidates = append(candidates, candidate)
		}

		err = Templates.ExecuteTemplate(w, "vote.html.tpl", map[string]any{
			"nameNum":      nameNum,
			"isMember":     true,
			"electionName": electionName,
			"candidates":   candidates,
		})
		if err != nil {
			log.Println("Failed to render template:", err)
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}
}

func (r *Router) adminVote(w http.ResponseWriter, req *http.Request) {
	userId, _ := getUserIDFromContext(req.Context())

	row := r.db.QueryRow("SELECT name_num FROM users WHERE buck_id = ?", userId)
	var nameNum string
	err := row.Scan(&nameNum)
	if err != nil {
		log.Println("Failed to get user:", err, userId)
		http.Redirect(w, req, "/signout", http.StatusTemporaryRedirect)
		return
	}

	type Election struct {
		ElectionId string
		Name       string
		DoneTime   string
		Done       bool
		Published  bool
	}
	rows, err := r.db.Query("SELECT election_id, name, timestamp FROM elections ORDER BY election_id DESC")
	if err != nil {
		log.Println("Failed to get elections:", err)
		http.Error(w, "Failed to get elections", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	ny, _ := time.LoadLocation("America/New_York")

	var elections []Election
	for rows.Next() {
		var election Election
		var timestamp int64
		err = rows.Scan(&election.ElectionId, &election.Name, &timestamp)
		if err != nil {
			log.Println("Failed to get election:", err)
			http.Error(w, "Failed to get election", http.StatusInternalServerError)
			return
		}

		election.DoneTime = time.Unix(timestamp, 0).In(ny).Format("Mon Jan _2, 2006 at 15:04")
		election.Done = timestamp > 1
		election.Published = timestamp == 1
		elections = append(elections, election)
	}

	err = Templates.ExecuteTemplate(w, "admin-vote.html.tpl", map[string]any{
		"nameNum":       nameNum,
		"path":          req.URL.Path,
		"pastElections": elections,
	})
	if err != nil {
		log.Println("Failed to render template:", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (r *Router) adminVoteEdit(w http.ResponseWriter, req *http.Request) {
	var electionId int
	if req.URL.Path == "/admin/vote/new" {
		row := r.db.QueryRow("INSERT INTO elections (name) VALUES ('New Election') RETURNING election_id")
		err := row.Scan(&electionId)
		if err != nil {
			log.Println("Failed to insert election:", err)
			http.Error(w, "Failed to insert election", http.StatusInternalServerError)
			return
		}

		_, err = r.db.Exec("INSERT INTO candidates (election_id, name) VALUES (?1, 'Candidate 1'), (?1, 'Candidate 2')", electionId)
		if err != nil {
			log.Println("Failed to insert candidates:", err)
			http.Error(w, "Failed to insert candidates", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Push-Url", fmt.Sprintf("/admin/vote/%d", electionId))
	} else {
		id, err := strconv.Atoi(req.PathValue("election_id"))
		if err != nil {
			log.Println("Failed to get election id:", err)
			http.Error(w, "Failed to get election id", http.StatusBadRequest)
			return
		}
		electionId = id
	}

	if req.Method == "PATCH" {
		electionName := req.FormValue("electionName")
		_, err := r.db.Exec("UPDATE elections SET name = ? WHERE election_id = ?", electionName, electionId)
		if err != nil {
			log.Println("Failed to update election:", err)
			http.Error(w, "Failed to update election", http.StatusInternalServerError)
		}
		return
	}

	var electionName string
	var timestamp int64
	err := r.db.QueryRow("SELECT name, timestamp FROM elections WHERE election_id = ?", electionId).Scan(&electionName, &timestamp)
	if err != nil {
		log.Println("Failed to get election:", err)
		http.Error(w, "Failed to get election", http.StatusInternalServerError)
		return
	}
	published := timestamp > 0
	done := timestamp > 1

	type Candidate struct {
		Id    int
		Name  string
		Votes int
	}

	var orderBy string
	if done {
		orderBy = "ORDER BY votes DESC"
	}
	rows, err := r.db.Query("SELECT candidate_id, name, votes FROM candidates WHERE election_id = ?"+orderBy, electionId)
	if err != nil {
		log.Println("Failed to get candidates:", err)
		http.Error(w, "Failed to get candidates", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var candidates []Candidate
	totalVotes := 0

	for rows.Next() {
		var candidate Candidate
		err = rows.Scan(&candidate.Id, &candidate.Name, &candidate.Votes)
		if err != nil {
			log.Println("Failed to get candidate:", err)
			http.Error(w, "Failed to get candidate", http.StatusInternalServerError)
			return
		}
		candidates = append(candidates, candidate)
		totalVotes += candidate.Votes
	}

	if !published && req.Method == "PUT" {
		candidate := Candidate{Name: "New Candidate"}
		err := r.db.QueryRow("INSERT INTO candidates (election_id, name) VALUES (?, 'New Candidate') RETURNING candidate_id", electionId).Scan(&candidate.Id)
		if err != nil {
			log.Println("Failed to insert candidate:", err)
			http.Error(w, "Failed to insert candidate", http.StatusInternalServerError)
			return
		}
		candidates = append(candidates, candidate)

		err = Templates.ExecuteTemplate(w, "admin-vote-edit-partial.html.tpl", map[string]any{
			"electionName": electionName,
			"electionId":   electionId,
			"candidates":   candidates,
		})
		if err != nil {
			log.Println("Failed to render template:", err)
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		return
	}

	userId, _ := getUserIDFromContext(req.Context())

	var nameNum string
	err = r.db.QueryRow("SELECT name_num FROM users WHERE buck_id = ?", userId).Scan(&nameNum)
	if err != nil {
		log.Println("Failed to get user:", err, userId)
		http.Redirect(w, req, "/signout", http.StatusTemporaryRedirect)
		return
	}

	err = Templates.ExecuteTemplate(w, "admin-vote-edit.html.tpl", map[string]any{
		"nameNum":      nameNum,
		"path":         "/admin/vote",
		"electionName": electionName,
		"electionId":   electionId,
		"candidates":   candidates,
		"published":    published,
		"done":         done,
		"timestamp":    timestamp,
		"totalVotes":   totalVotes,
	})
	if err != nil {
		log.Println("Failed to render template:", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (r *Router) adminVoteDelete(w http.ResponseWriter, req *http.Request) {
	electionId, err := strconv.Atoi(req.PathValue("election_id"))
	if err != nil {
		log.Println("Failed to get election id:", err)
		http.Error(w, "Failed to get election id", http.StatusBadRequest)
		return
	}

	_, err = r.db.Exec("DELETE FROM elections WHERE election_id = ?", electionId)
	if err != nil {
		log.Println("Failed to delete election:", err)
		http.Error(w, "Failed to delete election", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/admin/vote", http.StatusTemporaryRedirect)
}

func (r *Router) adminVoteCandidateEdit(w http.ResponseWriter, req *http.Request) {
	electionId, err := strconv.Atoi(req.PathValue("election_id"))
	if err != nil {
		log.Println("Failed to get election id:", err)
		http.Error(w, "Failed to get election id", http.StatusBadRequest)
		return
	}

	candidateId, err := strconv.Atoi(req.PathValue("candidate_id"))
	if err != nil {
		log.Println("Failed to get candidate id:", err)
		http.Error(w, "Failed to get candidate id", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "PATCH":
		candidateName := req.FormValue("candidateName")
		_, err = r.db.Exec("UPDATE candidates SET name = ? WHERE candidate_id = ?", candidateName, candidateId)
		if err != nil {
			log.Println("Failed to update candidate:", err)
			http.Error(w, "Failed to update candidate", http.StatusInternalServerError)
			return
		}
	case "DELETE":
		_, err = r.db.Exec("DELETE FROM candidates WHERE candidate_id = ?", candidateId)
		if err != nil {
			log.Println("Failed to delete candidate:", err)
			http.Error(w, "Failed to delete candidate", http.StatusInternalServerError)
			return
		}

		var electionName string
		err := r.db.QueryRow("SELECT name FROM elections WHERE election_id = ?", electionId).Scan(&electionName)
		if err != nil {
			log.Println("Failed to get election:", err)
			http.Error(w, "Failed to get election", http.StatusInternalServerError)
			return
		}

		type Candidate struct {
			Id   int
			Name string
		}

		rows, err := r.db.Query("SELECT candidate_id, name FROM candidates WHERE election_id = ?", electionId)
		if err != nil {
			log.Println("Failed to get candidates:", err)
			http.Error(w, "Failed to get candidates", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var candidates []Candidate

		for rows.Next() {
			var candidate Candidate
			err = rows.Scan(&candidate.Id, &candidate.Name)
			if err != nil {
				log.Println("Failed to get candidate:", err)
				http.Error(w, "Failed to get candidate", http.StatusInternalServerError)
				return
			}
			candidates = append(candidates, candidate)
		}

		err = Templates.ExecuteTemplate(w, "admin-vote-edit-partial.html.tpl", map[string]any{
			"electionName": electionName,
			"electionId":   electionId,
			"candidates":   candidates,
		})
		if err != nil {
			log.Println("Failed to render template:", err)
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}
}

func (r *Router) adminVotePublish(w http.ResponseWriter, req *http.Request) {
	electionId, err := strconv.Atoi(req.PathValue("election_id"))
	if err != nil {
		log.Println("Failed to get election id:", err)
		http.Error(w, "Failed to get election id", http.StatusBadRequest)
		return
	}

	_, err = r.db.Exec("UPDATE elections SET timestamp = 0 WHERE timestamp = 1", electionId)
	if err != nil {
		log.Println("Failed to update elections:", err)
		http.Error(w, "Failed to update elections", http.StatusInternalServerError)
		return
	}

	_, err = r.db.Exec("UPDATE elections SET timestamp = 1 WHERE election_id = ?", electionId)
	if err != nil {
		log.Println("Failed to update election:", err)
		http.Error(w, "Failed to update election", http.StatusInternalServerError)
		return
	}

	var electionName string
	err = r.db.QueryRow("SELECT name FROM elections WHERE election_id = ?", electionId).Scan(&electionName)
	if err != nil {
		log.Println("Failed to get election:", err)
		http.Error(w, "Failed to get election", http.StatusInternalServerError)
		return
	}

	type Candidate struct {
		Id   int
		Name string
	}

	rows, err := r.db.Query("SELECT candidate_id, name FROM candidates WHERE election_id = ?", electionId)
	if err != nil {
		log.Println("Failed to get candidates:", err)
		http.Error(w, "Failed to get candidates", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var candidates []Candidate

	for rows.Next() {
		var candidate Candidate
		err = rows.Scan(&candidate.Id, &candidate.Name)
		if err != nil {
			log.Println("Failed to get candidate:", err)
			http.Error(w, "Failed to get candidate", http.StatusInternalServerError)
			return
		}
		candidates = append(candidates, candidate)
	}

	err = Templates.ExecuteTemplate(w, "admin-vote-edit-partial.html.tpl", map[string]any{
		"published":    true,
		"electionName": electionName,
		"electionId":   electionId,
		"candidates":   candidates,
		"totalVotes":   0,
	})
	if err != nil {
		log.Println("Failed to render template:", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (r *Router) adminVoteClose(w http.ResponseWriter, req *http.Request) {
	electionId, err := strconv.Atoi(req.PathValue("election_id"))
	if err != nil {
		log.Println("Failed to get election id:", err)
		http.Error(w, "Failed to get election id", http.StatusBadRequest)
		return
	}

	_, err = r.db.Exec("UPDATE elections SET timestamp = strftime('%s', 'now') WHERE election_id = ?", electionId)
	if err != nil {
		log.Println("Failed to update election:", err)
		http.Error(w, "Failed to update election", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/admin/vote/%d", electionId), http.StatusTemporaryRedirect)
}
