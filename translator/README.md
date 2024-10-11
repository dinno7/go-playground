# Translator app

This cli app can help you to translate direct text or a file from a language to another one.

# How to use?

For getting help you can run

```bash
translate -h
```

### There is 2 way to translate you text:

1. translate direct text:

```bash
translate -text "Translate this text" --from en --to fa
```

2. Translate a whole file:

```bash
translate -path /path/of/yourfile.txt --from en --to fa
```

#### There is also `--sub` flag which help you translate your subtitle files in more efficient way:

**So if your file is a subtitle, provide `--sub` flag**

```bash
translate -path /path/of/your_subtitle_file.src --from en --to fa --sub
```

# Available options

`--from` or `-f` **string**
_The source language (default "en")_

`--to` or `-t` **string**
_The target language (default "fa")_

`--text` or `-x` **string**
_Text to translate_

`--path` or `-p` **string**
_File path to translate_

`--sub` or `-s` **boolean**
_If your file is a video subtitle, provide this option as true_

You can see complete list of language codes [HERE](https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes)
