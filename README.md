# Super Animal Royale

OpenPak's Super Animal Royale server.

| | |
| --- | --- |
| Game | Super Animal Royale |
| Network | Nintendo Switch |
| Title ID | `0100DB3011E8C000` |
| Middleware | PlayFab Client + Multiplayer Servers + Matchmaking/QoS + Party (no Photon realtime, no NEX) |
| Licence | AGPL-3.0-only |

## Status

Scaffold. Every PlayFab tenant route answers `501 Not Implemented` until the
title's PlayFab TitleId is recovered and `PLAYFAB_TITLE_ID` is set. Voice/chat
relay defers to the shared [`playfab`](https://github.com/openpak/playfab) service.

## What it serves

| Host | What | Served by |
| --- | --- | --- |
| `<title>.playfabapi.com` | Client login, multiplayer servers, matchmaking, QoS | this repo (`super-animal-royaled`), once the TitleId is known |
| `playfab.party` | Voice/chat/networks (Party native) | [`playfab`](https://github.com/openpak/playfab) (`internal/party`) |

## Run

```sh
go run ./cmd/super-animal-royaled
```

| Variable | Meaning |
| --- | --- |
| `SAR_HTTP_ADDR` | PlayFab tenant listener (default `:23620`) |
| `SAR_HEALTH_ADDR` | `GET /healthz` (default `:23621`) |
| `PLAYFAB_TITLE_ID` | PlayFab TitleId for this title; unset (tenant routes answer 501) |

## Docs

- [Protocol](docs/protocol.md)
- [Ports](docs/ports.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
