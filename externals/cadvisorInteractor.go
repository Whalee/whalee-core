package externals

import (
	"../models"
	"gopkg.in/jmcvetta/napping.v3"
	"log"
  "github.com/jmoiron/jsonq"
	"time"
	"strconv"
)

type CAInteractor struct {
	SlavesUrl []string;
}

func NewCAInteractor(slaves []string) (*CAInteractor) {
	return &CAInteractor {
    SlavesUrl: slaves,
  };
}

const base_url = "http://localhost:8080/api/v1.3/"

func (adv *CAInteractor) GetStatus(url string) models.DockerInfos {
	//TODO: remove the usage of base_url to use the slaves url in adv.
	route := base_url + url;
	log.Println("Querying " + route)
	 res := map[string]interface{}{}
	_, err := napping.Get(route, nil, &res, nil)
	if err != nil {
		log.Println("Error while napping " + route);
		log.Println(err)
	}
	jq := jsonq.NewQuery(res);
	id, _ := jq.String(url, "aliases", "1");
	stats, _ := jq.Array(url, "stats");
	cur_date, _ := jq.String(url, "stats", strconv.Itoa(len(stats)-1), "timestamp")
	cur_cpu, _ := jq.Int(url, "stats", strconv.Itoa(len(stats)-1), "cpu", "usage", "total");
	prev_date, _ := jq.String(url, "stats", strconv.Itoa(len(stats)-2), "timestamp")
	prev_cpu, _ := jq.Int(url, "stats", strconv.Itoa(len(stats)-2), "cpu", "usage", "total");
	durationInNs := getIntervalInNs(cur_date, prev_date)
	proc := models.Internals{
		Max:100,
		Cur: (float64(cur_cpu) - float64(prev_cpu)) / float64(durationInNs),
	}
	//TODO history
	//TODO memory
	//TODO hdd
	infos := models.DockerInfos{
		 Id: id,
		 Proc: proc,
	}
	return infos
}

func (adv *CAInteractor) Monitor() {
	var containers = adv.RetrieveContainers();
	for _, c := range containers {
		adv.GetStatus(c)
	}
}


func (adv *CAInteractor) RetrieveContainers() []string{
	var containers []string;
	for _, url := range adv.SlavesUrl {
		route := url + "containers/docker/"
		res := models.CAdvisorDockerList{}
		_, err := napping.Get(route, nil, &res, nil);
		if err != nil {
			log.Println(err);
		}
		log.Println(route);
		for _, container := range res.Subcontainers  {
			containers = append(containers, container.Name);
		}
	}
	return containers;
}


func getIntervalInNs(cur string, prev string) int64 {
	current_time,_ := time.Parse("2006-01-02T15:04:05.999999999Z", cur);
	previous_time,_ := time.Parse("2006-01-02T15:04:05.999999999Z", prev);

	return current_time.Sub(previous_time).Nanoseconds()
}
