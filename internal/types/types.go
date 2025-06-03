package types

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Todo struct {
	ID        string `json:"id"`
	Group     string `json:"group"`
	Urgency   int    `json:"urgency"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

type Config struct {
	Groups      []Group `json:"groups"`
	ActiveGroup string  `json:"active_group"`
	Todos       []Todo  `json:"todos"`
}
