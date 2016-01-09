package externals

import (
	"../models"
	"gopkg.in/jmcvetta/napping.v3"
	"log"
)

type CadvisorInteractor struct {
	SlavesUrl []string;
}

const base_url = "http://localhost:8080/api/v1.3/"

func (adv *CadvisorInteractor) GetStatus(url string) models.DockerInfos {
	route := base_url + url;
	log.Println("Querying " + route)
	res := models.CAdvisorDocker{}
	_, err := napping.Get(route, nil, &res, nil)
	if err != nil {
		log.Println(err)
	}
	//TODO: change the format to a DockerInfos

	// proc := models.Internals{
	// 	Max:100,
	// }
	infos := models.DockerInfos{
		Id: res.Aliases[len(res.Aliases)-1],
	}
	return infos
}

func (adv *CadvisorInteractor) Monitor() {
	var containers = adv.RetrieveContainers();
	for _, c := range containers {
		adv.GetStatus(c);
	}
}


func (adv *CadvisorInteractor) RetrieveContainers() []string{
	var containers []string;
	for _, url := range adv.SlavesUrl {
		route := url + "containers/docker/"
		res := models.CAdvisorDockerList{}
		_, err := napping.Get(route, nil, &res, nil);
		if err != nil {
			log.Println(err);
		}
		for _, container := range res.Subcontainers  {
			containers = append(containers, container.Name);
		}
	}
	return containers;
}
