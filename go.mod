module github.com/tlinden/ephemerup

go 1.18

require (
	github.com/alecthomas/repr v0.2.0
	github.com/gofiber/fiber/v2 v2.42.0
	github.com/gofiber/keyauth/v2 v2.1.32
	github.com/google/uuid v1.3.0
	github.com/knadh/koanf/parsers/hcl v0.1.0
	github.com/knadh/koanf/providers/env v0.1.0
	github.com/knadh/koanf/providers/file v0.1.0
	github.com/knadh/koanf/providers/posflag v0.1.0
	github.com/knadh/koanf/v2 v2.0.0
	github.com/spf13/pflag v1.0.5
	github.com/tlinden/ephemerup/common v0.0.0-00010101000000-000000000000
	go.etcd.io/bbolt v1.3.7
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20220530130905-52f3993e8d6d // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.44.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
)

replace github.com/tlinden/ephemerup/common => ./common
