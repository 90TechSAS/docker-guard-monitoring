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

___

#### GET /containers/probe/{name}

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

* [InfluxBD](https://github.com/influxdb/influxdb)
* [LoGo](https://github.com/Nurza/LoGo)

## License

MIT
