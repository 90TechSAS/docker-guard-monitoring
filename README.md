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

Example:
```bash
curl -XGET  -u "dgadmin:password" "http://127.0.0.1:8124/containers/probe/probe1"
```

Result:
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

___

#### GET /probes/{name}

___

#### GET /stats/probe/{name}

___

#### GET /stats/container/{id}

___

## How to contribute?


Feel free to fork the project a make a pull request!

## Thanks to

* [InfluxDB](https://github.com/influxdb/influxdb)
* [LoGo](https://github.com/Nurza/LoGo)

## License

MIT
