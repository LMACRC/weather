From https://github.com/gjr80/weewx-realtime_clientraw/blob/6cd06d8222a4a2229a333ed52e0db3bf671ff5a5/bin/user/rtcr.py#L954

| Field # | Name  | Unit  | Comment |
| :----   | :---- | :---- | :----   |
| 000     | preamble |       | "12345"      |
| 001     | avg speed | knots      |       |
| 002     | gust  | knots      |       |
| 003     | windDir | Deg      |       |
| 004     | outTemp) | Celsius      |       |
| 005     | outHumidity | %      |       |
| 006     | barometer | hPa      |       |
| 007     | daily rain  | mm      |       |
| 008     | monthly rain | mm      |       |
| 009     | yearly rain | mm      |       |
| 010     | rain rate  | mm / min      |       |
| 011     | max daily rainRate | mm / min       |       |
| 012     | inTemp  | Celsius      |       |
| 013     | inHumidity | %      |       |
| 014     | soil temperature | Celsius      |       |
| 015     | Forecast Icon |       |       |
| 016     | WMR968 extra temperature | Celsius      |       |
| 017     | WMR968 extra humidity | Celsius      |       |
| 018     | WMR968 extra sensor | Celsius      |       |
| 019     | yesterday rain | mm      |       |
| 020     | extra temperature sensor 1 | Celsius      |       |
| 021     | extra temperature sensor 2 | Celsius      |       |
| 022     | extra temperature sensor 3 | Celsius      |       |
| 023     | extra temperature sensor 4 | Celsius      |       |
| 024     | extra temperature sensor 5 | Celsius      |       |
| 025     | extra temperature sensor 6 | Celsius      |       |
| 026     | extra humidity sensor 1 | %      |       |
| 027     | extra humidity sensor 2 | %      |       |
| 028     | extra humidity sensor 3 | %      |       |
| 029     | hour |       |       |
| 030     | minute |       |       |
| 031     | seconds |       |       |
| 032     | station name |       |       |
| 033     | dallas lightning count |       |       |
| 034     | Solar Reading  |       | used as 'solar percent' in Saratoga dashboards      |
| 035     | Day |       |       |
| 036     | Month |       |       |
| 037     | WMR968/200 battery 1 |       |       |
| 038     | WMR968/200 battery 2 |       |       |
| 039     | WMR968/200 battery 3 |       |       |
| 040     | WMR968/200 battery 4 |       |       |
| 041     | WMR968/200 battery 5 |       |       |
| 042     | WMR968/200 battery 6 |       |       |
| 043     | WMR968/200 battery 7 |       |       |
| 044     | windchill | Celsius      |       |
| 045     | humidex | Celsius      |       |
| 046     | maximum day temperature | Celsius      |       |
| 047     | minimum day temperature | Celsius      |       |
| 048     | icon type |       |       |
| 049     | weather description |       |       |
| 050     | barometer trend | hPa      |       |
| 051     | windspeed hour 01-20 incl | knots      |       |
| 070     | windspeed hour (20) |       |       |
| 071     | maximum wind gust today |       |       |
| 072     | dewpoint  | Celsius      |       |
| 073     | cloud height  | foot      |       |
| 074     | date  | dd/mm/yyyy      |       |
| 075     | maximum day humidex  | Celsius      |       |
| 076     | minimum day humidex  | Celsius      |       |
| 077     | maximum day windchill  | Celsius      |       |
| 078     | minimum day windchill  | Celsius      |       |
| 079     | davis vp UV |       |       |
| 080     | hour wind speed 01-10 |       |       |
| 089     | hour wind speed (last) |       |       |
| 090     | hour temperature 01 | Celsius      |       |
| 091     | hour temperature 02-10 | Celsius      |       |
| 099     | hour temperature 10 | Celsius      |       |
| 100     | hour rain 01-10 | mm      |       |
| 109     | hour rain (10) | mm      |       |
| 110     | maximum day heatindex | Celsius      |       |
| 111     | minimum day heatindex | Celsius      |       |
| 112     | heatindex | Celsius      |       |
| 113     | maximum average speed | knot      |       |
| 114     | lightning count in last minute |       |       |
| 115     | time of last lightning strike  |       |       |
| 116     | date of last lightning strike  |       |       |
| 117     | wind average direction |       |       |
| 118     | nexstorm distance |       |       |
| 119     | nexstorm bearing  |       |       |
| 120     | extra temperature sensor 7 | Celsius      |       |
| 121     | extra temperature sensor 8 | Celsius      |       |
| 122     | extra humidity sensor 4 |       |       |
| 123     | extra humidity sensor 5 |       |       |
| 124     | extra humidity sensor 6 |       |       |
| 125     | extra humidity sensor 7 |       |       |
| 126     | extra humidity sensor 8 |       |       |
| 127     | vp solar |       |       |
| 128     | maximum inTemp (Celsius) |       |       |
| 129     | minimum inTemp (Celsius) |       |       |
| 130     | appTemp (Celsius) |       |       |
| 131     | maximum barometer (hPa) |       |       |
| 132     | minimum barometer (hPa) |       |       |
| 133     | maximum windGust last hour (knot) |       |       |
| 134     | maximum windGust in last hour time |       |       |
| 135     | maximum windGust today time |       |       |
| 136     | maximum day appTemp (Celsius) |       |       |
| 137     | minimum day appTemp (Celsius) |       |       |
| 138     | maximum day dewpoint (Celsius) |       |       |
| 139     | minimum day dewpoint (Celsius) |       |       |
| 140     | maximum windGust in last minute (knot) |       |       |
| 141     | current year |       |       |
| 142     | THSWS - will not implement |       |       |
| 143     | outTemp trend (logic) |       |       |
| 144     | outHumidity trend (logic) |       |       |
| 145     | humidex trend (logic) |       |       |
| 146     | hour wind direction 01-10 - will not implement |       |       |
| 155     | hour wind direction (10) - will not implement |       |       |
| 156     | leaf wetness |       |       |
| 157     | soil moisture |       |       |
| 158     | 10 minute average wind speed (knot) |       |       |
| 159     | wet bulb temperature (Celsius) |       |       |
| 160     | latitude (-ve for south) |       |       |
| 161     |  longitude (-ve for east) |       |       |
| 162     | 9am reset rainfall total (mm) |       |       |
| 163     | high day outHumidity |       |       |
| 164     | low day outHumidity |       |       |
| 165     | midnight rain reset total (mm) |       |       |
| 166     | low day windchill time |       |       |
| 167     | Current Cost Channel 1 - will not implement |       |       |
| 168     | Current Cost Channel 2 - will not implement |       |       |
| 169     | Current Cost Channel 3 - will not implement |       |       |
| 170     | Current Cost Channel 4 - will not implement |       |       |
| 171     | Current Cost Channel 5 - will not implement |       |       |
| 172     | Current Cost Channel 6 - will not implement |       |       |
| 173     | day windrun |       |       |
| 174     | record end (WD Version) |       | "!!EOR!!"      |
