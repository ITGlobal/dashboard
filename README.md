dashboard
=========

Dashboard app for Raspberry Pi

How to configure, build and run
-------------------------------

1. Clone this repository somewhere (e.g. `~/dashboard`)
2. Create configuration file `dashboard.json` in `config` directory. See below to find out about config file's format

   *Note that this config file will be ignored by git*
3. Make sure you have docker and docker-compose installed.
4. Type a docker-compose command to build and run **dashboard**:
```shell
# for x64 machines
docker-compose up -d

#  - or -
# for ARM machines (like Raspberry Pi)

docker-compose -f docker-compose.arm.yml up -d

#  - or -
# to run manually on linux

chmod +x ./build.sh
./run.sh

#  - or -
# to run manually on windows

./run.cmd
```
5. Open `http://localhost:8000` to see dashboard page.

Configuration file format
-------------------------

Configuration file has the following format:

```json
{
  "providers" : [
    // list of provider configurations
  ]
}
```

Each item in `providers` array represents one tile provider, while each tile provider may generate one or few tiles.

Common provider config has the following fields:

```json
{
  "type" : "PROVIDER_TYPE",
  "enabled" : true
}
```

Parameter `type` is required and should contain a valid provider type (see below), while parameter `enabled` is optional and defaults to `true`.

Dashboard tile providers
------------------------


## `sim` provider

Simulated tile provider (for debugging purposes). Generated and updates `N` tiles.

| Parameter | Type    | Is optional         | Description                 |
|-----------|---------|---------------------|-----------------------------|
| `count`   | integer | Yes, defaults to 10 | Amount of tiles to generate |


## `check` provider

Periodially checks specified HTTP/HTTPS URL and displays its status.

If the specified URL will respond to HTTP GET with any 2xx or 3xx HTTP status then the dashboard tile will be marked as `SUCCESS`.
Otherwise, it will be marked as `ERROR`.

| Parameter | Type   | Is optional                | Description        |
|-----------|--------|----------------------------|--------------------|
| `url`     | string | Required                   | URL to check       |
| `name`    | string | Yes, defaults to host name | Tile name          |
| `timer`   | string | Yes, defaults to `1m`      | URL check interval |

## `mongodb` provider

A provider that checks MongoDB replica status.

If all replica set nodes are healthy then the dashboard tile will be marked as `SUCCESS`.
Otherwise, it will be marked as `ERROR`. 

| Parameter | Type   | Is optional            | Description                        |
|-----------|--------|------------------------|------------------------------------|
| `name`    | string | Required               | Tile name                          |
| `url`     | string | Required               | MongoDB replica set connection URL |
| `timer`   | string | Yes, defaults to `20s` | MongoDB check interval             |

## `teamcity` provider

A provider that displays statuses of TeamCity builds.

Provider will populate one dashboard item for each TeamCity project which will indicate current build status.

| Parameter  | Type   | Is optional            | Description                   |
|------------|--------|------------------------|-------------------------------|
| `url`      | string | Required               | TeamCity URL                  |
| `username` | string | Optional               | Username for TeamCity access  |
| `password` | string | Optional               | Password for TeamCity access  |
| `timer`    | string | Yes, defaults to `20s` | TeamCity build fetch interval |


## `1cloud` provider

A provider that checks 1Cloud payment status and displays a warning if remaining account balance drops below limit.

Provider uses 1Cloud's *paid util* estimation.

| Parameter | Type   | Is optional                              | Description             |
|-----------|--------|------------------------------------------|-------------------------|
| `name`    | string | Required                                 | Tile name               |
| `token`   | string | Required                                 | 1Cloud API access token |
| `url`     | string | Yes, defaults to `https://api.1cloud.ru` | 1Cloud API URL          |
| `timer`   | string | Yes, defaults to `60m`                   | Balance check interval  |
