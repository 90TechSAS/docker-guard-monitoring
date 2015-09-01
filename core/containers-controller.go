package core

import (
	"encoding/json"
	"errors"
	"sync"

	"../utils"
)

/*
	Container info
*/
type Container struct {
	CID        string
	Probe      string
	Hostname   string
	Image      string
	IPAddress  string
	MacAddress string
}

const (
	// File to store containerList
	ContainerListFilePath = "./containers.json"
)

var (
	// map[PROBE_NAME] => map[CONTAINER_ID] => Container
	containerList map[string]*map[string]*Container
	// containerList's Mutex
	ContainerListMutex sync.Mutex
)

/*
	Initialize containers controller
*/
func InitContainersController() {
	// Make map
	containerList = make(map[string]*map[string]*Container)

	// Check if ContainerListFilePath exists
	if utils.FileExists(ContainerListFilePath) {
		// Load containerList from the file
		err := LoadListFromFile()
		if err != nil {
			l.Critical("Can't load containers list from file:", err)
		}
	}

}

/*
	Load containerList from a file
*/
func LoadListFromFile() error {
	// Lock / Unlock containerList
	ContainerListMutex.Lock()
	defer ContainerListMutex.Unlock()

	// Read the file
	content, err := utils.FileReadAllBytes(ContainerListFilePath)
	if err != nil {
		return errors.New("LoadListFromFile: Failed to read list in file: " + err.Error())
	}

	// Parse the file
	err = json.Unmarshal(content, &containerList)
	if err != nil {
		return errors.New("LoadListFromFile: Failed to unmarshal struct: " + err.Error())
	}

	return nil
}

/*
	Load containerList to a file
*/
func SaveListToFile() error {
	// Lock / Unlock containerList
	ContainerListMutex.Lock()
	defer ContainerListMutex.Unlock()

	// containerList => json
	tmpJSON, err := json.Marshal(containerList)
	if err != nil {
		return errors.New("SaveListToFile: Failed to marshal struct: " + err.Error())
	}

	// Write json to file
	err = utils.FileWriteAllBytes(ContainerListFilePath, tmpJSON)
	if err != nil {
		return errors.New("SaveListToFile: Failed to write list in file: " + err.Error())
	}

	return nil
}

/*
	Insert a Container in containerList
*/
func (c *Container) Insert() error {
	// Lock / Unlock containerList
	ContainerListMutex.Lock()
	defer func() {
		ContainerListMutex.Unlock()
		SaveListToFile()
	}()

	// Check if probe exists
	probe, ok := containerList[c.Probe]

	// If probe doesn't exist, create the map of the probe
	if !ok {
		tmpProbe := make(map[string]*Container)
		containerList[c.Probe] = &tmpProbe
		probe = containerList[c.Probe]
	}

	// Insert container in the map
	(*probe)[c.CID] = c

	return nil
}

/*
	Delete a container in containerList
*/
func (c *Container) Delete() error {
	// Lock / Unlock containerList
	ContainerListMutex.Lock()
	defer func() {
		ContainerListMutex.Unlock()
		SaveListToFile()
	}()

	// Check if probe exists
	probe, ok := containerList[c.Probe]

	// If probe doesn't exist, return an error
	if !ok {
		return errors.New("Delete: In containerList, probe " + c.Probe + " didn't exists")
	}

	// Delete probe in the map
	delete(*probe, c.CID)

	return nil
}

/*
	Get containers by probe name in containerList
*/
func GetContainersByProbe(probeName string) ([]Container, error) {
	var containers []Container // Containers to return

	// Lock / Unlock containerList
	ContainerListMutex.Lock()
	defer ContainerListMutex.Unlock()

	// Get map of the probe
	probe, ok := containerList[probeName]

	// If probe doesn't exist, return an error
	if !ok {
		return containers, errors.New("Not found")
	}

	// Create temporary list of containers to return
	containers = make([]Container, len(*probe))

	// Insert containers in this list
	var i = 0
	for _, c := range *probe {
		containers[i] = *c
		i++
	}

	return containers, nil
}

/*
	Get a container by cid in containerList
*/
func GetContainerByCID(cid string) (Container, error) {
	var container Container // Container to return

	// Lock / Unlock containerList
	ContainerListMutex.Lock()
	defer ContainerListMutex.Unlock()

	// Search container
	for _, p := range containerList {
		for _, c := range *p {
			if c.CID == cid {
				return *c, nil
			}
		}
	}

	// If this code is reached, the container doesn't exist => not found
	return container, errors.New("Not found")
}
