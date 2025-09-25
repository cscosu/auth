package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Structs matching the JSON structure

type Person struct {
	DisplayName   string        `json:"display_name"`
	FirstName     string        `json:"first_name"`
	MiddleName    *string       `json:"middle_name"`
	LastName      string        `json:"last_name"`
	LegalLastName string        `json:"legal_last_name"`
	NameSuffix    *string       `json:"name_suffix"`
	Username      string        `json:"username"`
	Email         string        `json:"email"`
	Address       Address       `json:"address"`
	Phone         Phone         `json:"phone"`
	Majors        []Major       `json:"majors"`
	Affiliations  []string      `json:"affiliations"`
	Appointments  []Appointment `json:"appointments"`
}

type Address struct {
	Building   *Building `json:"building"`
	RoomNumber *string   `json:"room_number"`
	Street1    string    `json:"street1"`
	Street2    *string   `json:"street2"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Zip        string    `json:"zip"`
}

type Building struct {
	Name   string `json:"name"`
	Number string `json:"number"`
	URL    string `json:"url"`
}

type Phone struct {
	AreaCode   string `json:area_code`
	Exchange   string `json:exchange`
	Subscriber string `json:subscriber`
	Formatted  string `json:formatted`
}

type Major struct {
	Major   string `json:"major"`
	College string `json:"college"`
}

type Appointment struct {
	JobTitle      string `json:job_title`
	WorkingTitle  string `json:working_title`
	Organization  string `json:organization`
	OrgCode       string `json:org_code`
	VpCollegeName string `json:vp_college_name`
}

func userExists(userid string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("https://directory.osu.edu/fpjson.php?name_n=%s", userid))
	if err != nil {
		return true, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return true, err
	}

	var people []Person
	err = json.Unmarshal(body, &people)
	if err != nil {
		return true, err
	}

	// if no one shows up, the user does not exist
	if len(people) == 0 {
		return false, nil
	}
	return true, nil
}

func alumnusCheckNextUser(b *DiscordBot) error {
	row := b.Db.QueryRow("SELECT buck_id, discord_id, nameNum FROM users WHERE student=1 AND alum=0 ORDER BY last_alum_check_timestamp ASC LIMIT(1)")
	var buckId string
	var discordId string
	var nameNum string
	err := row.Scan(&buckId, &discordId, &nameNum)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	exists, err := userExists(nameNum)
	if err != nil {
		return err
	}

	if !exists {
		err = b.Alumnify(discordId)
		if err != nil {
			return err
		}
	}

	res, err := b.Db.Exec("UPDATE users SET last_alum_check_timestamp=strftime('%s', 'now') WHERE buck_id=?", buckId)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("expected 1 row updated in alumnusCheckNextUser, got %d", rowsAffected)
	}

	return nil
}
