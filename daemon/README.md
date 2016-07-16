dashboard daemon
================

Dashboard app for Raspberry Pi - a daemon app

Compile
-------

Run [`build.cmd`](/../build.cmd) or [`build.sh`](/../build.sh) to compile all parts of **dashboard** including **dashd**

Alternalively, you can build only daemon using the following commands:

```bash
go get -u github.com/itglobal/dashboard/daemon
cd $GOPATH/src/github.com/itglobal/dashboard/daemon
go build -o dashd .
```

Configure
---------

Create a JSON file somewhere  with the following structure:

```json
[ {
    "type" : "PROVIDER KEY",
    // Place provider-specific parameters here
  }, {
    "type" : "PROVIDER KEY",
    // Place provider-specific parameters here
  },
  // ...
]
```

Each element in this list defines a dashboard data provider.

### Dashboard data providers

* [sim](/../providers/sim/README.md)
* [ping](/../providers/ping/README.md)
* [mongodb](/../providers/mongodb/README.md)
* [teamcity](/../providers/teamcity/README.md)


Run
---

Execute the following command to start daemon:

```bash
dashd [-addr ENDPOINT] [-config CONFIG_FILE]
```

### Parameters

* `ENDPOINT` specified address to listen. Default value is `0.0.0.0:8000`.
* `CONFIG_FILE` indicates path to config file. Default value is `./dashd.json`.

A running daemon can be stopped by pressing **Ctrl+C**.