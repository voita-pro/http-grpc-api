package api

import (
	"embed"
	"io/fs"
)

//go:embed openapi/*
var content embed.FS
var EmbedSwagger, _ = fs.Sub(content, "openapi")
