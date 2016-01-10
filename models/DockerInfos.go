package models




type Internals struct {
  Max int `json:"max"`
  Cur float64 `json:"cur"`
  Hist []float64 `json:"hist"`
}

type DockerInfos struct {
  Id  string `json:"id"`
  Proc Internals `json:"proc"`
  Disk Internals `json:"disk"`
  Memory Internals `json:"memory"`
}
