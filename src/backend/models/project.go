package models

type Project struct {
    ID         int      `json:"id"`
    Name       string   `json:"name"`
    Categories []string `json:"categories"`
}
