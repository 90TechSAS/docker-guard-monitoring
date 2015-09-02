package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	influxdb "github.com/influxdb/influxdb/client"

	"../utils"
)

/*
	Container's stats
*/
type Stat struct {
	ContainerID   string
	Time          time.Time
	SizeRootFs    uint64
	SizeRw        uint64
	SizeMemory    uint64
	NetBandwithRX uint64
	NetBandwithTX uint64
	CPUUsage      uint64
	Running       bool
}

/*
	HTTP GET options
*/
type Options struct {
	Since  int
	Before int
	Limit  int
}

const (
	StatsMeasurements = "cstats"
)

/*
	Client variables
*/
var (
	// DB
	DB *influxdb.Client
)

/*
	Initialize InfluxDB connection
*/
func InitDB() {
	var err error

	// Parse InfluxDB server URL
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", DGConfig.DockerGuard.InfluxDB.IP, DGConfig.DockerGuard.InfluxDB.Port))
	if err != nil {
		l.Critical("Can't parse InfluxDB config :", err)
	}

	// Make InfluxDB config
	conf := influxdb.Config{
		URL:      *u,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	}

	// Connect to InfluxDB server
	DB, err = influxdb.NewClient(conf)
	if err != nil {
		l.Critical("Can't connect to InfluxDB:", err)
	}

	// Test InfluxDB server connectivity
	dur, ver, err := DB.Ping()
	if err != nil {
		l.Critical("Can't ping InfluxDB:", err)
	}
	l.Verbose("Connected to InfluxDB! ping:", dur, "/ version:", ver)

	// Create DB if doesn't exist
	_, err = queryDB(DB, "create database "+DGConfig.DockerGuard.InfluxDB.DB)
	if err != nil {
		if err.Error() != "database already exists" {
			l.Critical("Create DB:", err)
		}
	}
}

/*
	Send a query to InfluxDB server
*/
func queryDB(con *influxdb.Client, cmd string) (res []influxdb.Result, err error) {
	q := influxdb.Query{
		Command:  cmd,
		Database: DGConfig.DockerGuard.InfluxDB.DB,
	}
	if response, err := con.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
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
	Insert a stat
*/
func (s *Stat) Insert() error {
	var pts = make([]influxdb.Point, 1) // InfluxDB point
	var err error                       // Error handling

	l.Silly("Insert stat:", s)
	// Make InfluxDB point
	pts[0] = influxdb.Point{
		Measurement: StatsMeasurements,
		Tags: map[string]string{
			"containerid": s.ContainerID,
		},
		Fields: map[string]interface{}{
			"sizerootfs":    s.SizeRootFs,
			"sizerw":        s.SizeRw,
			"sizememory":    s.SizeMemory,
			"netbandwithrx": s.NetBandwithRX,
			"netbandwithtx": s.NetBandwithTX,
			"cpuusage":      s.CPUUsage,
			"running":       s.Running,
		},
		Time:      time.Now(),
		Precision: "s",
	}

	// InfluxDB batch points
	bps := influxdb.BatchPoints{
		Points:          pts,
		Database:        DGConfig.DockerGuard.InfluxDB.DB,
		RetentionPolicy: "default",
	}

	// Write point in InfluxDB server
	timer := time.Now()
	_, err = DB.Write(bps)
	if err != nil {
		l.Error("Failed to write in InfluxDB:", bps, ". Error:", err)
	} else {
		l.Silly("Stat inserted in ", time.Since(timer), ":", bps)
	}

	return err
}

/*
	Insert some stats
*/
func InsertStats(stats []Stat) error {
	if len(stats) < 1 {
		return errors.New("len(stats) < 1")
	}

	var pts = make([]influxdb.Point, len(stats)) // InfluxDB point
	var err error                                // Error handling

	l.Silly("Insert stats:", stats)
	// Make InfluxDB points
	for i := 0; i < len(stats); i++ {
		pts[i] = influxdb.Point{
			Measurement: StatsMeasurements,
			Tags: map[string]string{
				"containerid": stats[i].ContainerID,
			},
			Fields: map[string]interface{}{
				"sizerootfs":    stats[i].SizeRootFs,
				"sizerw":        stats[i].SizeRw,
				"sizememory":    stats[i].SizeMemory,
				"netbandwithrx": stats[i].NetBandwithRX,
				"netbandwithtx": stats[i].NetBandwithTX,
				"cpuusage":      stats[i].CPUUsage,
				"running":       stats[i].Running,
			},
			Time:      time.Now(),
			Precision: "s",
		}
	}

	// InfluxDB batch points
	bps := influxdb.BatchPoints{
		Points:          pts,
		Database:        DGConfig.DockerGuard.InfluxDB.DB,
		RetentionPolicy: "default",
	}

	// Write points in InfluxDB server
	timer := time.Now()
	_, err = DB.Write(bps)
	if err != nil {
		l.Error("Failed to write in InfluxDB:", bps, ". Error:", err)
	} else {
		l.Silly("Stat inserted in ", time.Since(timer), ":", bps)
	}

	return err
}

/*
	Get container's last stat
*/
func (c *Container) GetLastStat() (Stat, error) {
	var stat Stat // Returned stat
	var err error // Error handling

	query := `	SELECT 	last(cpuusage),
						last(netbandwithrx),
						last(netbandwithtx),
						last(running),
						last(sizememory),
						last(sizerootfs),
						last(sizerw) 
				FROM cstats
				WHERE containerid = '` + c.CID + `'`

	// Send query
	res, err := queryDB(DB, query)
	if err != nil {
		return stat, err
	}

	// Get results
	for _, row := range res[0].Series[0].Values {
		var statValues [8]int64
		if len(row) != 8 {
			return stat, errors.New(fmt.Sprintf("GetLastStat: Wrong stat length: %d != 8", len(row)))
		}
		for i := 1; i <= 7; i++ {
			if i == 4 {
				continue
			}
			statValues[i], err = row[i].(json.Number).Int64()
			if err != nil {
				return stat, errors.New("GetLastStat: Can't parse value: " + row[i].(string))
			}
		}

		stat.ContainerID = c.CID
		stat.CPUUsage = uint64(statValues[1])
		stat.NetBandwithRX = uint64(statValues[2])
		stat.NetBandwithTX = uint64(statValues[3])
		stat.Running = row[4].(bool)
		stat.SizeMemory = uint64(statValues[5])
		stat.SizeRootFs = uint64(statValues[6])
		stat.SizeRw = uint64(statValues[7])
	}

	return stat, err
}

/*
	Get stats by container id
*/
func GetStatsByContainerCID(containerCID string, o Options) ([]Stat, error) {
	var stats []Stat  // List of stats to return
	var oS, oB string // Query options
	var err error     // Error handling

	query := `	SELECT *
				FROM cstats
				WHERE containerid = '` + containerCID + `'`

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
		query += fmt.Sprintf(" AND time > '%s' AND time < '%s'", oS, oB)
	}
	if o.Limit != -1 {
		query += fmt.Sprintf(" LIMIT %d", o.Limit)
	}

	// Send query
	res, err := queryDB(DB, query)
	if err != nil {
		return nil, err
	}

	// Get results
	for _, row := range res[0].Series[0].Values {
		var stat Stat
		var statValues [8]int64
		if len(row) != 8 {
			return nil, errors.New(fmt.Sprintf("GetLastStat: Wrong stat length: %d != 8", len(row)))
		}
		for i := 1; i <= 7; i++ {
			if i == 4 {
				continue
			}
			statValues[i], err = row[i].(json.Number).Int64()
			if err != nil {
				return nil, errors.New("GetLastStat: Can't parse value: " + row[i].(string))
			}
		}

		stat.Time, _ = time.Parse(time.RFC3339, row[0].(string))
		stat.ContainerID = containerCID
		stat.CPUUsage = uint64(statValues[1])
		stat.NetBandwithRX = uint64(statValues[2])
		stat.NetBandwithTX = uint64(statValues[3])
		stat.Running = row[4].(bool)
		stat.SizeMemory = uint64(statValues[5])
		stat.SizeRootFs = uint64(statValues[6])
		stat.SizeRw = uint64(statValues[7])

		stats = append(stats, stat)
	}
	return stats, nil
}

