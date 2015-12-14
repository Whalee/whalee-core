package routes
import(
  "fmt"
  "net/http"
  "../models"
  "io"
  "io/ioutil"
  "log"
  "encoding/json"
  "../externals"
)
/*
 * POST /projects
 */
func PostProjects(w http.ResponseWriter, r *http.Request) {
  var project models.ProjectRequest
  //Extract github repo url
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    panic(err);
  }
  if err := r.Body.Close(); err != nil {
    panic(err);
  }
  if err := json.Unmarshal(body, &project); err != nil {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(422) // unprocessable entity
    if err := json.NewEncoder(w).Encode(err); err != nil {
        panic(err)
    }
  }
  log.Println("Creating a docker");
  //local docker
  dockerClient := externals.NewLocalInteractor("unix:///var/run/docker.sock");
  config := externals.Config {
    User: "magicmicky",
    Project: "project1",
  }
  dockerClient.RunContainer(config);


  log.Println("Retrieving github " + project.Github + " from docker");
  //Call inside the executed docker the set up server

  fmt.Fprintln(w, "ok")
}

/*
 * GET /projects
 */
func GetProjects(w http.ResponseWriter, r * http.Request) {
  fmt.Fprintln(w, ":oops:")
}
