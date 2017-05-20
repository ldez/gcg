# GCG - GitHub Changelog Generator

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

Example:

```shell
./gcg -c"v1.3.0-rc1" -p"v1.2.0-rc1" -b"master" -o"containous" -r"traefik" -t"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

```shell
./gcg -c"v1.3.0-rc1" -p"v1.2.0-rc1" -b"master" -o"containous" -r"traefik" -t"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"\
--exclude-label="area/infrastructure" --enhancement-label="kind/enhancement" --doc-label="area/documentation" --bug-label="kind/bug/fix"
```
