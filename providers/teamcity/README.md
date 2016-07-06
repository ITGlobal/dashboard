# [teamcity] provider

A provider for Dash that displays statuses of TeamCity builds.
Provider will populate one dashboard item for each TeamCity project which will indicate current build status.

## Configuration

Provider key: `teamcity`.

Provider parameters:

* `url`   - TeamCity URL. This parameter is required and should contain a valid URL.
* `username`        - username for TeamCity access. This parameter is optional.
* `password`        - password for TeamCity access. This parameter is optional.
* `fetch-proj-time` - project list refresh timer period. This parameter is optional and will default to `10m` if ommited.
* `fetch-build-time`- build status refresh timer period. This parameter is optional and will default to `20s` if ommited.

Provider will connectto TeamCity as guest unless both `username` and `password` parameters are specified.