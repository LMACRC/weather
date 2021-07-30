# weather

## About

`weather` is a program designed to receive weather data from a personal weather station (PWS) and upload a CumulusMX
compatible realtime.txt file via FTP. Optionally, when running on a [Raspberry Pi][RPi], `weather` can capture images
from the camera module and upload these via FTP also, to serve as a webcam image.


## Features 

`weather` is written in Go and was designed to run on low-powered devices like a [Raspberry Pi][RPi].


## Usage

`weather` has been tested with a rebranded [Fine Offset WH2900][WH2900]. The WH2900 is a Wi-Fi enabled PWS, that can be
configured to send Ecowitt HTTP requests to a custom server. Configuring a custom server is done using the WS View app.
Navigate to the Custom Server page and specify the IP address and port of the computer running the `weather`
program in the custom server option. See the `server` section in the [weather.toml](etc/weather.toml) to determine 
the port.


## Configuration

`weather` is configured via a file named `weather.toml` in a format similar to INI called [TOML](https://toml.io).
See [weather.toml](etc/weather.toml) for a fully documented example of the available options.

## Modules

### Archive

The archive module is responsible for archiving historical weather observation data from the SQLite database to CSV 
files. 

[WH2900]: http://www.foshk.com/Wifi_Weather_Station/WH2900.html
[RPi]:    https://www.raspberrypi.org