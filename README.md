dashboard
=========

Dashboard app for Raspberry Pi

![Preview](/preview.png)

Install
-------

Run command
```shell
$ go get -u github.com/itglobal/dashboard
```

Configure
---------

Create a JSON file somwhere  with the following structure:
```javascript
{
  "theme" : {
    "style"  : "single",
    "colors" : "light"
  },
  "providers": [ 
     {
        "type" : "PROVIDER KEY",
        // Place provider-specific parameters here
     },
     {
        "type" : "PROVIDER KEY",
        // Place provider-specific parameters here
     },
     // ...
  ]
}
```

* Field `theme:style` defines style of dashboard items. Allowed values: `single` for single-line borders, `double` for double-line borders.
* Field `theme:colors` defines colors scheme for dashboard. Allowed values: `light`, `dark`.
* Fields `providers` defines list of dashboard data providers

### Dashboard data providers:

 * [sim](/providers/sim/README.md) 
 * [ping](/providers/ping/README.md) 
 * [mongodb](/providers/mongodb/README.md) 
 * [teamcity](/providers/teamcity/README.md) 


Run
---

Run command
```shell
$ dashboard --config CONFIG_FILE
```

Hit `ESC` key to exit
