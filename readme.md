# Airport Weather

As part of our Distributed Systems course at IMT Atlantique Engineering School, this project is about a system that collects, stores and exposes airports weather data coming from sensors.

## Running the project locally

## Requirements

You need to have `Docker` and `Docker Compose` running on your machine. Also, you need to have `Golang` installed on your system.

## Setting up the project

### Spin up MongoDB and Mosquitto

We use MongoDB as database, and Mosquitto as MQQT server. We run and connect them to other parts of the project using `Docker Compose`.

Open a terminal in the project root folder and run:

```sh
docker-compose up --build
```

### Run some sensors to send some data to Mosquitto

We simulate sensors with the `sensor.go` file. To run it, open a new terminal in the project root folder and run:

```sh
cd ./cmd/sensor && go run sensor.go CDG temperature 23
```

The above command is like saying: I need a temperature sensor on Paris-Charles de Gaulle Airport (CDG being the [IATA code](https://en.wikipedia.org/wiki/IATA_airport_code) of the airport). `23` is the base metric (here temperature). We use it to generate random data based on it.

To have a sensor on the same airport measuring pressure, open another terminal in the project root folder and run:

```sh
cd ./cmd/sensor && go run sensor.go CDG pressure 106
```

### Let's record the data sent by the sensors into CVS files

We save the data sent by the sensors in CSV files in the same folder as the corresponding script (one per day and per airport). For that, open another terminal in the project root folder and run:

```sh
cd ./cmd/file-recorder && go run file-recorder.go
```

### Let's record the data sent by the sensors into MongoDB

For that, we use the database lunched with `Docker` above. Open another terminal in the project root folder and run:

```sh
cd ./cmd/database-recorder && go run database-recorder.go
```

### Let's access the data in MongoDB with a REST API

We expose the stored data through a REST API running on `http://localhost:8080`. Open another terminal in the project root folder and run:

```sh
cd ./cmd/http-rest-server && go run http-rest-server.go
```

There is an OpenAPI based documentation of the API in `cmd/http-rest-server/OpenAPI.yaml`.

### Getting alerts when metrics reach some threshold

We send alerts to clients who are subscribed on the topic `airport/codeIata/alert/#` for all metrics related to an airport, or for a specific metric like the temperature with `airport/codeIata/alert/temperature`, when some threshold are reached.

To test it, open another terminal in the project root folder and run:

```sh
cd ./cmd/alert-manager && go run alert-manager.go
```

### Let's visualize the data in a Graphical Interface in a browser

We built a client app with `React`, `MUI` and `vite`, that consumes the REST API. To test it out, open another terminal in the project root folder and run (you need `node` and `npm` installed on your computer):

```sh
cd ./ihm && npm i && npm run dev
```

You would get the app on `http://localhost:5173/`, something similar to:

<img src = "./ihm.png"/>
