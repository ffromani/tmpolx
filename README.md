# TMPolX: Topology Manager Policy eXploration tool

A simple tool to test how kubernetes' topology manager behave, without need to run the real workload.
This tool uses the very same packages from upstream kubernetes, to give the closest representation as
possible as the real thing.

Usage example:
```bash
$ tmpolx -N 0-1 -P restricted \
	'nvidia.com/gpu:[{01 true} {11 false}]' \
	'openshift.io/intelsriov:[{10 true} {11 false}]' \
	'cpu:[{01 true} {10 true} {11 false}]'
using policy "restricted"
.	resource		hints				
.	nvidia.com/gpu		[{01 true} {11 false}]		
.	openshift.io/intelsriov	[{10 true} {11 false}]		
.	cpu			[{01 true} {10 true} {11 false}]
admit=false hint={01 false}
$ tmpolx -J -N 0-1 -P restricted \
	'{"R":"cpu", "H":[{"M":"01","P":true},{"M":"10","P":true},{"M":"11","P":false}]}' \
	'{"R":"nvidia.com/gpu", "H":[{"M":"01","P":true},{"M":"11","P":false}]}' \
	'{"R":"openshift.io/intelsriov", "H":[{"M":"10","P":true},{"M":"11","P":false}]}'
using policy "restricted"
.	resource		hints				
.	nvidia.com/gpu		[{01 true} {11 false}]		
.	openshift.io/intelsriov	[{10 true} {11 false}]		
.	cpu			[{01 true} {10 true} {11 false}]
admit=false hint={01 false}
$
```

## Command line explanation

`tmpolx` expects to receive positional arguments which represent the topology manager hints.
```json
{
	"R": "something", # resource name
	"H": [
		# array of hints
		{
			"M": "1111", # bitmask
			"P": true    # preferred flag
		}
	]
}
```

`tmpolx` accepts topology manager hints both in native go format (just copy/paste them from the kubelet logs) or in JSON format.
The go format is handier and simpler to use, but the support is still experimental. The JSON format is recommended to get
the maximum safety. The default is to use the go format.

### Topology Hints in JSON format

Use the `-J` flag to enable this format.

```bash
cpu:[{01 true} {10 true} {11 false}]
```
is
```
resource = "cpu"
hints = 
	{01 true}
	{10 true}
	{11 false}
# each hint is {mask preferred}
```
so we get
```json
{
	"R": "cpu",
	"H":[
		{
			"M": "01",
			"P": true
		},
		{
			"M": "10",
			"P": true
		},
		{
			"M": "11",
			"P": false
		}
	]
}
```

## license
(C) 2020 Red Hat Inc and licensed under the Apache License v2

## build
just run
```bash
make
```
