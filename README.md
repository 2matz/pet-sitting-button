pet sitting button
==================

How to use
----------

The package manager for this project is `dep`.
For instructions on how to install `dep`, please visit the `dep` [website](https://github.com/golang/dep).

```
$ make prepare
$ make build-linux-amd64
$ mv bin/pet-sitting-button_linux_amd64 bin/pet-sitting-button
$ zip bin/pet-sitting-button.zip bin/pet-sitting-button
```

and upload `bin/pet-sitting-button.zip` to aws as lambda function.