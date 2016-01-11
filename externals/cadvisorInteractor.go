package externals

import (
	"../models"
	"gopkg.in/jmcvetta/napping.v3"
	"log"
  "../jsonq"
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
	cur_cpu := extractCpuFromStats(jq, url, len(stats)-1);
	prev_cpu := extractCpuFromStats(jq, url, len(stats)-2);
	cur_mem := extractMemFromStats(jq, url, len(stats)-1);
	limit_mem, err := jq.Float64(url, "spec", "memory", "limit");
	proc := models.Internals{
		Max:100,
		Cur: getCpuUsage(cur_cpu, prev_cpu),
	}

	mem := models.Internals{
		Max: 100,
		Cur: float64(cur_mem.Value),
	}

	var p_hist, m_hist []float64
	for i:= maxInt(len(stats)-100,1); i<len(stats); i++ {
		p_cur := extractCpuFromStats(jq, url, i);
		p_prev := extractCpuFromStats(jq, url, i-1);
		p_hist = append(p_hist, getCpuUsage(p_cur, p_prev))
		m_cur := extractMemFromStats(jq,url,i)
		m_hist = append(m_hist, getMemUsage(m_cur.Value, limit_mem));
	}
	proc.Hist = p_hist
	mem.Hist = m_hist

	//TODO hdd
	infos := models.DockerInfos{
		 Id: id,
		 Proc: proc,
		 Memory:mem,
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

func getCpuUsage(cur models.TimestampedValue, prev models.TimestampedValue) float64 {
	durationInNs := getIntervalInNs(cur.Date, prev.Date)
	return (float64(cur.Value) - float64(prev.Value)) / float64(durationInNs)
}
func getMemUsage(cur int64, max float64) float64 {
	res := float64(float64(cur) / max) * 100
	return res
}

func extractCpuFromStats(jq *jsonq.JsonQuery, url string, i int) models.TimestampedValue {
	date, _ := jq.String(url, "stats", strconv.Itoa(i), "timestamp")
	cpu, _ := jq.Int(url, "stats", strconv.Itoa(i), "cpu", "usage", "total");
	return models.TimestampedValue{
		Value: int64(cpu),
		Date: date,
	}
}

func extractMemFromStats(jq *jsonq.JsonQuery, url string, i int) models.TimestampedValue {
	date, _ := jq.String(url, "stats", strconv.Itoa(i), "timestamp")
	mem, _ := jq.Int64(url, "stats", strconv.Itoa(i), "memory", "usage", );
	return models.TimestampedValue{
		Value: mem,
		Date: date,
	}}

func maxInt(i1 int, i2 int) int {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}
