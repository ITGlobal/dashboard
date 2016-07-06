# [ping] provider

A provider for Dash that checks specified URL periodically.

If the specified endpoint will respond to HTTP GET with any 2xx or 3xx HTTP status then the dashboard item will be marked as `SUCCESS`.
Otherwise, it will be marked as `ERROR`.

## Configuration

Provider key: `ping`.

Provider parameters:

* `url`   - endpoint to check. This parameter is required and should contain a valid URL.
* `timer` - endpoint check timer period. This parameter is optional and will default to `1m` if ommited.