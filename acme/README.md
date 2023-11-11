# acme

## Usage

```sh
acme [ -abr ] [ -m mtpt ] [ -c ncol ] [ -f varfont ] [ -l file | file... ]
```

```sh
9p ls acme

# read contents of current window
9p read acme/$winid/body

# update plumbing rules
9p write plumb/rules < $PLAN9/plumb/rules
9p read plumb/rules

fontsrv &
acme -f /mnt/font/GoMono/15a/font
```
