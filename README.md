# Docker Guard Monitoring 

## What?

Docker Guard is a powerful monitoring tool to watch your containers (running or not, memory/disk/netio usage, ...)

## Why ?

Because it's fast as hell! Docker Guard is a lightweight software and it can watch hundreds of containers (maybe thousands?).

## How it works?

![Architecture scheme](http://i.imgur.com/74qYu4z.png?1)

## How to install?

TODO

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


Feel free to fork the project a make a pull request!

## Thanks to

* [InfluxDB](https://github.com/influxdb/influxdb)
* [LoGo](https://github.com/Nurza/LoGo)

## License

MIT