/*
	Get stats by probe name
*/
func GetStatsByContainerProbeID(probeName string, o Options) ([]Stat, error) {
	var containers []Container // List of containers in the probe
	var stats []Stat           // List of stats to return
	var err error              // Error handling

	// Get list of containers in the probe
	containers, err = GetContainersByProbe(probeName)
	if err != nil {
		return nil, err
	}

	// Get stats for each containers
	for _, container := range containers {
		tmpStats, err := GetStatsByContainerCID(container.CID, o)
		if err != nil {
			return nil, err
		}
		for _, tmpStat := range tmpStats {
			stats = append(stats, tmpStat)
		}
	}

	return stats, nil
}

/*
	Get stats populated by probe name
*/
func GetStatsPByContainerProbeID(probeName string, o Options) ([]StatPopulated, error) {
	var containers []Container // List of containers in the probe
	var statsP []StatPopulated // List of stats populated to return
	var err error              // Error handling

	// Get list of containers in the probe
	containers, err = GetContainersByProbe(probeName)
	if err != nil {
		return nil, err
	}

	// Get stats for each containers
	for _, container := range containers {
		tmpStats, err := GetStatsByContainerCID(container.CID, o)
		if err != nil {
			return nil, err
		}
		for _, tmpStat := range tmpStats {
			statP := StatPopulated{
				Container:     container,
				Time:          tmpStat.Time,
				SizeRootFs:    tmpStat.SizeRootFs,
				SizeRw:        tmpStat.SizeRw,
				SizeMemory:    tmpStat.SizeMemory,
				NetBandwithRX: tmpStat.NetBandwithRX,
				NetBandwithTX: tmpStat.NetBandwithTX,
				CPUUsage:      tmpStat.CPUUsage,
				Running:       tmpStat.Running,
			}

			statsP = append(statsP, statP)
		}
	}

	return statsP, nil
}
