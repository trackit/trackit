# Trackit - Client (v2)

![trackit-logo](https://s3-us-west-2.amazonaws.com/trackit-public-artifacts/github-page/logo.png)

TrackIt helps you to optimize your AWS cloud

## Installation

### #1 Clone this repository

````sh
$> git clone https://github.com/trackit/trackit-client && cd trackit-client
````

### Via Docker

#### #2 Build Docker container

````sh
$> docker build -t trackit/ui .
````

#### #3 Start TrackIt

````sh
$> docker run --name Trackit_UI -e API_URL=http://localhost:8080trackit/ui
````

N.B. : If you are not running TrackIt on your local machine, you need to replace the URL of the API.

### Manually

#### #2 Update dependencies and build UI

````sh
$> yarn install
$> yarn run build
````

N.B. : You can also use `npm` instead of `yarn`

#### #3 Start TrackIt

You will need to set two environment variables :
- `API_URL` : URL of the API (Default value is `http://localhost:8080`)
- `UI_PORT` : Port of UI (Default value is `80`)

````sh
$> node production_server.js
````

N.B. : If you are not running TrackIt on your local machine, you need to replace the URL of the API.

# Screenshots

![account-namager](https://s3-us-west-2.amazonaws.com/trackit-public-artifacts/github-page/v2_account_wizard.png)

![cost-breakdown](https://s3-us-west-2.amazonaws.com/trackit-public-artifacts/github-page/v2_cost_breakdown_multi_charts.png)
