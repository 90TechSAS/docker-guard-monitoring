# Docker Guard Monitoring 

## What?

Docker Guard is a powerful monitoring tool to watch your containers (running or not, memory/disk/netio usage, ...)

## Why?

Because it's fast as hell! Docker Guard is a lightweight software and it can watch hundreds of containers (maybe thousands?).

## How does it work?

![Architecture scheme](http://i.imgur.com/74qYu4z.png?1)

## How to configure?

First, you need to copy the config example:

```bash
cp config.yaml.example config.yaml
```

Now you can edit the file ```config.yaml``` with your favorite editor before installing.

## How to install?

First, you need to install InfluxDB 0.9 or newer.
It's simple with Docker:

Make the InfluxDB data directory to make data persistent:
```bash
mkdir -p /var/lib/influxdb/data
```

And run the InfluxDB container from [tutumcloud/influxdb](https://github.com/tutumcloud/influxdb):
```bash
docker run -d -v "/var/lib/influxdb/data:/data" -p 8083:8083 -p 8086:8086 --expose 8090 --expose 8099 tutum/influxdb
```

Now, you can install Docker Guard Monitoring with docker:

Clone the project:
```bash
git clone https://github.com/90TechSAS/docker-guard-monitoring.git
```

Edit the file ```config.yaml``` at your own sweet will (see: "How to configure").
Type these commands to build a container with the Docker Guard Monitoring inside and run it!
Note that: when you are is the directory ```docker-guard-monitoring/docker``` and execute build.sh, this script will copy docker-guard-monitoring sources and config in the parent directory to the current directory.

```bash
cd docker-guard-monitoring/docker
./build.sh
./run.sh
```

Now, when you type ```docker ps``` you will see something like this:

```
CONTAINER ID        IMAGE                   COMMAND                  CREATED             STATUS              PORTS                                                                NAMES
865b1ccbb43b        dg-monitoring           "/bin/sh -c '/dgm/dg-"   3 minutes ago       Up 3 minutes        0.0.0.0:8124->8124/tcp                                               pensive_pare
9074d60353c2        tutum/influxdb:latest   "/run.sh"                13 days ago         Up 13 days          0.0.0.0:8083->8083/tcp, 8090/tcp, 0.0.0.0:8086->8086/tcp, 8099/tcp   jovial_bell
```

If you see the influxdb container + dg-monitoring container, it means that you did the job right.

## How to make my own transport?

First, what is a transport? A transport is an executable (script or binary) used for send an alert (like "OMG, this container is on fire!") on your favorite medium of communication (email, Slack, sms, webhook, ...).

But the best feature is: you can make your own transport!

To do this, you must create an executable in your transport directory (see: **How to configure?**).
This executable must have 5 parameters:

1. severity: The severity level (see the table bellow).
2. type: The alert type (see the table bellow).
3. target: The targeted system(s), it's generaly the concerned container's ID.
4. target_probe: The probe where the container is.
5. data: Additional data, it's detailed informations about the alert.

The transport will be executed like this example:

```./transports/mytransport.sh 1 CPUUsageOverload 169be77817 probe1 '8.45,6.12,2.89'```

**Severity levels:**

| Severity level | Description |
|--------------  |-------------|
| 0 			 | Notice 	   |
| 1 			 | Warning 	   |
| 2 			 | Critical	   |

**Alert types**

| Alert type              | Description                                               |
|-------------------------|-----------------------------------------------------------|
| DiskSpaceLimitReached   | The disk space limit of a container or probe is reached   |
| MemorySpaceLimitReached | The memory space limit of a container or probe is reached |
| ContainerStarted 		  | A container is started 									  |
| ContainerStopped 		  | A container is stopped 									  |
| ContainerCreated 		  | A container is created 									  |
| ContainerRemoved 		  | A container is removed 									  |
| NetBandwithOverload 	  | The net bandwith of a container overloaded                |
| CPUUsageOverload 		  | The cpu usage of a container or probe overloaded          |

**Example:**

This script is a transport example, it will log the alert in a file:

```bash
#!/bin/bash

LOG_FILE="test.log"

echo "Severity:     $1" >> $LOG_FILE
echo "Type:         $2" >> $LOG_FILE
echo "Target:       $3" >> $LOG_FILE
echo "Target_probe: $4" >> $LOG_FILE
echo "Data:         $5" >> $LOG_FILE
```

## API

#### GET /containers/{id}

**Description:**

Get one container's basic informations.
* $id : Container ID

**Example:**
```bash
curl -XGET  -u "dgadmin:password" "http://127.0.0.1:8124/containers/169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1"
```

**Result:**
```json
{
    "CID": "169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1",
    "Hostname": "169be7781716",
    "IPAddress": "172.17.0.1",
    "Image": "ubuntu",
    "MacAddress": "02:42:ac:11:00:01",
    "Probe": "probe1"
}
```

___

#### GET /containers/probe/{name}

**Description:**

Get one probe's list of containers.
* $name : Name of the probe

**Example:**
```bash
curl -XGET  -u "dgadmin:password" "http://127.0.0.1:8124/containers/probe/probe1"
```

**Result:**
```json
[
    {
        "CID": "169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1",
        "Hostname": "169be7781716",
        "IPAddress": "172.17.0.1",
        "Image": "ubuntu",
        "MacAddress": "02:42:ac:11:00:01",
        "Probe": "probe1"
    },
    {
        "CID": "33d62c50c2079d8b7d7cc18a235e7e7c24ef662ada953524f12047a3377de3c4",
        "Hostname": "33d62c50c207",
        "IPAddress": "172.17.0.2",
        "Image": "ubuntu",
        "MacAddress": "02:42:ac:11:00:02",
        "Probe": "probe1"
    }
]
```

___

#### GET /probes

**Description:**

Get the list of probes.

**Example:**
```bash
curl -XGET  -u "dgadmin:password" "http://127.0.0.1:8124/probes"
```

**Result:**
```json
[
    {
        "Containers": null,
        "DiskAvailable": 3530113024.0,
        "DiskTotal": 39277187072.0,
        "LoadAvg": "0.53,0.68,0.66",
        "MemoryAvailable": 15024922624.0,
        "MemoryTotal": 16807555072.0,
        "Name": "probe1",
        "Running": true
    },
    {
        "Containers": null,
        "DiskAvailable": 5468154694.0,
        "DiskTotal": 8644578945.0,
        "LoadAvg": "0.75,0.25,0.20",
        "MemoryAvailable": 12456842035.0,
        "MemoryTotal": 16807555072.0,
        "Name": "probe2",
        "Running": true
    }
]
```

___

#### GET /probes/{name}

**Description:**

Get one probe's basic informations + list of containers with basic informations.
* $name : Name of the probe

**Example:**
```bash
curl -XGET  -u "dgadmin:password" "http://127.0.0.1:8124/probes/probe1"
```

**Result:**
```json
{
    "Containers": [
        {
            "Hostname": "169be7781716",
            "IPAddress": "172.17.0.1",
            "Id": "169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1",
            "Image": "ubuntu",
            "MacAddress": "02:42:ac:11:00:01"
        },
        {
            "Hostname": "33d62c50c207",
            "IPAddress": "172.17.0.2",
            "Id": "33d62c50c2079d8b7d7cc18a235e7e7c24ef662ada953524f12047a3377de3c4",
            "Image": "ubuntu",
            "MacAddress": "02:42:ac:11:00:02"
        }
    ],
    "DiskAvailable": 3529875456.0,
    "DiskTotal": 39277187072.0,
    "LoadAvg": "0.45,0.62,0.66",
    "MemoryAvailable": 15016939520.0,
    "MemoryTotal": 16807555072.0,
    "Name": "probe1",
    "Running": true
}
```

___

#### GET /stats/probe/{name}

**Description:**

Get containers stats of a probe.
* $name : Name of the probe
 
GET parameters:

| Parameter     | Description                      | Example              | Default     |
|-------------- |----------------------------------|----------------------|-------------|
| since         | Date of the first stat (RFC3339) | 2015-09-02T09:27:41Z | now() - 24h |
| before        | Date of the last stat (RFC3339)  | 2015-09-02T09:27:41Z | now()       |
| limit         | Number of max stats returned     | 100                  | 20          |

**Example:**
```bash
curl -XGET  -u "dgadmin:password" "http://127.0.0.1:8124/stats/probe/probe1"
```

**Result:**
```json
[
    {
        "CPUUsage": 9,
        "ContainerID": "169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1",
        "NetBandwithRX": 460,
        "NetBandwithTX": 4518,
        "Running": true,
        "SizeMemory": 54644,
        "SizeRootFs": 386338816,
        "SizeRw": 386338816,
        "Time": "2015-09-02T09:27:38.142495446Z"
    },
    {
        "CPUUsage": 10,
        "ContainerID": "169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1",
        "NetBandwithRX": 56456,
        "NetBandwithTX": 566,
        "Running": true,
        "SizeMemory": 54678,
        "SizeRootFs": 386338816,
        "SizeRw": 386338816,
        "Time": "2015-09-02T09:27:39.231311889Z"
    },
    {
        "CPUUsage": 8,
        "ContainerID": "169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1",
        "NetBandwithRX": 456,
        "NetBandwithTX": 658,
        "Running": true,
        "SizeMemory": 54690,
        "SizeRootFs": 386338816,
        "SizeRw": 386338816,
        "Time": "2015-09-02T09:27:40.299626418Z"
    }
]
```

___

#### GET /stats/container/{id}

**Description:**

Get one container's stats.
* $id : Container ID
 
GET parameters:

| Parameter     | Description                      | Example              | Default     |
|-------------- |----------------------------------|----------------------|-------------|
| since         | Date of the first stat (RFC3339) | 2015-09-02T09:27:41Z | now() - 24h |
| before        | Date of the last stat (RFC3339)  | 2015-09-02T09:27:41Z | now()       |
| limit         | Number of max stats returned     | 100                  | 20          |

**Example:**
```bash
curl -XGET  -u "dgadmin:password" "http://127.0.0.1:8124/stats/container/169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1"
```

**Result:**
```json
[
    {
        "CPUUsage": 5,
        "ContainerID": "169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1",
        "NetBandwithRX": 456,
        "NetBandwithTX": 54,
        "Running": true,
        "SizeMemory": 78655,
        "SizeRootFs": 386338816,
        "SizeRw": 386338816,
        "Time": "2015-09-02T09:27:38.142495446Z"
    },
    {
        "CPUUsage": 6,
        "ContainerID": "169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1",
        "NetBandwithRX": 9789,
        "NetBandwithTX": 8965,
        "Running": true,
        "SizeMemory": 78461,
        "SizeRootFs": 386338816,
        "SizeRw": 386338816,
        "Time": "2015-09-02T09:27:39.231311889Z"
    },
    {
        "CPUUsage": 5,
        "ContainerID": "169be7781716d888835e0cafb46d7a0c3fc18a599406e45e6cf3816d345960d1",
        "NetBandwithRX": 6778,
        "NetBandwithTX": 78,
        "Running": true,
        "SizeMemory": 78765,
        "SizeRootFs": 386338816,
        "SizeRw": 386338816,
        "Time": "2015-09-02T09:27:40.299626418Z"
    }
]
```

___

## How to contribute?

Feel free to fork the project and make a pull request!

## Thanks to

* [InfluxDB](https://github.com/influxdb/influxdb)
* [LoGo](https://github.com/Nurza/LoGo)
* [Gorilla Mux](https://github.com/gorilla/mux)
* [Go-yaml](https://github.com/go-yaml/yaml)

## License

MIT
