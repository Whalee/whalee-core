package models

type ProjectRequest struct {
  User string `json:"user"`
  Project string `json:"project"`
}

type ProjectResponse struct {
  Id string `json:"id"`
}
