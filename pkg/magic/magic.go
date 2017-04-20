package magic

import "regexp"

var (
	Matcher = regexp.MustCompile(`(?m)^\s*(%%?[A-Za-z_][A-Za-z_0-9]*\b)\s*(.*)\s*$`)
)

// see https://github.com/JuliaLang/IJulia.jl/blob/master/src/magics.jl
