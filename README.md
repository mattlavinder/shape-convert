Shape File Convert
==================

**THIS IS A WORK IN PROGRESS**

Shape file conversion tool.

Requires go-shp package. Use following command to install:

    go get github.com/jonas-p/go-shp

Command Flags
-------------

| Flag | Description            | Example         |
|------|------------------------|-----------------|
| i    | Input file             | `-i alaska.shp` |
| o    | Output file            | `-o alaska.txt` |
| c    | Convert centers        | `-c`            |
| z    | Convert zip codes      | `-z`            |
| r    | Convert carrier route  | `-r`            |
| b    | Batch input directory  | `-b shape`      |
| v    | Verbose output         | `-v`            |

Command Usage
-------------

**Convert all routes**

    shape-convert -b 10-17_cr_shpGEN -o cr.txt -r

**Convert all zip codes**

    shape-convert -b 10-17_zip_shpGEN -o zip.txt -z

**Convert route centroids**

    shape-convert -i 10-17_cr_shpGEN/pcr_us.txt -o routecenters.txt -c

**Convert zip centroids**

    shape-convert -i 10-17_zip_shpGEN/zip_us.txt -o zipcenters.txt -c
