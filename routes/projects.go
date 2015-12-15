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
  "github.com/spf13/viper"
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
  if project.User != "" && project.Project != "" {
    log.Println("Creating a docker for "+  project.Project);
    port := startDocker(project.User, project.Project)
    deployApp(port, project.User, project.Project)
    //Call inside the executed docker the set up server

    fmt.Fprintln(w, "ok")
  }
}


func deployApp(managerPort string, user string, project string) {
  log.Println("Retrieving github " + project + "/" + user + " from docker");
  url := "http://localhost:" + managerPort + "/?url='http://github.com/"+ user + "/" + project + "'&action=clone";
  req, err := http.NewRequest("GET", url, nil)
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()
  fmt.Println("response Status:", resp.Status)
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))
  //
  // s := napping.Session{}
  // log.Printf("Url -> %s\n", url)
  // ans :=""
  // _, err := s.Get(url, nil, &ans, nil);
  // if err != nil {
  //   fmt.Printf("Error !!\n\t%s\n",err);
  // }
  // if resp.Status() == 200 {
  //   log.Printf("Everything's alright mate\n")
  //   return;
  // } else {
  //   log.Printf("Error while deploying user app\n\t%s",err);
  // }
}

func startDocker(user string, project string) (string) {
  var dockerClient *externals.DockerInteractor
  if viper.IsSet("dockerRemote") {
    dockerClient = externals.NewRemoteInteractor(viper.GetString("dockerRemote.ip"), viper.GetString("dockerRemote.port"));
  } else {
    dockerClient = externals.NewLocalInteractor("unix:///var/run/docker.sock");
  }
  config := externals.Config {
    User: user,
    Project: project,
  }
  _, managerPort, _ := dockerClient.RunContainer(config);
  return managerPort
}

/*
 * GET /projects
 */
func GetProjects(w http.ResponseWriter, r * http.Request) {
  fmt.Fprintln(w, ":oops:")
}
