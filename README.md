# learn-pub-sub-starter (Peril)

This is the starter code used in Boot.dev's [Learn Pub/Sub](https://learn.boot.dev/learn-pub-sub) course.

## Message Brokers

### Connect

The starter code contains a few interesting things

- The `internal/gamelogic` package which contains the game logic for the Peril game.
- The `internal/routing` package which contains some routing constants for the game.
- Stubs of the `cmd/client` and `cmd/server` packages, which are `main` packages that run the client and server for the game.
- The `rabbit.sh` script: a convenience for starting and stopping **RabbitMQ** with Docker. You can run:
    - `./rabbit.sh start` to start **RabbitMQ**
    - `./rabbit.sh stop` to stop it
    - `./rabbit.sh logs` to view the server logs

### MQTT and STOMP

1. Create a file called Dockerfile 
    ```sh
    touch Dockerfile
    ```
2. Build the image, and name it `rabbitmq-stomp`
    ```sh
    docker build -t rabbitmq-stomp .
    ```
3. Run the container.
    ```sh
    docker run -d --rm --name rabbitmq -p 61613:5672 -p 15672:15672 rabbitmq-stomp

## Publishers & Queues

### Exchanges and Queues

Let's update our server to publish pause/resume messages to an **exchange** on a specific routing key.
The server can then communicate with all the various players of the game to let them know when the game is paused or resumed. 

### Decoupling

## Subscribers & Routing

## Delivery

## Serialization

## Scalability
