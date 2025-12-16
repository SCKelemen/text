module github.com/SCKelemen/text

go 1.23

require (
	github.com/SCKelemen/unicode v0.0.0
	github.com/SCKelemen/units v0.0.0
)

replace (
	github.com/SCKelemen/unicode => ../unicode
	github.com/SCKelemen/units => ../units
)
