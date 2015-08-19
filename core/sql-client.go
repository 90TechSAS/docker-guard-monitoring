package core

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"

	"../utils"
)

/*
	Container info
*/
type Container struct {
	ID         string
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
	ContainerID string
	Time        int64
	SizeRootFs  uint64
	SizeRw      uint64
	SizeMemory  uint64
	Running     bool
}

/*
	SQL Variables
*/
var (
	// DB
	DB       *sql.DB
	ProbesID map[string]int

	// Prepared queries
	GetProbeIDStmt            *sql.Stmt
	InsertProbeStmt           *sql.Stmt
	InsertContainerStmt       *sql.Stmt
	DeleteContainerStmt       *sql.Stmt
	InsertStatStmt            *sql.Stmt
	DeleteStatStmt            *sql.Stmt
	GetContainerByIdStmt      *sql.Stmt
	GetStatsByContainerIdStmt *sql.Stmt
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
	InsertContainerStmt, err = DB.Prepare("INSERT INTO containers VALUES (?,?,?,?,?,?)")
	if err != nil {
		l.Critical("Can't create InsertContainerStmt:", err)
	}
	DeleteContainerStmt, err = DB.Prepare("DELETE FROM containers WHERE id=?")
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
	GetContainerByIdStmt, err = DB.Prepare("SELECT * FROM containers WHERE id=?")
	if err != nil {
		l.Critical("Can't create GetContainerByIdStmt:", err)
	}
	GetStatsByContainerIdStmt, err = DB.Prepare("SELECT * FROM stats WHERE containerid=?")
	if err != nil {
		l.Critical("Can't create DeleteStatStmt:", err)
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
func (c *Container) Insert() error {
	var err error // Error handling

	_, err = InsertContainerStmt.Exec(c.ID, c.ProbeID, c.Hostname, c.Image, c.IPAddress, c.MacAddress)

	return err
}

/*
	Delete a container
*/
func (c *Container) Delete() error {
	var err error // Error handling

	_, err = DeleteContainerStmt.Exec(c.ID)

	return err
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
	for _, i := range []string{"id", "hostname", "image", "ip", "mac"} {
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
		l.Error("GetContainersBy:", err)
		return containers, err
	}
	defer rows.Close()
	for rows.Next() {
		var tmpContainer Container

		if err = rows.Scan(&tmpContainer.ID,
			&tmpContainer.ProbeID,
			&tmpContainer.Hostname,
			&tmpContainer.Image,
			&tmpContainer.IPAddress,
			&tmpContainer.MacAddress); err != nil {
			l.Error("GetContainersBy:", err)
			return containers, err
		}

		containers = append(containers, tmpContainer)
	}
	if err = rows.Err(); err != nil {
		l.Error("GetContainersBy:", err)
		return containers, err
	}

	return containers, nil
}

/*
	Get a container by id
*/
func GetContainerById(id string) (Container, error) {
	var container Container // Container to return
	var err error           // Error handling

	err = GetContainerByIdStmt.QueryRow(id).Scan(&container.ID, &container.ProbeID, &container.Hostname, &container.Image, &container.IPAddress, &container.MacAddress)
	if err != nil {
		l.Error("GetContainerById:", err)
		return container, err
	}

	return container, nil
}

/*
	Insert a stat
*/
func (s *Stat) Insert() error {
	var err error // Error handling

	_, err = InsertStatStmt.Exec(s.ContainerID, s.Time, s.SizeRootFs, s.SizeRw, s.SizeMemory, s.Running)

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
func GetStatsByContainerID(containerID string) ([]Stat, error) {
	var stats []Stat   // List of stats to return
	var err error      // Erro handling
	var tmpStat Stat   // Temporary stat
	var rows *sql.Rows // Temporary sql rows

	rows, err = GetStatsByContainerIdStmt.Query(containerID)
	if err != nil {
		l.Error("GetStatsByContainerID:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&tmpStat.ContainerID, &tmpStat.Time, &tmpStat.SizeRootFs, &tmpStat.SizeRw, &tmpStat.SizeMemory, &tmpStat.Running)
		if err != nil {
			l.Error("GetStatsByContainerID:", err)
			return nil, err
		}
		stats = append(stats, tmpStat)
	}
	err = rows.Err()
	if err != nil {
		l.Error("GetStatsByContainerID:", err)
		return nil, err
	}

	return stats, nil
}
