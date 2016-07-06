# [mongodb] provider

A provider for Dash that checks MongoDB replica status.

If all replica set nodes are healthy then the dashboard item will be marked as `SUCCESS`.
Otherwise, it will be marked as `ERROR`. 

## Configuration

Provider key: `mongodb`.

Provider parameters:

* `name`  - dashboard item name. This parameter is required.
* `url`   - MongoDB replica set connection URL. This parameter is required and should contain a valid URL.
* `timer` - endpoint check timer period. This parameter is optional and will default to `30s` if ommited.