package smashgg

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"net/http"
	"net/http/cookiejar"

	"github.com/groundstorm/cobalt/apps/go/src/models"
	"github.com/groundstorm/cobalt/apps/go/src/util"
)

const ()

var (
	cookieJar *cookiejar.Jar
)

func init() {
	var err error
	cookieJar, err = cookiejar.New(nil)
	if err != nil {
		log.Fatalf("failed to create http cookie jar for requests: %s", err)
	}
}

// ImportGroup pulls in data from smash.gg for the specified group id
func ImportGroup(t *models.Tournament, groupID int) (*models.Group, error) {
	query := "https://api.smash.gg/phase_group/%d?expand[]=sets&expand[]=seeds"
	url := fmt.Sprintf(query, groupID)
	bytes, err := util.GetURL(url, log)
	if err != nil {
		return nil, err
	}

	// Parse the result.
	var q GroupQuery
	err = json.Unmarshal(bytes, &q)
	if err != nil {
		return nil, err
	}

	g := &models.Group{
		ID:      models.GroupID(q.Entities.Group.ID),
		PhaseID: models.PhaseID(q.Entities.Group.PhaseID),
		Sets:    []*models.Set{},
	}
	t.Groups[g.ID] = g

	// Read in all the Sets for this Group
	for _, qs := range q.Entities.Sets {
		s := &models.Set{
			ID: models.SetID(qs.ID),
		}
		t.Sets[s.ID] = s
		g.Sets = append(g.Sets, s)
	}

	return g, nil
}

// ImportEvent pulls in data from smash.gg for the specified event id
func ImportEvent(t *models.Tournament, eventID int) error {
	query := "https://api.smash.gg/event/%d?expand[]=phase&expand[]=groups"
	url := fmt.Sprintf(query, eventID)
	bytes, err := util.GetURL(url, log)
	if err != nil {
		return err
	}

	// Parse the result.
	var q EventQuery
	err = json.Unmarshal(bytes, &q)
	if err != nil {
		return err
	}

	e := &models.Event{
		ID:   models.EventID(q.Entities.Event.ID),
		Slug: q.Entities.Event.Slug,
	}
	// Add in all the phases
	for _, qp := range q.Entities.Phases {
		p := &models.Phase{
			ID:     models.PhaseID(qp.ID),
			Name:   qp.Name,
			Groups: []*models.Group{},
		}
		t.Phases[p.ID] = p
		e.Phases = append(e.Phases, p)
	}

	// Add all the groups to those phases.
	for _, qg := range q.Entities.Groups {
		g, err := ImportGroup(t, qg.ID)
		if err != nil {
			return err
		}
		p := t.Phases[g.PhaseID]
		p.Groups = append(p.Groups, g)
	}

	t.Events[e.ID] = e
	return nil
}

// ImportTournament pulls in data from smash.gg for the specified tournament
func ImportTournament(slug string) (*models.Tournament, error) {
	// Get all the data.
	query := "https://api.smash.gg/tournament/%s?expand[]=event"
	url := fmt.Sprintf(query, slug)
	bytes, err := util.GetURL(url, log)
	if err != nil {
		return nil, err
	}

	// Parse the wrapper.
	var q TournamentQuery
	err = json.Unmarshal(bytes, &q)
	if err != nil {
		return nil, err
	}

	// Build the result as we go
	t := &models.Tournament{
		ID:           models.TournamentID(q.Entities.Tournament.ID),
		Name:         q.Entities.Tournament.Name,
		Slug:         q.Entities.Tournament.Slug,
		Events:       map[models.EventID]*models.Event{},
		Phases:       map[models.PhaseID]*models.Phase{},
		Groups:       map[models.GroupID]*models.Group{},
		Sets:         map[models.SetID]*models.Set{},
		Slots:        map[models.SlotID]*models.Slot{},
		Participants: map[models.ParticipantID]*models.Participant{},
		Entrants:     map[models.EntrantID]*models.Entrant{},
	}
	for _, qe := range q.Entities.Events {
		err = ImportEvent(t, qe.ID)
		if err != nil {
			return nil, err
		}
	}

	// We need to send additional queries to get specific info about the events
	return t, nil
}

// GetTournamentID sets the tournament ID from the slug.
func GetTournamentID(slug string) (models.TournamentID, error) {
	url := fmt.Sprintf("https://api.smash.gg/tournament/%s", slug)
	bytes, err := util.GetURL(url, log)
	if err != nil {
		return models.TournamentID(0), err
	}
	q := struct {
		Entities struct {
			Tournament struct {
				ID int
			}
		}
	}{}
	err = json.Unmarshal(bytes, &q)
	if err != nil {
		return models.TournamentID(0), err
	}
	return models.TournamentID(q.Entities.Tournament.ID), nil
}

func FetchAttendees(slug string) ([]byte, error) {
	client := &http.Client{
		Jar: cookieJar,
	}

	// login.
	const loginFmt = `{"requests":{"g0":{"resource":"gg_api./user/login","operation":"create","params":{"email":"%s","password":"%s","rememberMe":true,"validationKey":"LOGIN_userlogin","expand":[]}}},"context":{}}`
	email := "tony.cannon@gmail.com"
	password := "ww31lvamm"
	body := fmt.Sprintf(loginFmt, email, password)

	response, err := client.Post("https://smash.gg/api/-", "application/json", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Now grab the attendees
	id, err := GetTournamentID(slug)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://smash.gg/api-proxy/tournament/%d/export_attendees", id)
	exportResponse, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer exportResponse.Body.Close()
	return ioutil.ReadAll(exportResponse.Body)
}

func LoadAttendees(str string) (*models.Attendees, error) {
	// Parse!
	r := csv.NewReader(strings.NewReader(str))

	r.LazyQuotes = true

	// The exporter for evo 2017 data for some reason is not returning the same number
	// of fields for all records.  This will turn off that error
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	pos := map[string]int{}
	for i, name := range header {
		pos[name] = i
	}
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	a := &models.Attendees{
		Participants: make([]models.Participant, len(records), len(records)),
	}

	// evo 2017 exports have a single "name" field.
	name_offset := pos["Name"]

	// evo 2017 exports have both first name and last name
	first_offset := pos["First Name"]
	last_offset := pos["Last Name"]

	id_offset := pos["Id"]
	email_offset := pos["Email"]
	for i, record := range records {
		id, _ := strconv.Atoi(record[id_offset])
		p := models.Participant{
			ID: models.ParticipantID(id),
			Player: models.Player{
				FirstName: record[first_offset],
				LastName:  record[last_offset],
				Email:     models.Email(record[email_offset]),
			},
		}

		if name_offset > 0 {
			names := strings.SplitN(record[name_offset], " ", 2)
			p.Player.FirstName = names[0]
			if len(names) > 1 {
				p.Player.LastName = names[1]
			} else {
				p.Player.LastName = ""
			}
		} else {
			p.Player.FirstName = record[first_offset]
			p.Player.LastName = record[last_offset]
		}
		a.Participants[i] = p
	}

	return a, nil
}
