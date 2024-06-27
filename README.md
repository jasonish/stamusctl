# Stamus control

## Description
stamusctl is a Command-Line Interface application written in GoLang by Stamus Networks that provides various functionalities to:
- manage Stamus stack configuration files
- deploy Stamus stack

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
stamusctl [commands] [flags] [args]
```
If not, you can:
```
./stamusctl [commands] [flags] [args]
```

## Commands

### Compose
`stamusctl compose` is the command used to manage the containerized Stamus stack deployement.

- `init` is used to initiate configuration files
  - `--folder ` to select folder to save configuration
  - `--default` to disable the interactive prompting
    - `[key]=[value]` are the args to set configuration values (ex: scirius.token=AwesomeToken)
- `config` is used to modify configuration files
  - `get` to display the current configuration values
    - `[key]` to get specific configuration values (ex: scirius)
  - `set` to modify configuration
    - `--reload` to reset arbitrary values
    - `[key]=[value]` to set configuration values (ex: scirius.token=AwesomeToken)
- `up` to start current configuration
  - `--file` to select specific configuration
  - `--detach` to launch in detached mode
- `down` to stop current configuration
  - `--file` to select specific configuration
  - `--volumes` to remove named volumes
  - `--remove-orphans` to remove not defined containers
- `update` to update the configuration
  - `--folder` to select folder
  - `--registry` to select registry
  - `--user` to input user
  - `--pass` to input password
  - `--version` to input version


## Examples
```
// Init via user prompt
stamusctl compose init

// Init default settings
stamusctl compose init --default scirius.token=AwesomeToken

// Get current config
stamusctl compose config get

// Set parameter in current config
stamusctl compose config set scirius.token=AnotherAwesomeToken

// Start current configuration
stamusctl compose up -d
```

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request.

## License
This project is licensed under the GPL3 License.