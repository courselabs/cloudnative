
## Reference

- [Docker Compose Sock Shop documentation ](https://microservices-demo.github.io/deployment/docker-compose.html)


## Run the app

The Compose file is in `sockshop/docker`.

There are lots of container images so you may want to pull them all before you run the app.

```
docker-compose -f .\sockshop\docker\docker-compose.yml pull

docker-compose -f .\sockshop\docker\docker-compose.yml up -d
```

When you run the project you'll see a warning about the MySQL password; you can ignore that.

Check all the containers are running and browse to the app. It should look like this:

![The Sock Shop web application](/img/sock-shop.png)

You can log in with these credentials:

- username: `user`
- password: `password`

Then you can add an item to your cart and checkout.

## Explore the Compose file

Model:

- why do you see a warning about the password?
- what component is receiving traffic?
- why are two ports published?
- why doesn't the front end container receive traffic directly?

Compose features:

- no dependency mapping
- capabilities dropped & added
- restart flag
- read_only setting

## Manage the app

_Check the application logs_.

You can use the Compose command line to check container logs, and map out the workflow when an order is placed.

docker-compose -f .\sockshop\docker\docker-compose.yml logs front-end

- user, payment show useful logs

Scale up the front-end service up to three containers. Check the app - are you still logged in? Remove the first container and check again.

docker-compose -f .\sockshop\docker\docker-compose.yml up -d --scale front-end=3

docker rm -f docker_front-end_1

## Clean up

Close down the app with Compose.


docker-compose -f .\sockshop\docker\docker-compose.yml down