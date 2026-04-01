package store
import("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{*sql.DB}
type CachedPackage struct{ID int64 `json:"id"`;Registry string `json:"registry"`;Package string `json:"package"`;Version string `json:"version"`;SizeBytes int64 `json:"size_bytes"`;HitCount int `json:"hit_count"`;CachedAt time.Time `json:"cached_at"`}
func Open(d string)(*DB,error){os.MkdirAll(d,0755);dsn:=filepath.Join(d,"parcelproxy.db")+"?_journal_mode=WAL&_busy_timeout=5000";db,err:=sql.Open("sqlite",dsn);if err!=nil{return nil,fmt.Errorf("open: %w",err)};db.SetMaxOpenConns(1);migrate(db);return &DB{db},nil}
func migrate(db *sql.DB){db.Exec(`CREATE TABLE IF NOT EXISTS cached_packages(id INTEGER PRIMARY KEY AUTOINCREMENT,registry TEXT NOT NULL,package TEXT NOT NULL,version TEXT NOT NULL,size_bytes INTEGER DEFAULT 0,hit_count INTEGER DEFAULT 1,cached_at DATETIME DEFAULT CURRENT_TIMESTAMP,UNIQUE(registry,package,version))`)}
func(db *DB)RecordHit(registry,pkg,version string,size int64){db.Exec(`INSERT INTO cached_packages(registry,package,version,size_bytes)VALUES(?,?,?,?) ON CONFLICT(registry,package,version) DO UPDATE SET hit_count=hit_count+1`,registry,pkg,version,size)}
func(db *DB)List()([]CachedPackage,error){rows,_:=db.Query(`SELECT id,registry,package,version,size_bytes,hit_count,cached_at FROM cached_packages ORDER BY hit_count DESC LIMIT 500`);defer rows.Close();var out[]CachedPackage;for rows.Next(){var c CachedPackage;rows.Scan(&c.ID,&c.Registry,&c.Package,&c.Version,&c.SizeBytes,&c.HitCount,&c.CachedAt);out=append(out,c)};return out,nil}
func(db *DB)Delete(id int64){db.Exec(`DELETE FROM cached_packages WHERE id=?`,id)}
func(db *DB)Stats()(map[string]interface{},error){var total int;var size int64;db.QueryRow(`SELECT COUNT(*),COALESCE(SUM(size_bytes),0) FROM cached_packages`).Scan(&total,&size);return map[string]interface{}{"cached_packages":total,"total_bytes":size},nil}
