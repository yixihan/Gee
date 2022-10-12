module "Gee/day06-template"

go 1.18

require (
	gee v0.0.1
)

replace (
	gee => ./gee
)