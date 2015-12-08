package main

import(
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "./models"
)

func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Welcome!")
}


func ForceDeploy(githubRepo string) {
  //Force the dheployment inside the dockers
}

func GetInfos(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  i := string(vars["id"])
  h := []int{0,10}
  proc := models.Internals{Max:100, Cur:10, Hist:h}
  dock := models.DockerInfos{Id:i,Proc:proc, Disk:proc, Memory:proc}
  if err := json.NewEncoder(w).Encode(dock); err != nil {
    panic(err);
  }
}
