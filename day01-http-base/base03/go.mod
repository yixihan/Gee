module "Gee/day01-http-base/base03"

go 1.18

require gee v0.0.1

// 使用 replace 将 gee 指向 ./gee
replace gee => ./gee