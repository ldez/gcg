# GCG - GitHub Changelog Generator

[![Build Status](https://travis-ci.org/ldez/gcg.svg?branch=master)](https://travis-ci.org/ldez/gcg)


```shell
GCG is a GitHub Changelog Generator.
	
Usage: gcg [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: gcg [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Flags:
-b, --base-branch       Base branch name. PR branch destination.                  (default "master")
    --bug-label         Bug Label.                                                (default "bug")
-c, --current-ref       Current commit reference. Can be a tag, a branch, a SHA.  
    --doc-label         Documentation Label.                                      (default "documentation")
    --enhancement-label Enhancement Label.                                        (default "enhancement")
    --exclude-label     Label to exclude.                                         
-f, --future-ref-name   The future name of the current reference.                 
    --output-type       Output destination type. (file|Stdout)                    (default "file")
-o, --owner             Repository owner.                                         
-p, --previous-ref      Previous commit reference. Can be a tag, a branch, a SHA. 
-r, --repo-name         Repository name.                                          
-t, --token             GitHub Token                                              
-h, --help              Print Help (this message) and exit  
```

## Examples

```bash
./gcg -c"v1.3.0-rc1" -p"v1.2.0-rc1" -o"containous" -r"traefik" -t"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

### Next major release (RC1 on master)

ex: for the (non-existing) next version 1.4.0-rc1
```bash
go run gcg.go -c"master" -p"v1.3.0-rc1" -f"v1.4.0-rc1" -b"master" \
-o"containous" -r"traefik" -t"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
--exclude-label="area/infrastructure" --enhancement-label="kind/enhancement" --doc-label="area/documentation" --bug-label="kind/bug/fix" \
--debug
```

### Next RC release (RC2 on a specific branch)

ex: for the (non-existing) version 1.3.0-rc2
```bash
go run gcg.go -c"v1.3" -p"v1.3.0-rc1" -b"v1.3" \
-o"containous" -r"traefik" -t"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
--exclude-label="area/infrastructure" --enhancement-label="kind/enhancement" --doc-label="area/documentation" --bug-label="kind/bug/fix" \
--debug
```

### Previous major release (Pre RC1)

ex: for the (existing) version 1.3.0-rc1
```bash
go run gcg.go -c"v1.3.0-rc1" -p"v1.2.0-rc1" -b"master" \
-o"containous" -r"traefik" -t"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
--exclude-label="area/infrastructure" --enhancement-label="kind/enhancement" --doc-label="area/documentation" --bug-label="kind/bug/fix" \
--debug
```

### Previous RC release (between RC1 and RC2)

ex: for the (existing) version 1.3.0-rc2
```bash
go run gcg.go -c"v1.3.0-rc2" -p"v1.3.0-rc1" -b"v1.3" \
-o"containous" -r"traefik" -t"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
--exclude-label="area/infrastructure" --enhancement-label="kind/enhancement" --doc-label="area/documentation" --bug-label="kind/bug/fix" \
--debug
```
