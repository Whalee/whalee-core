package routes
import(
  "fmt"
  "net/http"
  // "github.com/gorilla/mux"
  "../externals"
  "github.com/spf13/viper"
)


func GetInfos(w http.ResponseWriter, r *http.Request) {
  var cad *externals.CAInteractor
  // vars := mux.Vars(r)
  // id := vars["id"]
  urls := viper.GetStringSlice("cadvisorUrl")
  cad = externals.NewCAInteractor(urls)
  dockers := cad.RetrieveContainers()
  res:= cad.GetStatus(dockers[0])
  fmt.Println(res)
}
