# id3ed

Edit ID3 tags as JSON in vim.

> ID3 is a metadata container most often used in conjunction with the MP3 audio file format. It allows information such as the title, artist, album, track number, and other information about the file to be stored in the file itself. ([Wikipedia](https://en.wikipedia.org/wiki/ID3))

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
