# Stamus control

## Description
stamusctl is a Command-Line Interface application written in GoLang that provides various functionalities to:
- manage stamus stack configuration files

It will also (soon) help:
- manage stamus stack deployement

## Installation
To install stamusctl, you can:
- build it from source
- get it from your favorite package manager

#### Build from source
```
git clone https://git.stamus-networks.com/devel/stamus-ctl
```
and follow [golang documentation](https://go.dev/doc/tutorial/compile-install)

## Usage
If you have the binary in your path, you can:
```
stamusclt [command] [flag] [args]
```
If not, you can:
```
./stamusclt [command] [flag] [args]
```

## Commands

```
stamusclt compose init
stamusclt compose config get
stamusclt compose config set
```

## Examples
```
// Init via user prompt
stamusclt compose init

// Init default settings
stamusclt compose init --default scirius.token=AwesomeToken

// Get current config
stamusclt compose config get

// Set parameter in current config
stamusclt compose config set scirius.token=AnotherAwesomeToken
```

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request.

## License
This project is licensed under the GPL3 License.