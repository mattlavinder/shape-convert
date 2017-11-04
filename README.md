Shape File Convert
==================

##THIS IS A WORK IN PROGRESS

Shape file conversion tool for BRL Mapping.

Requires go-shp package. Use following command to install:

    go get github.com/jonas-p/go-shp

Command usage:

| Flag | Description            | Example         |
|------|------------------------|-----------------|
| i    | File to convert        | `-i alaska.shp` |
| o    | File to create         | `-o alaska.csv` |
| c    | Convert all centers    | `-c`            |
| z    | Convert all zip codes  | `-z`            |
| r    | Convert all routes     | `-r`            |
| v    | Verbose output         | `-v`            |
