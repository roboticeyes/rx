# rx - CLI for REX
[![Build Status](https://travis-ci.org/breiting/rx.svg?branch=master)](https://travis-ci.org/breiting/rx)

**rx** is a command line tool for accessing the REX cloud API. rx is using the [rex](https://github.com/breiting/rex) GO library.

Make sure that you have a valid API key which can be retrieved from https://rex.robotic-eyes.com.

## Compile or download your binary

If you have Go installed on your machine, you can directly build our executable by running:

```
go build
./rx
```

You can also download a pre-compiled version of `rx` from [here](https://github.com/breiting/rx/releases).

## Setup your config file

You need to create a configuration file in order to run `rx` properly. Create a new file (e.g. `rx.yml`) and add the following information:

```
BaseURL: "https://rex.robotic-eyes.com"
ClientID: "<your-client-id>"
ClientSecret: "<your-client-secret"
```

## Check authentication

To test your settings, you can run the following for getting your user information:

```
./rx --config rx.yml users --me
```

If you can see your user details, then the credentials are set properly.

## Create a new project

You can create a new project using the interactive shell:

```
./rx --config rx.yml projects new
```

The program will ask you several questions and stores the project in the REX cloud. Here is an example output of the interactive session. Please notice, that I am using OpenStreetMap for getting the geo-location (see license notice below). Based on the given address, OpenStreetMap may already propose a proper LAT/LON information.

Once the project is created, `rx` offers a link on the map, again using OpenStreetMap.

```
Name                : My first project
AddressLine1        : Stremayrgasse 16
AddressLine2        : 
AddressLine3        : 
AddressLine4        : 
Postcode            : 8010
City                : Graz
Region              : 
Country             : AT
OpenStreetMap data

LICENSE notice:    Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright
Geolocation name:  Technische Universität Graz, 16, Stremayrgasse, Schönau, Jakomini, Graz, Steiermark, 8010, Österreich
Geolocation lat:   47.0650332
Geolocation lon:   15.4522326353324

Lat [47.065033]     : 
Lon [15.452233]     : 
Height [m]          : 350
Northing (0=north, 90=east): 0


Project My first project added successfully!


Map link: https://www.openstreetmap.org?mlat=47.065033&mlon=15.452233#map=17/47.065033/15.452233&layers=N

```

For further information, please use `./rx --help`.
