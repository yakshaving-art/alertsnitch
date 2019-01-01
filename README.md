# Alert Snitch

Captures Prometheus alertmanager alerts and writes them in a MySQL instance

## How to run

Run using docker in this very registry, for ex.

```sh
$ docker run --rm \
    -p 8080:8080 \
    -e MYSQL_DSN \
    registry.gitlab.com/yakshaving.art/alertsnitch
```

To run it requires a MySQL database to write to.

## Usage

Once AlertSnitch is up and running, configure the Prometheus Alert Manager to
forward every alert to it on the `/webhooks` path.

```yaml
To Be Provided configuration sample
```

## Readiness probe

AlertSnitch offers a `/-/ready` endpoint which will return 200 if the
application is ready to accept webhook posts.

During startup AlertSnitch will probe the MySQL database and the database
model. If everything works correctly it will set itself as ready.

In case of failure it will return a 500 and will write the error in the
response payload.

## Liveliness probe

AlertSnitch offers a `/-/health` endpoint which will return 200 as long as
the MySQL database is reachable.

In case of error it will return a 500 and will write the error in the
response payload.

## Metrics

AlertSnitch provides Prometheus metrics on `/metrics` as per prometheus
convention.

## Security

There is no offering of security of any kind. AlertSnitch is not ment to be
exposed to the internet but to be executed in an internal network reachable
by the alert manager.

## Grafana Compatibility

AlertSnitch writes alerts in such a way thay they can be explored using
Grafana's MySQL Data Source plugin. Refer to Grafana documentation for
further instructions.