package models




type Internals struct {
  Max int `json:"max"`
  Cur int `json:"cur"`
  Hist []int `json:"hist"`
}

type DockerInfos struct {
  Id  string `json:"id"`
  Proc Internals `json:"proc"`
  Disk Internals `json:"disk"`
  Memory Internals `json:"memory"`
}
