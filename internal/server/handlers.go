package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-parcelproxy/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.CachedPackage{}};writeJSON(w,200,list)}
func(s *Server)handleRecord(w http.ResponseWriter,r *http.Request){var req struct{Registry string `json:"registry"`;Package string `json:"package"`;Version string `json:"version"`;SizeBytes int64 `json:"size_bytes"`};json.NewDecoder(r.Body).Decode(&req);if req.Package==""{writeError(w,400,"package required");return};if req.Registry==""{req.Registry="npm"};s.db.RecordHit(req.Registry,req.Package,req.Version,req.SizeBytes);writeJSON(w,201,map[string]string{"status":"recorded"})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
