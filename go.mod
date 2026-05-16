module github.com/persianopencart/fleetssl-cpanel-new

go 1.25.0

require (
	github.com/boltdb/bolt v1.3.1
	github.com/domainr/dnsr v0.0.0-20260501085227-a7299c7b79fd
	github.com/eggsampler/acme/v3 v3.8.1
	github.com/fatih/color v1.19.0
	github.com/fsnotify/fsnotify v1.10.1
	github.com/juju/ratelimit v1.0.2
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/kardianos/service v1.2.2
	github.com/letsencrypt-cpanel/cpanelgo v1.2.1
	github.com/sirupsen/logrus v1.9.4
	golang.org/x/crypto v0.51.0
	golang.org/x/net v0.54.0
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/ini.v1 v1.67.2
	gopkg.in/urfave/cli.v1 v1.20.0
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.72 // indirect
	golang.org/x/mod v0.35.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	golang.org/x/tools v0.44.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260511170946-3700d4141b60 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

// Dependencies we have hard-forked
replace github.com/kardianos/service => ./internal/github.com/kardianos/service
