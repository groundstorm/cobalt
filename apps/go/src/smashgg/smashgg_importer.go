package smashgg

import (
	"fmt"

	"github.com/groundstorm/cobalt/apps/go/src/models"
	"github.com/groundstorm/cobalt/apps/go/src/util"
)

// ImportEvent pulls in data from smash.gg for the specified tournament slug
func ImportEvent(slug string) (*models.Event, error) {
	// Get all the data.
	query := "https://api.smash.gg/tournament/%s?expand[]=event&expand[]=phase&expand[]=groups&expand[]=stations"
	url := fmt.Sprintf(query, slug)
	data, err := util.GetJSON(url)

	// Parse the blob.   It's largely undocumented >_<
	evt := &models.Event{}
	return evt, nil
}
