# DumpSBS

Golang tool to connect to ADS-B receivers like [readsb](https://github.com/wiedehopf/readsb) or [dump1090](https://github.com/flightaware/dump1090) via [SBS protocol](http://woodair.net/sbs/article/barebones42_socket_data.htm) and dump data to files.
New file will be created hourly.

Add this to cron.daily to compress those files:

```sh
find /opt/dumpsbs/logs -mtime +1 -name '*.csv' -exec xz -3 {} \;
```

```go
type SBS struct {
	MessageType      string
	TransmissionType int
	SessionID        int       // With dump1090, this is always 1
	AircraftID       int       // With dump1090, this is always 1
	HexIdent         string    // Aircraft Mode S hexadecimal code
	FlightID         int       // With dump1090, this is always 1
	TimeGenerated    time.Time // Protocol does not contain timezone data
	TimeLogged       time.Time // Protocol does not contain timezone data

	CallSign          string // Eight digit flight ID. Can be anything at all
	Altitude          int    // Mode C altitude. Height relative to 1013.2mb (Flight Level). Not height AMSL.
	GroundSpeed       float64
	Track             float64 // In dump1090, this is the aircraft's heading. Elsewhere, the track of the craft derived from the velocity E/W and velocity N/S
	Latitude          float64
	Longitude         float64
	VerticalRate      int
	Squawk            string // Assigned Mode A Squawk code
	SquawkChangeAlert int    // Flag to indicate the squawk has changed
	Emergency         int    // Flag to indicate the emergency code has been set
	SPI               int    // Flag to indicate the transponder ident has been activated
	IsOnGround        int    // Flag to indicate the ground squat switch is active
}
```