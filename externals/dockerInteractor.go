package externals
import (
  "fmt"
  "github.com/fsouza/go-dockerclient"
  "github.com/spf13/viper"
  "time"
  "gopkg.in/jmcvetta/napping.v3"
  	"log"
    "strconv"
    "../jsonq"
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
func (dtor *DockerInteractor) RunContainer(config Config) (string, string, string, error) {
  id, err := dtor.createDefaultContainer(config);
   if err != nil {
     fmt.Printf("Error while creating default container\n\t%s", err)
     return "","", "", err
  } else {
    dtor.startContainer(id);
    appPort, managerPort, ip, err := dtor.retrieveExposedPort(id, config.User + "@" + config.Project);
    if err != nil {
      fmt.Printf("Error while starting the container\n\t%s", err)
    }
    //Wait till container started
    //TODO: check for a better way to do that.
    time.Sleep(5 * time.Second)
    return appPort, managerPort, ip, err
  }
}

func (dtor *DockerInteractor) ListContainers(project, user string) {
  opts :=  docker.ListContainersOptions{
    Filters: map[string][]string{
      "label":{"project=" + project, "user=" + user,},
    },
  }
  containers, err := dtor.client.ListContainers(opts)
  if err != nil {
    fmt.Printf("Error while listing containers\n\t%s", err)
  }
  fmt.Println(containers);
}
func (dtor *DockerInteractor) StartDRCoN(project Config, consulip, consulport string) {
  contid, err := dtor.createDRCoNContainer(project, consulip, consulport);
  if err != nil {
    // fmt.Printf("Error while creating default container\n\t%s", err)
    // return "","", err
    return
  } else {
   dtor.startContainer(contid);
  }
}

func (dtor *DockerInteractor) createDRCoNContainer(conf Config, consulip, consulport string) (string, error) {
  service_name :=  conf.User + "@"+conf.Project
  createContOpts := docker.CreateContainerOptions {
    Name: "DRCoN_"+conf.User + "_"+conf.Project,
    Config: &docker.Config {
      Image: "apox0/drcon",
      ExposedPorts: map[docker.Port]struct{} {
        "80/tcp":{},
      },
      Env: []string {
        "VIRTUAL_HOST="+conf.User+".whalee.io",
        "constraint:function==master",
        "CONSUL="+consulip+":"+consulport,
        "SERVICE_NAME=DRCoN_" + service_name,
        "SERVICE="+service_name,
      },
      Labels: map[string]string {
        "project": conf.Project,
        "user": conf.User,
      },
    },
    HostConfig: &docker.HostConfig{
      Binds: []string{"/var/run:/var/run", "/sys:/sys", "/var/lib/docker:/var/lib/docker"},
    },
  }
  cont, err := dtor.client.CreateContainer(createContOpts)
  if err != nil {
    fmt.Printf("CreateContainer failed for DRCoN\n\t%s\n",err)
    return "",err
  } else {
    fmt.Printf("DRCoN Container created\n");
    return cont.ID,nil
  }
}
/*
 * Create a container from a given name
 */
func (dtor *DockerInteractor) createDefaultContainer(config Config) (string, error) {

  portBindings :=  map[docker.Port][]docker.PortBinding{
     "80/tcp": {{HostIP: "0.0.0.0", HostPort: "0"}},
    "8081/tcp": {{HostIP: "0.0.0.0", HostPort: "0"}},
  }
  createContHostConfig := docker.HostConfig{
    Binds:           []string{"/var/run:/var/run", "/sys:/sys", "/var/lib/docker:/var/lib/docker"},
    PortBindings:    portBindings,
    PublishAllPorts: true,
    Privileged:      false,
}

  var createContOpts = docker.CreateContainerOptions {
    Config: &docker.Config {
      Image: viper.GetString("containerBase"),
      ExposedPorts: map[docker.Port]struct{} {
        "80/tcp":{},
        "8081/tcp": {},
      },
      Env: [](string){
        "SERVICE_NAME=" + config.User + "@"+config.Project,
        "SERVICE_8081_IGNORE=1",
        "SERVICE_8080_IGNORE=1",
        "constraint:function==node"},
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

func (dtor *DockerInteractor) retrieveExposedPort(ctid, service_name string) (string, string, string, error) {
  cont, err :=dtor.client.InspectContainer(ctid);
  if err != nil {
    fmt.Printf("Error while Inspecting container \n\t%s", err)
    return "", "", "", err
  }
  port1 := cont.NetworkSettings.Ports["80/tcp"][0].HostPort
  managerPort :=cont.NetworkSettings.Ports["8081/tcp"][0].HostPort
  ip := "localhost"
  if viper.IsSet("consul") {
    route := "http://" + viper.GetString("consul.ip") + ":"+  viper.GetString("consul.port") + "/v1/catalog/service/" + service_name
    res := []interface{}{}
    _,err := napping.Get(route, nil, &res, nil);
    if err != nil {
      log.Println("Error while napping " + route);
      log.Println(err)
    }
    for  i := 0; i < len(res) ; i++ {
      jq := jsonq.NewQuery(res[i]);
      port, _ := jq.Int("ServicePort")
      if(strconv.Itoi(port) == port1) {
        ipjq, _ := jq.String("Address")
        ip =ipjq
      }
    }
  }
  fmt.Printf("Two interesting ports: 3000 -> %s, 8081 -> %s\n",port1,managerPort)
  return port1, managerPort, ip, nil
}


/*
 * Attach the logs somewhere
 */
func (dtor *DockerInteractor) attachLogs() {
  // TODO
}
