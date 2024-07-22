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
  - `--values` to use a `values.yaml` as configuration
  - `--fromFile` to use a file as value for a specific key
  - `--apply` to relaunch the configuration
  - `--default` to disable the interactive prompting
    - `[key]=[value]` are the args to set configuration values (ex: scirius.token=AwesomeToken)
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

### Config
`stamusctl config` is the command used to manage the configuration files.

- `config` is used to modify configuration files
  - `get` to display the current configuration values
    - `[key]` to get specific configuration values (ex: scirius)
    - `content` to get configuration architecture
      - `[key]` to get specific configuration files (ex: nginx)
  - `set` to modify configuration
    - `--reload` to reset arbitrary values
    - `--apply` to relaunch the configuration
    - `[key]=[value]` to set configuration values (ex: scirius.token=AwesomeToken)
    - `content` to set configuration files
      - `[host folder]:[configuration folder]` to set specific configuration files (ex: ./nginx:/nginx)

### Login
`stamusctl login` is the command used to login to image registries.
This way, you can pull images from private registries.

- `login` is used to login to image registries
  - `--registry` to select registry
  - `--user` to input user
  - `--pass` to input password
  

## Examples
```
// Init via user prompt
stamusctl compose init

// Init default settings
stamusctl compose init --default scirius.token=AwesomeToken

// Get current config
stamusctl config get

// Set parameter in current config
stamusctl config set scirius.token=AnotherAwesomeToken

// Start current configuration
stamusctl compose up -d
```

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request.

## License
This project is licensed under the GPL3 License.