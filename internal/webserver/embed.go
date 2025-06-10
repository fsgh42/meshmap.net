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

// https://github.com/meshtastic/design/tree/master/Meshtastic%20Powered%20Logo
//
//go:embed m-pwrd_bw_noborder.png
var iconPoweredByMeshtastic []byte

//go:embed meshhessen-cropped-MH_nM-32x32.png
var iconMeshHessen32 []byte

//go:embed meshhessen-cropped-MH_nM-192x192.png
var iconMeshHessen192 []byte
