dashboard terminal
==================

Dashboard app for Raspberry Pi - a terminal client app

Compile
-------

Run [`build.cmd`](/../build.cmd) or [`build.sh`](/../build.sh) to compile all parts of **dashboard** including **dasht**

Alternalively, you can build only daemon using the following commands:

```bash
go get -u github.com/itglobal/dashboard/terminal
cd $GOPATH/src/github.com/itglobal/dashboard/terminal
go build -o dasht .
```

Run
---

Execute the following command to start terminal client:

```bash
dasht [--url URL] [-colors COLORS] [-style STYLE]
```

### Parameters

* `URL` points to `dashd`'s endpoint. Default value is `http://127.0.0.0:8000`.
* `COLORS` might be either `dark` or `light`. Default value is `dark`.
* `STYLE` might be either `single` or `double`. Default value is `single`.

A running client can be stopped by pressing **ESC**.
