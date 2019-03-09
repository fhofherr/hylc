# hylc

`hylc` is login and consent provider for [Ory
Hydra](https://github.com/ory/hydra). Its main purpose is to serve as
a testbed for learning and evaluating Hydra. `hylc` is also intended to
be safe for production use for cases where its limited functionality is
all that is necessary.

## Features

`hylc` implements the minimal set of features necessary to make it an
useful login and consent provider for Hydra.

### Login Provider

`hylc` serves as a login provider for Hydra. The login functionality is
available under the `/login` path.

### Consent Provider

`hylc` serves as a login provider for Hydra. The login functionality is
available under the `/consent` path.

### Error Endpoint

`hylc` provides an error endpoint which displays login related error
messages to the end user.

### User management

`hylc` provides a simple commandline interface for user management.

## Deploy `hylc`

`hylc` is packaged as a Docker image. It is not yet available on docker
hub. Users must build the image themselves using `make image`.

## License

Copyright © 2019 Ferdinand Hofherr

Distributed under the MIT License.
