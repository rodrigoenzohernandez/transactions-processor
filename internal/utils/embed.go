// /internal/utils/embed.go
package utils

import (
	"embed"
)

//go:embed templates/balance.html
var Templates embed.FS
