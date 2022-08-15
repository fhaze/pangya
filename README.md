# Go Pangya Server

Just another pangya server, but this time created for the Go Language.

## 🤔 Considerations

This project is an experiment and probably will not be finished.

## 📦 Requirements

- go 1.18
- go-migrate cli
- tmux
- docker

## 💻 Environment Setup

I'm currently using WSL with Ubuntu, but it should work for any Linux distro and macOS as well. Windows is not supported and never will (not by me 😉).


Create a `dotenv` from template

```bash
cp ./.env-template ./.env
```

Downloads deps

```bash
make deps
```

Build everything

```bash
make build
```

## 🚀 Running Locally

I created a make target called `run` to quickly build and run all the servers locally. You should use it, cuz it's simpler and convenient

```bash
make run
```

> 📢 All servers are executed inside tmux, just tap `CTRL+C` to close each server.

## License

[WTFPL](https://choosealicense.com/licenses/wtfpl/)
