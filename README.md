# id3ed

Edit ID3 tags as JSON in vim.

> ID3 is a metadata container most often used in conjunction with the MP3 audio file format. It allows information such as the title, artist, album, track number, and other information about the file to be stored in the file itself. ([Wikipedia](https://en.wikipedia.org/wiki/ID3))

## Installation

```console
$ go install ./cmd/id3ed
```

## Usage

```console
$ id3ed --help
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

Modify them, save the file, then quit the editor.

Accepts [JWCC](https://nigeltao.github.io/blog/2021/json-with-commas-comments.html).

### Scripts

Edit every file in the present working directory that doesn't have a set title.

```bash
for file in *; do
    len=$(id3ed --inspect "$file" | jq '.title | length')
    if [ $len == "0" ]; then
        id3ed "$file";
    fi
done
```