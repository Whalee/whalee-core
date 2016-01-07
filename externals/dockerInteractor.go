package externals
import (
  "fmt"
  "github.com/fsouza/go-dockerclient"
  "github.com/spf13/viper"
  "time"
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
  dcl, err   := docker.NewClient(file);
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
func (dtor *DockerInteractor) RunContainer(config Config) (string, string, error) {
  id, err := dtor.createDefaultContainer(config);
   if err != nil {
     fmt.Printf("Error while creating default container\n\t%s", err)
     return "","", err
  } else {
    dtor.startContainer(id);
    appPort, managerPort, err := dtor.retrieveExposedPort(id);
    if err != nil {
      fmt.Printf("Error while starting the container\n\t%s", err)
    }
    //Wait till container started
    //TODO: check for a better way to do that.
    time.Sleep(5 * time.Second)
    return appPort, managerPort, err
  }
}

/*
 * Create a container from a given name
 */
func (dtor *DockerInteractor) createDefaultContainer(config Config) (string, error) {

  portBindings :=  map[docker.Port][]docker.PortBinding{
        "3000/tcp": {{HostIP: "0.0.0.0", HostPort: "0"}},
        "8081/tcp": {{HostIP: "0.0.0.0", HostPort: "0"}}}
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
        "8081/tcp": {},
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

func (dtor *DockerInteractor) retrieveExposedPort(ctid string) (string, string, error) {
  cont, err :=dtor.client.InspectContainer(ctid);
  if err != nil {
    fmt.Printf("Error while Inspecting container \n\t%s", err)
    return "", "", err
  }
  port1 := cont.NetworkSettings.Ports["3000/tcp"][0].HostPort
  managerPort :=cont.NetworkSettings.Ports["8081/tcp"][0].HostPort
  fmt.Printf("Two interesting ports: 3000 -> %s, 8081 -> %s\n",port1,managerPort)
  return port1, managerPort, nil
}

/*
 * Attach the logs somewhere
 */
func (dtor *DockerInteractor) attachLogs() {
  // TODO
}
