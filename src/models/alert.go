package models

type Alert struct {
	Id   string `json:"id"`
	Body string `json:"body"`
}

func (a *Alert) Execute() error {
	return nil
}
