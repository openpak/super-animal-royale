# Ports

| Port | Use |
| --- | --- |
| 23620 | plain HTTP PlayFab tenant API (behind whatever terminates TLS for `<title>.playfabapi.com`) |
| 23621 | `/healthz` |
| 23622–23629 | reserved (Party/QoS listeners, once the traffic shape is known) |

23620–23629 is this title's pre-assigned block; do not change it.
