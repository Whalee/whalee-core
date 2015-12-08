package main

import(
  "gopkg.in/jmcvetta/napping.v3"
  "./models"
  "log"
)
const base_url= "http://localhost:8080/api/v1.3/"

func getStatus(dockerId string) models.CAdvisorDocker {
  route:=base_url +"containers/docker/" + dockerId
  log.Println("Querying " + route);
  res:= models.CAdvisorDocker{}
  _, err:= napping.Get(route, nil, &res, nil);
  if(err!= nil) {
    log.Println(err);
  }
  //TODO: change the format to a DockerInfos

  return res;
}
