package core

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

/*
	Container info
*/
type Container struct {
	ID         string
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
	DB *sql.DB // DB

	// Prepared queries
	InsertContainerStmt  *sql.Stmt
	DeleteContainerStmt  *sql.Stmt
	InsertStatStmt       *sql.Stmt
	DeleteStatStmt       *sql.Stmt
	GetContainerByIdStmt *sql.Stmt
)

func InitSQL() {
	var err error // Error handling

	// Connect DB
	DB, err = sql.Open("mysql", "root:toor@tcp(172.17.0.4:3306)/dgs")
	if err != nil {
		l.Critical(err)
	}
	l.Verbose("Connected to SQL database")

	// Prepare DB queries
	InsertContainerStmt, err = DB.Prepare("INSERT INTO containers VALUES (?,?,?,?,?)")
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
}

/*
	Insert a container
*/
func (c *Container) Insert() error {
	var err error // Error handling

	_, err = InsertContainerStmt.Exec(c.ID, c.Hostname, c.Image, c.IPAddress, c.MacAddress)

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
	Get a container by $field
*/
func GetContainersBy(field string, value interface{}) []Container {
	var containers []Container

	rows, err := DB.Query("SELECT * FROM containers WHERE "+field+"=?", value)
	if err != nil {
		l.Critical(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tmpContainer Container
		var id string
		var hostname string
		var image string
		var ipAddress string
		var macAddress string

		if err := rows.Scan(&id, &hostname, &image, &ipAddress, &macAddress); err != nil {
			l.Critical(err)
		}

		tmpContainer.ID = id
		tmpContainer.Hostname = hostname
		tmpContainer.Image = image
		tmpContainer.IPAddress = ipAddress
		tmpContainer.MacAddress = macAddress

		containers = append(containers, tmpContainer)
	}
	if err := rows.Err(); err != nil {
		l.Critical(err)
	}

	return containers
}

/*
	Get a container by id
*/
func GetContainerById(id string) Container {
	var container Container
	var iD string
	var hostname string
	var image string
	var ipAddress string
	var macAddress string

	err := GetContainerByIdStmt.QueryRow(id+"b").Scan(&iD, &hostname, &image, &ipAddress, &macAddress)
	if err != nil {
		l.Critical(err)
	}

	container.ID = id
	container.Hostname = hostname
	container.Image = image
	container.IPAddress = ipAddress
	container.MacAddress = macAddress

	return container
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
