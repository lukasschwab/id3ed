# id3ed

Edit ID3 tags as JSON in vim. See [mikkyang/id3-go](https://github.com/mikkyang/id3-go)

## Usage

```console
$ go install ./cmd/id3ed
$ id3ed ~/Desktop/some-file.mp3
```

Opens an editor for the following metadata:

```json
{
    "title": "",
    "artist": "",
    "album": "",
    "year": "",
    "genre": ""
}
```

Accepts [JWCC](https://nigeltao.github.io/blog/2021/json-with-commas-comments.html).
