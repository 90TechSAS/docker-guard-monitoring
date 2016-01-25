package core

import (
	"encoding/json"
	"errors"
	"sync"

	dguard "github.com/90TechSAS/libgo-docker-guard"

	"../utils"
)

const (
	// File to store containerList
	ContainerListFilePath = "./containers.json"
)

var (
	// map[PROBE_NAME] => map[CONTAINER_ID] => Container
	containerList map[string]*map[string]*dguard.Container
	// containerList's Mutex
	ContainerListMutex sync.Mutex
)

/*
	Initialize containers controller
*/
func InitContainersController() {
	// Make map
	containerList = make(map[string]*map[string]*dguard.Container)

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
func InsertContainer(c *dguard.Container) error {
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
		tmpProbe := make(map[string]*dguard.Container)
		containerList[c.Probe] = &tmpProbe
		probe = containerList[c.Probe]
	}

	// Insert container in the map
	(*probe)[c.ID] = c

	return nil
}

/*
	Delete a container in containerList
*/
func DeleteContainer(c *dguard.Container) error {
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
	delete(*probe, c.ID)

	return nil
}

/*
	Get containers by probe name in containerList
*/
func GetContainersByProbe(probeName string) ([]dguard.Container, error) {
	var containers []dguard.Container // Containers to return

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
	containers = make([]dguard.Container, len(*probe))

	// Insert containers in this list
	var i = 0
	for _, c := range *probe {
		containers[i] = *c
		i++
	}

	return containers, nil
}

/*
	Get []dguard.SimpleContainer by probe name in containerList
*/
func GetSimpleContainersByProbe(probeName string) ([]dguard.SimpleContainer, error) {
	var tmpContainers []dguard.Container
	var simpleContainers []dguard.SimpleContainer
	var err error // Error handling

	tmpContainers, err = GetContainersByProbe(probeName)
	if err != nil {
		return nil, err
	}

	for _, c := range tmpContainers {
		var simpleContainer = dguard.SimpleContainer{
			ID:         c.ID,
			Hostname:   c.Hostname,
			Image:      c.Image,
			IPAddress:  c.IPAddress,
			MacAddress: c.MacAddress,
		}
		simpleContainers = append(simpleContainers, simpleContainer)
	}

	return simpleContainers, nil
}

/*
	Get a container by cid in containerList
*/
func GetContainerByCID(cid string) (dguard.Container, error) {
	var container dguard.Container // Container to return

	// Lock / Unlock containerList
	ContainerListMutex.Lock()
	defer ContainerListMutex.Unlock()

	// Search container
	for _, p := range containerList {
		for _, c := range *p {
			if c.ID == cid {
				return *c, nil
			}
		}
	}

	// If this code is reached, the container doesn't exist => not found
	return container, errors.New("Not found")
}

/*
	Get list of probes infos
*/
func GetProbesInfos() []dguard.ProbeInfos {
	var probes []dguard.ProbeInfos // List of probes infos to return

	for _, probe := range Probes {
		if probe == nil {
			l.Error("GetProbesInfos: probe can't be nil")
			continue
		}
		probes = append(probes, *((*probe).Infos))
	}

	return probes
}
