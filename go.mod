module github.com/dh1tw/streamdeck-examples

go 1.14

replace github.com/karalabe/hid => github.com/dh1tw/hid v1.0.1-0.20190713013040-4404f3fb150e

replace github.com/dh1tw/streamdeck-buttons => /Users/user/go/src/github.com/dh1tw/streamdeck-buttons

require (
	github.com/dh1tw/hid v1.0.1-0.20190713013040-4404f3fb150e
	github.com/dh1tw/streamdeck v0.1.0
	github.com/dh1tw/streamdeck-buttons v0.0.0-00010101000000-000000000000
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/markbates/pkger v0.17.0
	github.com/spf13/cobra v1.0.0
)
