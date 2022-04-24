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