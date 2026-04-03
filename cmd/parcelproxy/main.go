package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-parcelproxy/internal/server";"github.com/stockyard-dev/stockyard-parcelproxy/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./parcelproxy-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("parcelproxy: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Parcel Proxy — Self-hosted package registry proxy\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("parcelproxy: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
