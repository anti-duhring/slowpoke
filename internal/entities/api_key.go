package entities

type ApiKey struct {
  ID int `json:"id"`
  Name string `json:"name"`
  UserID int `json:"user_id"`
}
