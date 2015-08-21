package core

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"../utils"
)

/*
	Container info
*/
type Container struct {
	ID         int
	CID        string
	ProbeID    int
	Hostname   string
	Image      string
	IPAddress  string
	MacAddress string
}

/*
	Container's stats
*/
type Stat struct {
	ContainerID int
	Time        int64
	SizeRootFs  uint64
	SizeRw      uint64
	SizeMemory  uint64
	Running     bool
}

/*

 */
type Options struct {
	Since  int
	Before int
	Limit  int
}

/*
	SQL Variables
*/
var (
	// DB
	DB       *sql.DB
	ProbesID map[string]int

	// Prepared queries
	GetProbeIDStmt             *sql.Stmt
	InsertProbeStmt            *sql.Stmt
	InsertContainerStmt        *sql.Stmt
	DeleteContainerStmt        *sql.Stmt
	GetLastStatStmt            *sql.Stmt
	GetBetweenStatsStmt        *sql.Stmt
	InsertStatStmt             *sql.Stmt
	DeleteStatStmt             *sql.Stmt
	GetContainerByCIDStmt      *sql.Stmt
	GetStatsByContainerCIDStmt *sql.Stmt
)

func InitSQL() {
	var err error // Error handling

	// Connect DB
	sqlc := DGConfig.DockerGuard.SQLServer // SQL Config
	DB, err = sql.Open("mysql", sqlc.User+":"+sqlc.Pass+"@tcp("+sqlc.IP+":"+utils.I2S(sqlc.Port)+")/"+sqlc.DB)
	if err != nil {
		l.Critical(err)
	}
	err = DB.Ping()
	if err != nil {
		l.Critical(err)
	}
	l.Verbose("Connected to SQL database")

	// Prepare DB queries
	GetProbeIDStmt, err = DB.Prepare("SELECT id FROM probes WHERE name=?")
	if err != nil {
		l.Critical("Can't create InsertContainerStmt:", err)
	}
	InsertProbeStmt, err = DB.Prepare("INSERT INTO probes VALUES (DEFAULT,?)")
	if err != nil {
		l.Critical("Can't create InsertContainerStmt:", err)
	}
	InsertContainerStmt, err = DB.Prepare("INSERT INTO containers VALUES (DEFAULT,?,?,?,?,?,?)")
	if err != nil {
		l.Critical("Can't create InsertContainerStmt:", err)
	}
	DeleteContainerStmt, err = DB.Prepare("DELETE FROM containers WHERE containerid=?")
	if err != nil {
		l.Critical("Can't create DeleteContainerStmt:", err)
	}
	GetLastStatStmt, err = DB.Prepare("SELECT * FROM stats WHERE containerid=? ORDER BY time DESC LIMIT 1")
	if err != nil {
		l.Critical("Can't create DeleteContainerStmt:", err)
	}
	GetBetweenStatsStmt, err = DB.Prepare("SELECT * FROM stats WHERE containerid=? AND time BETWEEN ? AND ? ORDER BY time")
	if err != nil {
		l.Critical("Can't create DeleteContainerStmt:", err)
	}
	InsertStatStmt, err = DB.Prepare("INSERT INTO stats VALUES (?,?,?,?,?,?)")
	if err != nil {
		l.Critical("Can't create InsertStatStmt:", err)
	}
	DeleteStatStmt, err = DB.Prepare("DELETE FROM stats WHERE containerid=? AND time=?")
	if err != nil {
		l.Critical("Can't create DeleteStatStmt:", err)
	}
	GetContainerByCIDStmt, err = DB.Prepare("SELECT * FROM containers WHERE containerid=?")
	if err != nil {
		l.Critical("Can't create GetContainerByIdStmt:", err)
	}

	// Get probes ID
	ProbesID = make(map[string]int)
	for _, probe := range DGConfig.Probes {
		id, err := GetProbeID(probe.Name)
		if err != nil {
			l.Critical("Error GetProbeID ("+probe.Name+"):", err)
		}
		ProbesID[probe.Name] = id
	}
}

/*
	Parse Options
*/
func GetOptions(r *http.Request) Options {
	var options Options // Returned options
	var err error       // Error handling

	// Get url parameters
	oS := r.URL.Query().Get("since")
	oB := r.URL.Query().Get("before")
	oL := r.URL.Query().Get("limit")

	// Format parameters to int and set options
	oSInt, err := utils.S2I(oS)
	if err != nil {
		options.Since = -1
	} else {
		options.Since = oSInt
	}
	oBInt, err := utils.S2I(oB)
	if err != nil {
		options.Before = -1
	} else {
		options.Before = oBInt
	}
	oLInt, err := utils.S2I(oL)
	if err != nil {
		options.Limit = -1
	} else {
		options.Limit = oLInt
	}

	return options
}

/*
	Get a probe ID from sql server
	If does not exist: create it and return the probe ID
*/
func GetProbeID(name string) (int, error) {
	var id int64  // Probe ID to return
	var err error // Error handling

	err = GetProbeIDStmt.QueryRow(name).Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			result, err := InsertProbeStmt.Exec(name)
			if err != nil {
				l.Error("GetProbeID:", err)
				return 0, err
			}
			id, err = result.LastInsertId()
			if err != nil {
				l.Error("GetProbeID:", err)
				return 0, err
			}
		} else {
			l.Error("GetProbeID:", err)
			return 0, err
		}
	}

	return int(id), nil
}

