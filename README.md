![TrackIt](https://s3-us-west-2.amazonaws.com/trackit-public-artifacts/github-page/logo.png)

[![Docker Pulls](https://img.shields.io/docker/pulls/msolution/trackit2-api.svg)](https://hub.docker.com/r/msolution/trackit2-api)
[![CircleCI](https://img.shields.io/circleci/build/github/trackit/trackit-server.svg)](https://circleci.com/gh/trackit/trackit-server)
[![GitHub](https://img.shields.io/github/license/trackit/trackit-server.svg)](LICENSE)

TrackIt is a tool to optimize your AWS cloud usage and spending.

## Features

- Easy account setup

![account-setup](https://s3.us-west-2.amazonaws.com/trackit-public-artifacts/github-page/v2_account_wizard.png)

- AWS Cost Breakdown

![cost-breakdown](https://s3-us-west-2.amazonaws.com/trackit-public-artifacts/github-page/v2_cost_breakdown_multi_charts.png)

- AWS Tags overview

![tags](https://s3-us-west-2.amazonaws.com/trackit-public-artifacts/github-page/v2_tags.png)

- Events alerts

![events](https://s3-us-west-2.amazonaws.com/trackit-public-artifacts/github-page/v2_events.png)

## How to use

### With Docker Compose

You can start using TrackIt by using the `docker-compose.yml` template available in this repository. It will pull Docker images from Docker Registry.

````sh
$> docker-compose up -d
````

You can also build locally the needed Docker images by using the `docker-compose.yml` file available in `docker/` folder.

````sh
$> docker-compose up -d -f docker/docker-compose.yml
````

### Manually

#### 0. Be sure all requirements below are met

- [Docker](https://docs.docker.com/engine/installation/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html) and [configure your credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html)
- [Golang](https://golang.org/doc/install)
- [Govender](https://github.com/kardianos/govendor)

#### 1. Clone this repository

````sh
$> mkdir -p $HOME/go/src/github.com/trackit
$> cd $HOME/go/src/github.com/trackit
$> git clone https://github.com/trackit/trackit
$> cd trackit-server
````

#### 2. Check out dependencies

````sh
$> govendor sync -v
````

#### 3. Start TrackIt

````sh
$> ./start.sh
````

Note: On most operating systems, you will need to [increase the mmap limit](https://www.elastic.co/guide/en/elasticsearch/reference/current/vm-max-map-count.html) to allow elasticsearch to run properly:

````sh
$> sudo sysctl -w vm.max_map_count=262144
````

#### 4. Now you can use TrackIt

TrackIt API is now listening on `localhost:8580`

## Web UI

A Web UI made with React is available here: [TrackIt Client](https://github.com/trackit/trackit2-client)

## API documentation

The API exposes its own documentation on the `GET /docs` route, in JSON format.
Also, the documentation for each route can be retrieved by an `OPTIONS`
request. We are working on an actual viewer for this.

## Recommendation plugins

Trackit uses a plugin system to easily implement new recommendation checks.
Informations on how to write plugins are available in a README in the `plugins` directory.
