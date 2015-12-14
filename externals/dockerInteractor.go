package externals
import (
  "fmt"
  "github.com/fsouza/go-dockerclient"
  "github.com/spf13/viper"

)

type DockerInteractor struct {
  Endpoint string;
  client *docker.Client;
}

type Config struct {
  User string;
  Project string;
}

func NewRemoteInteractor(ip string, port string) (*DockerInteractor) {
  endpoint :="tcp://" + ip + ":" + port;
  dcl, err := docker.NewClient(endpoint);
  if err != nil {
    fmt.Printf("An error happened while creating the Docker Client:\n\t%s", err);
  }

  return &DockerInteractor {
    Endpoint: endpoint,
    client: dcl,
  };
}

func NewLocalInteractor(file string) (*DockerInteractor) {
  dcl, err := docker.NewClient(file);
  if err != nil {
    fmt.Printf("Docker client creation error:\n\t%s\n", err);
  }

  return &DockerInteractor {
    Endpoint: file,
    client: dcl,
  }
}
/*
 * Run a given container
 */
func (dtor *DockerInteractor) RunContainer(config Config) {
  id, err := dtor.createDefaultContainer(config);
   if err != nil {
  } else {
    dtor.startContainer(id);
    // dtor.attachLogs();
  }
}

/*
 * Create a container from a given name
 */
func (dtor *DockerInteractor) createDefaultContainer(config Config) (string, error) {

  portBindings :=  map[docker.Port][]docker.PortBinding{
        "3000/tcp": {{HostIP: "0.0.0.0", HostPort: "8080"}}}
  createContHostConfig := docker.HostConfig{
    Binds:           []string{"/var/run:/var/run", "/sys:/sys", "/var/lib/docker:/var/lib/docker"},
    PortBindings:    portBindings,
    PublishAllPorts: true,
    Privileged:      false,
}

  var createContOpts = docker.CreateContainerOptions {
    Name: config.User + "-" + config.Project,
    Config: &docker.Config {
      Image: viper.GetString("containerBase"),
      ExposedPorts: map[docker.Port]struct{} {
        "3000/tcp": {},
      },
      Env: [](string){"VIRTUAL_HOST=" + config.User + ".whalee.io/"+config.Project},
    },
    HostConfig: &createContHostConfig,
  }
  cont, err := dtor.client.CreateContainer(createContOpts)
  if err != nil {
      fmt.Printf("CreateContainer error:\n\t%s\n", err);
      return "", err
  } else {
    fmt.Printf("Container created @%s\n", cont.ID);
    return cont.ID, nil
  }

}

/*
 * Start a previously created container
 */
func (dtor *DockerInteractor) startContainer(ctid string) {
  err := dtor.client.StartContainer(ctid, nil);
  if err != nil {
    fmt.Printf("Error while starting container\n\t%s", err);
  }
  fmt.Println("Container started");
}

/*
 * Attach the logs somewhere
 */
func (dtor *DockerInteractor) attachLogs() {
  // TODO
}
