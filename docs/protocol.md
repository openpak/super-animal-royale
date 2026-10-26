# Super Animal Royale (Switch) — wire facts

Title ID `0100DB3011E8C000`, Unity 2019.4.27f1 (IL2CPP), SDK 11.4.4.0.
Every fact below comes from the recon inputs; keep this file in step with
`cmd/super-animal-royaled`.

Sources:

- `scratch/decompile-super-animal-royale/RECON.md`
- `scratch/decompile-super-animal-royale/findings.json`

## Middleware: PlayFab only

- PlayFab C# SDK `UnitySDK-2.84.200402` (`X-PlayFabSDK` header), hosts
  `*.playfabapi.com` / `*.playfabsandbox.com` (suffixes only — the TitleId is
  unresolved; it lives in the compressed `PlayFabSharedSettings` asset inside
  `data.unity3d`).
- Login: `LoginWithNintendoSwitchDeviceId` and
  `LoginWithNintendoServiceAccount`/`LinkNintendoServiceAccount`, fed by
  `nn::account::EnsureNetworkServiceAccountAvailable` +
  `EnsureNetworkServiceAccountIdTokenCacheAsync`. The title also has its own
  `SuperAnimalAccount` email login and Modus publisher UI.
- Multiplayer Servers + matchmaking surface: `Create/Get/CancelMatchmakingTicket`,
  server-backfill tickets, `GetMultiplayerServerDetails`, `GetGameServerRegions`,
  `ListPartyQosServers`, `GetRemoteLoginEndpoint`.
- PlayFab Party native (statically linked in `main`): `PlayFabParty.dll`,
  `PartyXbl*` exports, `playfab.party[.internal]`, `com.playfab.events.default`;
  speech transcription via `TranscriberWebSocket*` + `tts.speech.microsoft.com`
  cognitive-services token. C# side `PartyCSharpSDK[.Interop]` +
  `PartyXBLCSharpSDK[.Interop]`; game layer `LobbyManagerPlayFabParty`
  (Create/Initialize/CleanUp) with network descriptor `superanimalparty`.

## Explicitly absent

- Photon realtime: only PlayFab SDK API names (`GetPhotonAuthenticationToken`,
  `PhotonApplicationId`) plus an XR `motionToPhoton` false positive;
  `ExitGames`/`LoadBalancing`/`PhotonPeer` = 0.
- Lidgren: dead PC-legacy reference (assembly ref + type names in metadata, zero
  string literals in `main`).
- UNet: unused engine code (native `NetworkTransport::*` bindings +
  `PlayerConnection` only; no HLAPI game types).
- NEX/PRUDP/`srv.nintendo.net`, EOS, Demonware, Steamworks, RakNet, gRPC,
  Mirror, MLAPI: all 0.
- No `PartyNintendo*` symbols. Xbox strings are InControl profiles + PlayFab
  Xbox API surface — no XSAPI login on Switch, so no Xbox gate is needed.

## Transport

- HTTPS via Unity (curl/7.64.1 + UnityTls) over the Switch `ssl` service, no
  bundled certs. Game-traffic sockets are runtime-allocated; no hardcoded game
  hosts or IPs (loopback only).
- NPDM network services: `acc:u0`, `bsd:u`, `bsdcfg`, `nsd:u`, `nifm:u`,
  `sfdnsres`, `ssl`, `ldn:u` (no PIA use — `nn::pia` = 0), `friend:u`, `bcat:u`,
  `pctl`.
- Non-game URLs: `discordapp.com` invite widget, `discord.gg`, `modusgames.com`,
  `go.microsoft.com`.

## Unresolved (from recon)

1. PlayFab TitleId + vertical/region defaults (needs Unity asset extraction).
2. Gameplay traffic shape: Party network UDP vs Multiplayer-Servers dedicated
   process, ports/protocols (needs IL2CPP dump + capture).
3. Lidgren confirm-dead via IL2CPP method refs.
4. `UNetUpdate` owner + `ldn:u` requester (minor).
5. No photon/PlayFab registry or SERVER entry for this title yet (needs TitleId).
