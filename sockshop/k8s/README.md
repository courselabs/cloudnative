# Sock Shop on Kubernetes

Sock Shop is a microservices demo app. It runs as lots of small services which each own their own data, using the appropriate technology stack and communication pattern for that component.

Each part of the app is packaged and published as container images on Docker Hub, and you can model and run the whole app using Kubernetes.

## Reference

- [Sock Shop microservices demo app homepage](https://microservices-demo.github.io) 

- [Sock Shop deployment repo](https://github.com/microservices-demo/microservices-demo)

- [Sock Shop organization](https://github.com/microservices-demo) - contains source code repos for each component

- [Kubernetes Sock Shop documentation ](https://microservices-demo.github.io/deployment/kubernetes-start.html)

## Run the app

The Kubernetes manifests are in `sockshop/k8s`.

📋 Deploy the application.

<details>
  <summary>Not sure how?</summary>

```
kubectl apply -f ./sockshop/k8s/
```

</details><br />

📋 Check the Pods and Services and browse to the app. 

<details>
  <summary>Not sure how?</summary>

The model deploys multiple Pods and Services - they all have the label `app=sockshop`:

```
kubectl get po -l app=sockshop

kubectl get svc -l app=sockshop
```

_There is no dependency mapping in the app specs, but there are health checks. It will take at least 1 minutes before all the Pods are in the ready state_.

When the Pods are all ready you'll be able to browse to http://localhost:30001 to reach the front-end NodePort Service.

</details><br />

Your local Sock Shop deployment should look like this:

![The Sock Shop web application](/img/sock-shop.png)

You can log in with these credentials:

- username: `user`
- password: `password`

Then you can add an item to your cart and checkout. The full user experience should work correctly and you should see your order status changed to _shipped_.

## Explore the Kubernetes manifests

There is a lot of YAML for this app :)

Can you answer these questions about the app from the specs?

- check the front end Deployment in [09-front-end-dep.yaml](09-front-end-dep.yaml) - what might the `SESSION` configuration be used for?
- why is the user service in [25-user-dep.yaml](25-user-dep.yaml) unavailable for at least one minute after starting up?
- why are there two containers defined for the RabbitMQ message queue Pod in [19-rabbitmq-dep.yaml](19-rabbitmq-dep.yaml)

How does the model use these features of Kubernetes?

- resource limits and requests
- security contexts
- container probes
- node selectors
- configuration
- read only filesystems
- pod controllers

Do those features mean the model is production-ready?

## Managing the app

_Check the application logs_.

📋 List the Pods in the namespace and you'll see they all have component name labels. You can use those to print logs for different services.

<details>
  <summary>Not sure how?</summary>

```
kubectl get po --show-labels

kubectl logs -l name=user 

kubectl logs -l name=payment

kubectl logs -l name=catalogue 
```


</details>


📋 Scale the front-end service up to three replicas. Carry on using the app - are you still logged in? Delete all the Pods; when they get recreated are you still logged in to your session?

<details>
  <summary>Not sure how?</summary>

```
# scale up - the front end component is stateless:
k scale deploy/front-end --replicas=3 

# remove all the Pods - the new ones can still load your web session:
kubectl delete po -l name=front-end 
```

</details>

> The front end is configured to be stateless, it stores session state in the Redis Pod which can be read by all front end Pods.

## Clean up

📋 Remove the whole app with a single Kubectl command.

<details>
  <summary>Not sure how?</summary>

```
kubectl delete -f ./sockshop/k8s/
```

</details>