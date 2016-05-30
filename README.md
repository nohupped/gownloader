# gownloader
A go download accelerator
Run as

```gownloader http://url.com 3```

This write to the filename, and each goroutine seeks to its respective offset and starts writing from there.
