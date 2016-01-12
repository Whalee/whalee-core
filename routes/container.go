package routes
import(
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
  "../externals"
  "github.com/spf13/viper"
)


func GetInfos(w http.ResponseWriter, r *http.Request) {
  var cad *externals.CAInteractor
  vars := mux.Vars(r)
  id := vars["id"]
  urls := viper.GetStringSlice("cadvisorUrl")
  cad = externals.NewCAInteractor(urls)
  res:= cad.GetStatus("/docker/"+id)
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(w).Encode(res); err != nil {
      panic(err)
  }
}


func monitor(w http.ResponseWriter, r *http.Request) {
  var cad *externals.CAInteractor
  fmt.Fprintf(w, "ok");
}
