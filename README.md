# Ecommerce Notification Service

This is the notification service that handles all notifications in the ecommerce application. Requests are made to the notification service through NATS messaging

## Applications

* ## [Jaeger](https://www.jaegertracing.io/)

  * Jaeger is an **open source software for tracing transactions between distributed services**.
* ## [Nats](https://nats.io)

  * NATS is **an open-source messaging system** (sometimes called message-oriented middleware).
  * All requests to the notification service is made through NATS, this is useful for background processing and retries.

### Usage

To install / run the user microservice run the command below:

```bash
docker-compose up
```

## Requirements

The application requires the following:

* Go (v 1.5+)
* Docker (v3+)
* Docker Compose

## Other Micro-Services / Resources

* #### [Product Service](https://github.com/wisdommatt/ecommerce-microservice-product-service)
* #### [User Service](https://github.com/wisdommatt/ecommerce-microservice-user-service)
* #### [Cart Service](https://github.com/wisdommatt/ecommerce-microservice-cart-service)
* #### [Shared](https://github.com/wisdommatt/ecommerce-microservice-shared)

## Public API

The public graphql API that interacts with the microservices internally can be found in [https://github.com/wisdommatt/ecommerce-microservice-public-api](https://github.com/wisdommatt/ecommerce-microservice-public-api).
