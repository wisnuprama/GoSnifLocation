# GoSnifLocation

> Where are you, scammers? :)

Use this to track scammer location.

Please use it ethically.

## Drawbacks

Doing cat and mouse with scammers is tricky. Using this tool, we actually *scam* back the scammers. This tool requires location permission, we will ask scammer to grant the location. We won't be able to get their location if they are too aware.

## How to use

### Installation

1. Make sure you have install Go
2. Install [Ngrok](https://ngrok.com/download) for proxy the local server to public

### Build & Run

1. Run `go build` and it will make a binary version `server`
2. Run `./server`

### Development

1. Run `go run .` to start server

## Example 

**Logging**

![Logging](./assets/logs.png "Logging")

**JSON File**

```json
{
  "geo_info": {
    "longitude": "13.105939",
    "latitude": "145.508881"
  },
  "ip_address": "[::1]:60316"
}
```