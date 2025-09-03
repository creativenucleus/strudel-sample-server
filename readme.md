# Strudel Sample Server

A local sample server for Strudel. I tried a couple of other solutions, but they didn't suit.

This was built for me, so it's not refined for other users just now!

`strudel-sample-server.exe`

Starts a local server to provide samples to Strudel. Example:

```cli
.\strudel-sample-server.exe -port 5000 -sources /my/samples/strudel.json
```

```strudel
samples('http://localhost:5000')

$: s("noice")
```

## Arguments

### -port (optional)

`-port 1234`

### -sources

`-sources /path/to/strudel.json`

Currently, max one source, could be multiple in the future
