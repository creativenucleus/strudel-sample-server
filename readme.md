# Strudel Sample Server

A local sample server for Strudel. I tried a couple of other solutions, but they didn't suit.

This was built for me, so it's not refined for other users just now!

`strudel-sample-server.exe`

Starts a local server to provide samples to Strudel. Example:

```cli
.\strudel-sample-server.exe --port 5000 --sources "banginsamples|/my/samples/strudel.json"
```

```strudel
samples('http://localhost:5000/banginsamples')

$: s("noice")
```

## Arguments

### --port (optional)

`--port 1234`

### --sources

`--sources "alias|/path/to/strudel.json"`

Multiple sources can be provided, and these will be served from top-level folders from the endpoint.
