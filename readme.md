# Airportly


As part of our Distributed Systems course at IMT Atlantique Engineering School, Airportly is about a system that collects, stores and exposes airports weather data coming from sensors.


## Overview


To run the project locally, you need to have `Docker` and `Docker Compose` running on your machine. Also, you need to have the `Go` programming language installed on your system. 

The project is composed of a MQTT server (Mosquitto), a NoSQL database (MongoDB), a REST API (Go), a file recorder (Go), a database recorder (Go), an alert manager (Go), and a UI (TypeScript & React).


Step one, open a terminal in the project root folder and run:


```sh
docker-compose up --build
```
This will start the MQTT server and the MongoDB database. Then, if your system can run Bash scripts, you can start all the remaining services with:

```sh
./start.sh
```


We simulate sensors with the `sensor` package. For example, to have a temperature sensor on Paris-Charles de Gaulle Airport (CDG being the [IATA code](https://en.wikipedia.org/wiki/IATA_airport_code) of the airport), we do:


```sh
./bin/sensor CDG temperature 23
```


`23` is the base metric (here temperature). We use it to generate random data based on it. We save the data sent by the sensors in CSV files in the project root folder (one per day and per airport), and into the MongoDB. Furthermore, we generate alerts when some thresholds are reached.


We expose the stored data through a REST API. The server will listen on `http://localhost:8080` (where you get an OpenAPI based documentation, known as Swagger). The UI  consumes the REST API, and runs on `http://localhost:4173`, something similar to:


<img src = "./ihm.png"/>

