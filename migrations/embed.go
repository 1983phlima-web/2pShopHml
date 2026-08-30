// Package migrations embeds the SQL migration files into the compiled
// binary so the application can apply them automatically at startup,
// without depending on a separate migration step or shell access in the
// (distroless) runtime image.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
