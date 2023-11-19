module github.com/lucebac/winreg-tasks

go 1.21

require github.com/kaitai-io/kaitai_struct_go_runtime v0.10.0

require (
	github.com/alecthomas/kong v0.8.1
	github.com/google/uuid v1.4.0
	golang.org/x/sys v0.14.0
	golang.org/x/text v0.14.0
	www.velocidex.com/golang/regparser v0.0.0-20221020153526-bbc758cbd18b
)

require github.com/davecgh/go-spew v1.1.1 // indirect

// use my fork until https://github.com/Velocidex/regparser/pull/6 is merged
replace www.velocidex.com/golang/regparser => github.com/lucebac/regparser v0.0.0-20231118234522-9507c49eccc2
