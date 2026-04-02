package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Upstream struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Registry string `json:"registry"`
	URL string `json:"url"`
	CacheEnabled int `json:"cache_enabled"`
	CacheHits int `json:"cache_hits"`
	CacheMisses int `json:"cache_misses"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"parcelproxy.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS upstreams(id TEXT PRIMARY KEY,name TEXT NOT NULL,registry TEXT DEFAULT '',url TEXT DEFAULT '',cache_enabled INTEGER DEFAULT 1,cache_hits INTEGER DEFAULT 0,cache_misses INTEGER DEFAULT 0,status TEXT DEFAULT 'active',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Upstream)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO upstreams(id,name,registry,url,cache_enabled,cache_hits,cache_misses,status,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Registry,e.URL,e.CacheEnabled,e.CacheHits,e.CacheMisses,e.Status,e.CreatedAt);return err}
func(d *DB)Get(id string)*Upstream{var e Upstream;if d.db.QueryRow(`SELECT id,name,registry,url,cache_enabled,cache_hits,cache_misses,status,created_at FROM upstreams WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Registry,&e.URL,&e.CacheEnabled,&e.CacheHits,&e.CacheMisses,&e.Status,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Upstream{rows,_:=d.db.Query(`SELECT id,name,registry,url,cache_enabled,cache_hits,cache_misses,status,created_at FROM upstreams ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Upstream;for rows.Next(){var e Upstream;rows.Scan(&e.ID,&e.Name,&e.Registry,&e.URL,&e.CacheEnabled,&e.CacheHits,&e.CacheMisses,&e.Status,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Upstream)error{_,err:=d.db.Exec(`UPDATE upstreams SET name=?,registry=?,url=?,cache_enabled=?,cache_hits=?,cache_misses=?,status=? WHERE id=?`,e.Name,e.Registry,e.URL,e.CacheEnabled,e.CacheHits,e.CacheMisses,e.Status,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM upstreams WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM upstreams`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Upstream{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,registry,url,cache_enabled,cache_hits,cache_misses,status,created_at FROM upstreams WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Upstream;for rows.Next(){var e Upstream;rows.Scan(&e.ID,&e.Name,&e.Registry,&e.URL,&e.CacheEnabled,&e.CacheHits,&e.CacheMisses,&e.Status,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM upstreams GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
