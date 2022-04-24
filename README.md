# hue-exporter

hue-exporter is a Prometheus exporter, which will read all light information from a Hue bridge, and publish it on the
`/metrics` endpoint for Prometheus to scrape.

The exporter only calls Hue bridge when Prometheus calls  the `/metrics` endpoint.

The exporter has to run **locally**, on the same LAN as the Hue bridge.

## Configuration

Configuration is done using the following environment variables:

| Name          | Description                                                 |
|---------------|-------------------------------------------------------------|
| HUE_CLIENT_ID | The Hue client ID. See how to get one in the section below. |
| HUE_BRIDGE_IP | The IP for the Hue Bridge.                                  |

## Get Hue Client ID

In order to get a Hue Client ID, you should send a `POST` request to the bridge's `/api` endpoint.
The body should consist of JSON object with a `devicetype` property.

```json
{
  "devicetype": "hue-exporter"
}
```

CURL:
```shell
curl -k -X POST -d '{"devicetype": "hue-exporter"}' -H "Content-Type: application/json" https://<bridgeIP>/api
```

Before sending this request, you need to press the link button (the big round button on the bridge), otherwise the Hue
Bridge will answer
```json
[{"error":{"type":101,"address":"","description":"link button not pressed"}}]
```