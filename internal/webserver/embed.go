package webserver

import (
	_ "embed"
)

//go:embed meshmap.html
var meshmapHTML []byte

//go:embed site.webmanifest
var webmanifest []byte

//go:embed favicon.ico
var iconFav []byte

//go:embed android-chrome-192x192.png
var iconChrome192 []byte

//go:embed android-chrome-512x512.png
var iconChrome512 []byte

//go:embed apple-touch-icon.png
var iconAppleTouch []byte
