package models

import "time"

type Task struct {
    ID            int       `json:"id"`
    Name          string    `json:"name"`
    Description   string    `json:"description"`
    Category      string    `json:"category"`
    ProjectID     int       `json:"project_id"`
    Status        string    `json:"status"`
    AssignedUsers []int     `json:"assigned_users"`
    EstimatedTime int       `json:"estimated_time"`
    RealTime      int       `json:"real_time"`
    CreatedAt     time.Time `json:"created_at"`
    AvailableFrom time.Time `json:"available_from"`
}
