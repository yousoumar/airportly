# Airport Weather

To test what's done so far, first we need a running Mosquitto on `localhost:1985`. For that, we use Docker. Open a terminal in the project root folder and run:

```sh
docker-compose up --build
```

Then, for example to test our Alert Manager, do the followings:

Open a new terminal in the project root folder and run:

```sh
cd ./cmd/alert-manager && go run alert-manager.go
```

Open another terminal in the project root folder and run:

```sh
cd ./cmd/sensor/pressure && go run pressure.go
```

Go back to the first terminal, see the message.
