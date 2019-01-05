package smashgg

import (
	"encoding/json"
	"fmt"

	"github.com/groundstorm/cobalt/apps/go/src/models"
	"github.com/groundstorm/cobalt/apps/go/src/util"
)

// ImportGroup pulls in data from smash.gg for the specified group id
func ImportGroup(groupID int) (*models.Group, error) {
	query := "https://api.smash.gg/phase_group/%d?expand[]=sets"
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
	return g, nil
}

// ImportEvent pulls in data from smash.gg for the specified event id
func ImportEvent(eventID int) (*models.Event, error) {
	query := "https://api.smash.gg/event/%d?expand[]=phase&expand[]=groups"
	url := fmt.Sprintf(query, eventID)
	bytes, err := util.GetURL(url, log)
	if err != nil {
		return nil, err
	}

	// Parse the result.
	var q EventQuery
	err = json.Unmarshal(bytes, &q)
	if err != nil {
		return nil, err
	}

	e := &models.Event{
		ID:   models.EventID(q.Entities.Event.ID),
		Slug: q.Entities.Event.Slug,
	}
	// Add in all the phases
	phases := map[models.PhaseID]*models.Phase{}
	for _, qp := range q.Entities.Phases {
		p := &models.Phase{
			ID:     models.PhaseID(qp.ID),
			Name:   qp.Name,
			Groups: []*models.Group{},
		}
		phases[p.ID] = p
		fmt.Printf("-> %d %v\n", qp.ID, phases[p.ID])
		e.Phases = append(e.Phases, p)
	}

	// Add all the groups to those phases.
	for _, qg := range q.Entities.Groups {
		g, err := ImportGroup(qg.ID)
		if err != nil {
			return nil, err
		}
		p := phases[g.PhaseID]
		p.Groups = append(p.Groups, g)
	}
	return e, nil
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
		ID:   models.TournamentID(q.Entities.Tournament.ID),
		Name: q.Entities.Tournament.Name,
		Slug: q.Entities.Tournament.Slug,
	}
	for _, qe := range q.Entities.Events {
		e, err := ImportEvent(qe.ID)
		if err != nil {
			return nil, err
		}
		t.Events = append(t.Events, e)
	}

	// We need to send additional queries to get specific info about the events
	return t, nil
}