/*
	Insert a container
*/
func (c *Container) Insert() (int64, error) {
	var err error         // Error handling
	var result sql.Result // SQL result

	result, err = InsertContainerStmt.Exec(c.CID, c.ProbeID, c.Hostname, c.Image, c.IPAddress, c.MacAddress)

	if err == nil {
		id, err := result.LastInsertId()
		if err != nil {
			return 0, err
		} else {
			return id, nil
		}
	}

	return 0, err
}

/*
	Delete a container
*/
func (c *Container) Delete() error {
	var err error // Error handling

	_, err = DeleteContainerStmt.Exec(c.CID)

	return err
}

/*
	Get container's last stat
*/
func (c *Container) GetLastStat() (Stat, error) {
	var stat Stat // Returned stat
	var err error // Error handling

	err = GetLastStatStmt.QueryRow(c.ID).Scan(&stat.ContainerID, &stat.Time, &stat.SizeRootFs, &stat.SizeRw, &stat.SizeMemory, &stat.Running)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			l.Error("GetLastStat:", err)
		}
		return stat, err
	}

	return stat, err
}

/*
	Get container's stats between two dates
*/
func (c *Container) GetBetweenStats(begin, end int) ([]Stat, error) {
	var stats []Stat   // Returned stats
	var rows *sql.Rows // SQL Rows
	var err error      // Error handling

	rows, err = GetBetweenStatsStmt.Query(c.ID, begin, end)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			l.Error("GetBetweenStats:", err)
		}
		return stats, err
	}

	defer rows.Close()
	for rows.Next() {
		var tmpStat Stat

		if err = rows.Scan(
			&tmpStat.ContainerID,
			&tmpStat.Time,
			&tmpStat.SizeRootFs,
			&tmpStat.SizeRw,
			&tmpStat.SizeMemory,
			&tmpStat.Running); err != nil {
			l.Error("GetBetweenStats: Can't scan row:", err)
			return stats, err
		}

		stats = append(stats, tmpStat)
	}
	if err = rows.Err(); err != nil {
		l.Error("GetBetweenStatStmt: Rows error:", err)
		return stats, err
	}

	return stats, err
}

/*
	Get containers by $field
*/
func GetContainersBy(field string, value interface{}) ([]Container, error) {
	var containers []Container // Containers to return
	var rows *sql.Rows         // SQL Rows
	var err error              // Error handling

	// Protection against SQL injection
	var fieldExists bool = false
	for _, i := range []string{"id", "containerid", "probeid", "hostname", "image", "ip", "mac"} {
		if field == i {
			fieldExists = true
		}
	}
	if !fieldExists {
		l.Error("GetContainersBy: Field (" + field + ") does not exist.")
		return containers, errors.New("GetContainersBy: Field (" + field + ") does not exist.")
	}

	// Get containers
	rows, err = DB.Query("SELECT * FROM containers WHERE "+field+"=?", value)
	if err != nil {
		l.Error("GetContainersBy: Can't get rows:", err)
		return containers, err
	}
	defer rows.Close()
	for rows.Next() {
		var tmpContainer Container

		if err = rows.Scan(&tmpContainer.ID,
			&tmpContainer.CID,
			&tmpContainer.ProbeID,
			&tmpContainer.Hostname,
			&tmpContainer.Image,
			&tmpContainer.IPAddress,
			&tmpContainer.MacAddress); err != nil {
			l.Error("GetContainersBy: Can't scan row:", err)
			return containers, err
		}

		containers = append(containers, tmpContainer)
	}
	if err = rows.Err(); err != nil {
		l.Error("GetContainersBy: Rows error:", err)
		return containers, err
	}

	return containers, nil
}

/*
	Get a container by id
*/
func GetContainerByID(id int) (Container, error) {
	var container Container // Container to return
	var err error           // Error handling

	err = GetContainerByCIDStmt.QueryRow(id).Scan(&container.ID, &container.CID, &container.ProbeID, &container.Hostname, &container.Image, &container.IPAddress, &container.MacAddress)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return container, err
		} else {
			l.Error("GetContainerByID:", err)
			return container, err
		}
	}

	return container, nil
}

/*
	Get a container by cid
*/
func GetContainerByCID(cid string) (Container, error) {
	var container Container // Container to return
	var err error           // Error handling

	err = GetContainerByCIDStmt.QueryRow(cid).Scan(&container.ID, &container.CID, &container.ProbeID, &container.Hostname, &container.Image, &container.IPAddress, &container.MacAddress)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return container, err
		} else {
			l.Error("GetContainerByCID:", err)
			return container, err
		}
	}

	return container, nil
}

/*
	Insert a stat
*/
func (s *Stat) Insert() error {
	var err error // Error handling
	var timeInsert int64 = time.Now().Unix()

	_, err = InsertStatStmt.Exec(s.ContainerID, timeInsert, s.SizeRootFs, s.SizeRw, s.SizeMemory, s.Running)

	return err
}

