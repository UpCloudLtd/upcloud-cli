package supabasechart

import (
	"embed"
)

// Note: go:embed charts will not embed _*.tpl files because of a bug or implementation detail: https://github.com/golang/go/issues/43854

//go:embed charts/supabase/*
//go:embed charts/supabase/templates/*
//go:embed charts/supabase/templates/*/*
var supabaseChartFS embed.FS

// ChartFS exposes the embedded supabase chart files.
var ChartFS = supabaseChartFS
