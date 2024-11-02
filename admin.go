package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func (r *Router) EnforceAdminMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userId, _ := getUserIDFromContext(req.Context())

		row := r.db.QueryRow("SELECT is_admin FROM users WHERE idm_id = ?", userId)
		var isAdmin bool
		err := row.Scan(&isAdmin)
		if err != nil {
			log.Println("Failed to get if user is admin:", err, userId)
			http.Redirect(w, req, "/signout", http.StatusTemporaryRedirect)
			return
		}

		if !isAdmin {
			log.Println("User is not admin:", userId)
			http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
			return
		}

		handler.ServeHTTP(w, req)
	})
}

func (r *Router) admin(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/admin/users", http.StatusTemporaryRedirect)
}

type AdminUserList struct {
	IDMID       string
	DiscordID   int64
	BuckID      int64
	DisplayName string
	NameNum     string
}

func (r *Router) adminUsers(w http.ResponseWriter, req *http.Request) {
	userId, _ := getUserIDFromContext(req.Context())

	row := r.db.QueryRow("SELECT name_num FROM users WHERE idm_id = ?", userId)
	var nameNum string
	err := row.Scan(&nameNum)
	if err != nil {
		log.Println("Failed to get user:", err, userId)
		http.Redirect(w, req, "/signout", http.StatusTemporaryRedirect)
		return
	}

	page := req.URL.Query().Get("page")
	pageNum, _ := strconv.Atoi(page)

	users := []AdminUserList{}
	rows, err := r.db.Query(`
		SELECT idm_id, discord_id, buck_id, name_num, display_name, last_seen_timestamp, last_attended_timestamp, added_to_mailinglist, student, alum, employee, faculty
		FROM users
		ORDER BY last_seen_timestamp ASC
		LIMIT 100
		OFFSET ?1
	`, pageNum*100)
	if err != nil {
		log.Println("Failed to get users:", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var idmID, nameNum, displayName string
		var discordID, lastAttendedTimestamp sql.NullInt64
		var buckID, lastSeenTimestamp int64
		var addedToMailingList, student, alum, employee, faculty bool
		err := rows.Scan(&idmID, &discordID, &buckID, &nameNum, &displayName, &lastSeenTimestamp, &lastAttendedTimestamp, &addedToMailingList, &student, &alum, &employee, &faculty)
		if err != nil {
			log.Println("Failed to scan user:", err)
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			return
		}
		users = append(users, AdminUserList{
			IDMID:       idmID,
			DiscordID:   discordID.Int64,
			BuckID:      buckID,
			DisplayName: displayName,
			NameNum:     nameNum,
		})
	}

	err = Templates.ExecuteTemplate(w, "admin-users.html.tpl", map[string]interface{}{
		"nameNum": nameNum,
		"path":    req.URL.Path,
		"users":   users,
	})
	if err != nil {
		log.Println("Failed to render template:", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (r *Router) adminVote(w http.ResponseWriter, req *http.Request) {
	userId, _ := getUserIDFromContext(req.Context())

	row := r.db.QueryRow("SELECT name_num FROM users WHERE idm_id = ?", userId)
	var nameNum string
	err := row.Scan(&nameNum)
	if err != nil {
		log.Println("Failed to get user:", err, userId)
		http.Redirect(w, req, "/signout", http.StatusTemporaryRedirect)
		return
	}

	err = Templates.ExecuteTemplate(w, "admin-vote.html.tpl", map[string]interface{}{
		"nameNum": nameNum,
		"path":    req.URL.Path,
	})
	if err != nil {
		log.Println("Failed to render template:", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
