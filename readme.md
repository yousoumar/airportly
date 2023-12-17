# Airport Weather

To test, what's done so fare, first, make sure you have a MQTT server like Mosquitto running on your system, running on `localhost` and on port `1883`. Then, for example to test our Alert Manager, do the followings:

Open a terminal in the project root folder and run:
````sh
cd ./cmd/alert-manager && go run alert-manager.go
````
Open another terminal in the project root folder and run:
````sh
cd ./cmd/sensor/pressure && go run pressure.go
````

Go back to the first terminal, see the message.