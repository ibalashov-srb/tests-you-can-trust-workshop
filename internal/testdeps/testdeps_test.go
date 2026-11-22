// Package testdeps keeps the test libraries in go.mod while the repository has no tests yet.
package testdeps

import (
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
)
