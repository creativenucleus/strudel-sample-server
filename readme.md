# Strudel Tasting Tray

Serves samples for Strudel locally. I tried a couple of other solutions, but they didn't suit.

This was built for me, so it's not refined for other users just now!

`strudel-tasting-tray.exe`

Starts a local server to provide samples to Strudel. Example:

```cli
.\strudel-tasting-tray.exe --port 5000 --sources "banginsamples<-./testdata/samplepack/strudel.json"
```

```strudel
samples('http://localhost:5000/banginsamples')

$: s("noice")
```

## Arguments

### --port (optional)

`--port 1234`

### --sources

`--sources "alias<-/path/to/strudel.json"`

Multiple sources can be provided, and these will be served from top-level folders from the endpoint.

`--sources "mybreaks<-/breaks/pathto/strudel.json" --sources "myvox<-/vox/pathto/strudel.json"`
