<script setup>
import { data as libVersion } from '../version.data.js'
</script>

# Install

Latest version: **{{ libVersion }}**

## Go Install

This is the recommended way to use Goverter.

1. Install the binary

    ```bash-vue
    $ go install github.com/jmattheis/goverter/cmd/goverter@{{ libVersion }}
    ```

1. Run the binary.

    ```bash
    $ goverter --help
    ```

This method installs the binary inside your `$GOPATH/bin`, ensure that this
path is on your `$PATH`. 

## Go Run

You can `go run` goverter like this:

```bash-vue
$ go run github.com/jmattheis/goverter/cmd/goverter@{{ libVersion }} --help
```

This method is the easiest, as you don't have to install a binary on your
system. The command may take some time to execute, because Go has to compile
goverter before executing it. Go will cache the build process, but it may be
invalidated sometimes. 

## Dependency

1. Create a go modules project if you haven't done so already

    ```bash
    $ go mod init module-name
    ```

1. Add goverter as dependency:

    ```bash-vue
    $ go get github.com/jmattheis/goverter@{{ libVersion }}
    ```

1. Run the binary.

    ```bash
    $ go run github.com/jmattheis/goverter/cmd/goverter --help
    ```

This method allows you to have the goverter dependency defined inside the
`go.mod`. The benefit is that all developers will use the same goverter
version.
