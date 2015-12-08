package routes
import(
  "fmt"
  "net/http"
  "strings"
  "os/exec"
  "../models"
  "io"
  "io/ioutil"
  "log"
  "encoding/json"
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
  //Exec a command
  // exe_cmd("docker run ...")
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


func exe_cmd(cmd string) {
  fmt.Println("command is ",cmd)
  // splitting head => g++ parts => rest of the command
  parts := strings.Fields(cmd)
  head := parts[0]
  parts = parts[1:len(parts)]

  out, err := exec.Command(head,parts...).Output()
  if err != nil {
    fmt.Printf("%s", err)
  }
  fmt.Printf("%s", out)
}
