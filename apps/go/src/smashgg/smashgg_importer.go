package smashgg

import (
	"bytes"
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
		Name: q.Entities.Event.Name,
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

// GetTournamentRegistrationInfo gets the tournament info required to interpret the
// export of registered players.
func GetTournamentRegistrationInfo(slug string) (*TournamentRegistrationsQuery, error) {
	q := &TournamentRegistrationsQuery{}
	url := fmt.Sprintf("https://api.smash.gg/tournament/%s?expand[]=event", slug)
	bytes, err := util.GetURL(url, log)
	if err != nil {
		return q, err
	}
	err = json.Unmarshal(bytes, q)
	if err != nil {
		return q, err
	}
	return q, nil
}

// LoadAttendeesRaw downloads the raw attendee data from smash.gg (without processing)
func LoadAttendeesRaw(info *TournamentRegistrationsQuery) ([]byte, error) {
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
	url := fmt.Sprintf("https://smash.gg/api-proxy/tournament/%d/export_attendees", info.Entities.Tournament.ID)
	exportResponse, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer exportResponse.Body.Close()
	return ioutil.ReadAll(exportResponse.Body)
}

// LoadAttendees fetches the list of attendees from smash.gg and converts it into
// our attendees model.
func LoadAttendees(info *TournamentRegistrationsQuery) (*models.Attendees, error) {
	data, err := LoadAttendeesRaw(info)
	if err != nil {
		return nil, err
	}

	// Parse!
	r := csv.NewReader(bytes.NewReader(data))

	r.LazyQuotes = true

	// The exporter for evo 2017 data for some reason is not returning the same number
	// of fields for all records.  This will turn off that error
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	columns := map[string]int{}
	for i, name := range header {
		columns[name] = i
	}
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	a := &models.Attendees{
		Registrations: map[models.ParticipantID]models.Registration{},
		Events:        map[models.EventID]models.Event{},
	}

	for _, e := range info.Entities.Events {
		eid := models.EventID(e.ID) // XXX: can we just store this as models.EventID in the Query type?
		a.Events[eid] = models.Event{
			ID:   eid,
			Name: e.Name,
		}
	}

	// evo 2017 exports have a single "name" field.
	iName := columns["Name"]
	// evo 2017 exports have both first name and last name
	iFirstName := columns["First Name"]
	iLastName := columns["Last Name"]

	iID := columns["Id"]
	iEmail := columns["Email"]

	for _, record := range records {
		id, _ := strconv.Atoi(record[iID])

		// Do the easy stuff first.
		r := models.Registration{
			Participant: models.Participant{
				ID: models.ParticipantID(id),
				Player: models.Player{
					FirstName: record[iFirstName],
					LastName:  record[iLastName],
					Email:     models.Email(record[iEmail]),
				},
			},
			Events: map[models.EventID]bool{},
		}
		// add in all the games.
		for _, e := range info.Entities.Events {
			iEvent := columns[e.Name]
			if iEvent > 0 {
				r.Events[models.EventID(e.ID)] = true
			}
		}

		// The format of the name changed between 2017 and 2018.  Do that now.
		if iName > 0 {
			names := strings.SplitN(record[iName], " ", 2)
			r.Participant.Player.FirstName = names[0]
			if len(names) > 1 {
				r.Participant.Player.LastName = names[1]
			} else {
				r.Participant.Player.LastName = ""
			}
		} else {
			r.Participant.Player.FirstName = record[iFirstName]
			r.Participant.Player.LastName = record[iLastName]
		}

		a.Registrations[r.Participant.ID] = r
	}

	return a, nil
}
