Dashboard API
=============

Get current dashboard state
---------------------------

### Request

```http
GET /data.json
```

### Response

```json
{
    "version" : "VERSION_HASH",
    "tiles"   : [
        // Array of tiles
    ]
}
```

Field `tiles` contains tile objects.

Tile object
-----------

```json
{
    "id"          : "UNIQUE_TILE_ID",
    "source"      : "PROVIDER_TYPE",
    "updated"     : "LAST_TILE_UPDATE_TIME",
    "type"        : "TILE_TYPE",
    "size"        : "TILE_SIZE",
    "state"       :"TILE_STATE",
    "titleText"   : "TILE_TITLE_TEXT",
    "descrText"   : "TILE_DESCRIPTION_TEXT",
    "statusValue" : 100
}
```

| Field         | Type                 | Required | Decription                                     |
|---------------|----------------------|----------|------------------------------------------------|
| `id`          | `string`             | Yes      | Tile ID. It's is stable and should not change. |
| `source`      | `string`             | Yes      | Corresponding tile provider type.              |
| `updated`     | `datetime`           | Yes      | Tile's last update time.                       |
| `type`        | `string` (see below) | Yes      | Tile layout type.                              |
| `size`        | `string` (see below) | Yes      | Tile size.                                     |
| `state`       | `string` (see below) | Yes      | Tile state.                                    |
| `titleText`   | `string`             | Yes      | Tile title text.                               |
| `descrText`   | `string`             | No       | Tile description text.                         |
| `statusValue` | `number?`            | No       | Tile's progress bar position.                  |

> **Note:** `id` is stable and should not change

### Tile types

Different tile types support different fields.
More specifically, fields `titleText`, `descrText` and `statusValue` might be supported or not supported by various tile layouts.

| Type              | Supported fields         | Providers                   | Description                                               |
|-------------------|--------------------------|-----------------------------|-----------------------------------------------------------|
| `text`            | `titleText`              |                             | Tile with one line of text                                |
| `text-status`     | `titleText`, `descrText` | `check`                     | Tile with header and a line of text                       |
| `text-status-2`   | `titleText`, `descrText` | `1cloud`                    | Tile with small header and a large line of text           |
| `text-status-bar` | `titleText`, `descrText` | `sim`, `mongdb`, `teamcity` | Tile with header, a large line of text and a progress bar |

### Tile sizes

Valid values are:

* `1x` - tile that takes 1x1 grid cell
* `2x` - tile that takes 2x1 grid cells (horizontally)
* `4x` - tile that takes 2x2 grid cells

### Tile states

Valid values are:

* `default` - default *gray* tile
* `success` - success *green* tile
* `indeterminate` - indeterminate *cyan* tile
* `warning` - warning *orange* tile
* `error` - default *red* tile