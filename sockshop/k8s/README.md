
## Reference

- [Kubernetes Sock Shop documentation ](https://microservices-demo.github.io/deployment/kubernetes-start.html)


## Run the app

The Kubernetes manifests are in `sockshop/k8s`.

There are lots of container images so you may want to pull them all before you run the app.

```
k apply -f .\sockshop\k8s\
```

Check the Pods and Services and browse to the app. It should look like this:

![The Sock Shop web application](/img/sock-shop.png)

You can log in with these credentials:

- username: `user`
- password: `password`

Then you can add an item to your cart and checkout.


## Explore the manifests

Model:

- check the front end Deployment - what might the session configuration be used for?
- can the databases be scaled?
- why are two ports published?

Kubernetes features:

- resource blocks
- security contexts
- probes
- nodeSelector


## Managing the app

_Check the application logs_.

List the Pods in the namespace and you'll see they all have component name labels. You can use those to print logs.

- user, payment show useful logs

Scale the front-end service up to two replicas. Carry on using the app. What's the user experience?

k scale deploy/front-end --replicas=3 -n sock-shop


## Clean up

Remove the app.

k delete ns sock-shop