/*
	Delete a stat
*/
func (s *Stat) Delete() error {
	var err error // Error handling

	_, err = DeleteStatStmt.Exec(s.ContainerID, s.Time)

	return err
}

/*
	Get stats by container id
*/
func GetStatsByContainerCID(containerCID string, o Options) ([]Stat, error) {
	var stats []Stat    // List of stats to return
	var err error       // Erro handling
	var tmpStat Stat    // Temporary stat
	var rows *sql.Rows  // Temporary sql rows
	var sqlQuery string // SQL query
	var oS, oB string   // SQL options

	sqlQuery = "SELECT * FROM stats WHERE containerid=(SELECT id FROM containers WHERE containerid=?)" // Base sql query

	// Add options
	if o.Since != -1 || o.Before != -1 {
		if o.Since != -1 && o.Before != -1 {
			oS = fmt.Sprintf("%d", o.Since)
			oB = fmt.Sprintf("%d", o.Before)
		} else if o.Since == -1 || o.Before != -1 {
			oS = fmt.Sprintf("%d", 0)
			oB = fmt.Sprintf("%d", o.Before)
		} else if o.Since != -1 || o.Before == -1 {
			oS = fmt.Sprintf("%d", o.Since)
			oB = fmt.Sprintf("%d", 2000000000)
		}
		sqlQuery += fmt.Sprintf(" AND time BETWEEN %s AND %s", oS, oB)
	}
	if o.Limit != -1 {
		sqlQuery += fmt.Sprintf(" LIMIT %d", o.Limit)
	}

	// Exec query
	rows, err = DB.Query(sqlQuery, containerCID)
	if err != nil {
		l.Error("GetStatsByContainerCID:", err)
		return nil, err
	}
	defer rows.Close()

	// Get results
	for rows.Next() {
		err = rows.Scan(&tmpStat.ContainerID, &tmpStat.Time, &tmpStat.SizeRootFs, &tmpStat.SizeRw, &tmpStat.SizeMemory, &tmpStat.Running)
		if err != nil {
			l.Error("GetStatsByContainerCID:", err)
			return nil, err
		}
		stats = append(stats, tmpStat)
	}
	err = rows.Err()
	if err != nil {
		l.Error("GetStatsByContainerCID:", err)
		return nil, err
	}

	return stats, nil
}

/*
	Get stats by probe id
*/
func GetStatsByContainerProbeID(probeID string, o Options) ([]Stat, error) {
	var stats []Stat       // List of stats to return
	var err error          // Erro handling
	var containerIDs []int // Array of container ID
	var tmpStat Stat       // Temporary stat
	var rows *sql.Rows     // Temporary sql rows
	var sqlQuery string    // SQL query
	var oS, oB string      // SQL options

	sqlQuery = "SELECT * FROM stats WHERE containerid=?" // Base sql query

	// Add options
	if o.Since != -1 || o.Before != -1 {
		if o.Since != -1 && o.Before != -1 {
			oS = fmt.Sprintf("%d", o.Since)
			oB = fmt.Sprintf("%d", o.Before)
		} else if o.Since == -1 || o.Before != -1 {
			oS = fmt.Sprintf("%d", 0)
			oB = fmt.Sprintf("%d", o.Before)
		} else if o.Since != -1 || o.Before == -1 {
			oS = fmt.Sprintf("%d", o.Since)
			oB = fmt.Sprintf("%d", 2000000000)
		}
		sqlQuery += fmt.Sprintf(" AND time BETWEEN %s AND %s", oS, oB)
	}
	if o.Limit != -1 {
		sqlQuery += fmt.Sprintf(" LIMIT %d", o.Limit)
	}

	// Get containers' id
	rows, err = DB.Query("SELECT id FROM containers WHERE probeid=?", probeID)
	if err != nil {
		l.Error("GetStatsByContainerProbeID:", err)
		return nil, err
	}
	defer rows.Close()

	// Get results
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			l.Error("GetStatsByContainerProbeID:", err)
			return nil, err
		}
		containerIDs = append(containerIDs, id)
	}
	err = rows.Err()
	if err != nil {
		l.Error("GetStatsByContainerProbeID:", err)
		return nil, err
	}

	// Get containers' stats
	for _, id := range containerIDs {
		// Exec query
		rows, err = DB.Query(sqlQuery, id)
		if err != nil {
			l.Error("GetStatsByContainerProbeID:", err)
			return nil, err
		}
		defer rows.Close()

		// Get results
		for rows.Next() {
			err = rows.Scan(&tmpStat.ContainerID, &tmpStat.Time, &tmpStat.SizeRootFs, &tmpStat.SizeRw, &tmpStat.SizeMemory, &tmpStat.Running)
			if err != nil {
				l.Error("GetStatsByContainerProbeID:", err)
				return nil, err
			}
			stats = append(stats, tmpStat)
		}
		err = rows.Err()
		if err != nil {
			l.Error("GetStatsByContainerProbeID:", err)
			return nil, err
		}
	}

	return stats, nil
}
