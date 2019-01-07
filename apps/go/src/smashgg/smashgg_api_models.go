package smashgg

type TournamentRegistrationsQuery struct {
	Entities struct {
		Tournament struct {
			ID int `json:"id"`
		} `json:"tournament"`
		Events []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"event"`
	} `json:"entities"`
}

type TournamentQuery struct {
	Entities struct {
		Tournament struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"tournament"`
		Events []struct {
			ID int `json:"id"`
		} `json:"event"`
	} `json:"entities"`
}

type EventQuery struct {
	Entities struct {
		Event struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Slug        string `json:"slug"`
		} `json:"event"`
		Phases []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			CreatedAt int    `json:"createdAt"`
			UpdatedAt int    `json:"updatedAt"`
		} `json:"phase"`
		Groups []struct {
			ID int `json:"id"`
		} `json:"groups"`
	} `json:"entities"`
}
type GroupQuery struct {
	Entities struct {
		Group struct {
			ID                   int    `json:"id"`
			PhaseID              int    `json:"phaseId"`
			PoolRefParticipantID int    `json:"poolRefId"`
			StationID            int    `json:"stationId"`
			Name                 string `json:"displayIdentifier"`
			StartAt              int    `json:"startAt"`
		} `json:"groups"`
		Sets []struct {
			ID int `json:"id"`
		} `json:"sets"`
		Seeds []struct {
			ID        int  `json:"id"`
			EntrantID int  `json:"entrantId"`
			Num       int  `json:"seedNum"`
			IsBye     bool `json:"isBye"`
		} `json:"seeds"`
	} `json:"entities"`
}
