module main.go
go 1.22
replace gopkg.in/alecthomas/kingpin.v2 => github.com/alecthomas/kingpin/v2 v2.4.0
replace github.com/u-root/u-root => github.com/u-root/u-root v0.14.1-0.20241107071304-f908619d0238
require (
    github.com/alecthomas/kingpin/v2 v2.4.0
    github.com/u-root/u-root v0.14.1-0.20241107071304-f908619d0238
    // Add other dependencies here
)
