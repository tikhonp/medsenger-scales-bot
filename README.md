<!--suppress HtmlDeprecatedAttribute -->
<div align="center">
    <br>
    <h1>⚖️ Medsenger Xiaomi Mi Scales bot</h1>
</div>

The __GO__ Medsenger bot for Xiaomi Mi Scales. It allows you to get your weight and body composition data from the scales to the Medsenger chat.

# 📦 Development

1. Install __docker__ and __make__

2. Create configuration file on `.env`

### Run Development

```sh
make
```

or

```sh
make dev
```

or

```sh
make build-dev # preferred if config files were changed, so it rebuilds image
```

### HTML templating

I use [templ](https://github.com/a-h/templ) as template engine. After changing `*.templ` files regenerate go code using:

```sh
make templ
```

> development docker container must be active

### Enter server container shell

There is shortcut for this:

```sh
make go-to-server-container
```

# Deploying

To deploy you also need __docker__ and __make__. In project root run:

```sh
make prod
```

It will create prod containers and run it in detached mode.

To stop run:

```sh
make fprod
```

To view logs in real time:

```sh
make logs-prod
```

## 💼 License

Created by Tikhon Petrishchev

Copyright © 2024 OOO Telepat. All rights reserved.
