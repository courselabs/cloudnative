# Sock Shop on Docker Compose

Sock Shop is a microservices demo app. It runs as lots of small services which each own their own data, using the appropriate technology stack and communication pattern for that component.

Each part of the app is packaged and published as container images on Docker Hub, and you can model and run the whole app using Docker Compose.

## Reference

- [Sock Shop microservices demo app homepage](https://microservices-demo.github.io) 

- [Sock Shop deployment repo](https://github.com/microservices-demo/microservices-demo)

- [Sock Shop organization](https://github.com/microservices-demo) - contains source code repos for each component

- [Docker Compose Sock Shop documentation](https://microservices-demo.github.io/deployment/docker-compose.html)


## Run the app

The Compose file is in `sockshop/docker`.

There are lots of container images so it will be useful to pull them all before you run the app.

📋 Pull the images and start the app with Compose.

<details>
  <summary>Not sure how?</summary>

```
# pulls all the images in the model:
docker-compose -f ./sockshop/docker/docker-compose.yml pull

# starts all containers in detached mode:
docker-compose -f ./sockshop/docker/docker-compose.yml up -d
```

</details><br />

> When you run the project you'll see a warning about the MySQL password - you can ignore it.

📋 Check all the containers are running, find the entrypoint to the app and browse to it. and browse to the app. 

<details>
  <summary>Not sure how?</summary>

```
# this will show container status and published ports:
docker-compose -f ./sockshop/docker/docker-compose.yml ps
```

</details><br />

Your local Sock Shop deployment should look like this:

![The Sock Shop web application](/img/sock-shop.png)

You can log in with these credentials:

- username: `user`
- password: `password`

Then you can add an item to your cart and checkout. The full user experience should work correctly and you should see your order status changed to _shipped_.

## Explore the Compose file

Read through the Compose model in [docker-compose.yml](docker-compose.yml).

Can you answer these questions about the app?

- why do you see a warning about the password?
- what component is receiving traffic?
- why are two ports published?
- why doesn't the front end container receive traffic directly?

How does the model use these features of Compose?

- dependency mapping
- container networks
- OS capabilities 
- container restarts
- read only filesystems

Do those features mean the model is production-ready?

## Manage the app

_Check the application logs_

You can use the Compose command line to check container logs, and map out the workflow when an order is placed.

📋 Start by checking the logs from the front-end website.

<details>
  <summary>Not sure how?</summary>

```
docker-compose -f ./sockshop/docker/docker-compose.yml logs front-end
```

</details>

> The user and payment services also show useful logs.

_Scale the app_

Scale up the front-end service up to three containers. Check the app - are you still logged in? 

📋 Remove the first container and check again.

<details>
  <summary>Not sure how?</summary>

```
docker-compose -f ./sockshop/docker/docker-compose.yml up -d --scale front-end=3

docker rm -f docker_front-end_1
```

</details>

> The front end is a stateful app, it stores user sessions locally in each container.

## Clean up

📋 Remove the app containers with Compose.

<details>
  <summary>Not sure how?</summary>

```
docker-compose -f ./sockshop/docker/docker-compose.yml down
```

</details>