package models

type User struct {
    ID            int    `json:"id"`
    Username      string `json:"username"`
    Email         string `json:"email"`
    PasswordHash  string `json:"-"`
    DiscordID     string `json:"discord_id"`
    DiscordPseudo string `json:"discord_pseudo"`
}
