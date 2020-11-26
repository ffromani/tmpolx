# TMPolX: Topology ManagerPolixy eXplanation tool

A simple tool to test how kubernetes' topology manager behave, without need to run the real workload.
This tool uses the very same packages from upstream kubernetes, to give the closest representation as
possible as the real thing.

Usage example:
```bash
$ tmpolx -N 0-1 -P restricted \
	'{"R":"cpu", "H":[{"M":"01","P":true},{"M":"10","P":true},{"M":"11","P":false}]}' \
	'{"R":"nvidia.com/gpu", "H":[{"M":"01","P":true},{"M":"11","P":false}]}' \
	'{"R":"openshift.io/intelsriov", "H":[{"M":"10","P":true},{"M":"11","P":false}]}'
{restricted: map[cpu:[{01 true} {10 true} {11 false}] nvidia.com/gpu:[{01 true} {11 false}] openshift.io/intelsriov:[{10 true} {11 false}]]}
admit=false hint={01 false}
$
```

## license
(C) 2020 Red Hat Inc and licensed under the Apache License v2

## build
just run
```bash
make
```
