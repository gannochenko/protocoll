# Protocoll

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Protocoll is used to generate Postman/Insomnia collection out of a set of proto files.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the application, run the installation script. The script will install the most recent version.

The application is only available for Linux, MacOS and Windows, amd64 CPU architecture. The script will choose the operating system and the CPU architecture automatically.

~~~bash
curl -sSL https://raw.githubusercontent.com/gannochenko/protocoll/installer/script/install.sh | sh
~~~

After the script is done with its job, add the following path to the PATH environment variable:

~~~
export PATH=~/.protocoll/bin:$PATH;
~~~

The previous versions are available for download on [the release pages](https://github.com/gannochenko/protocoll/releases).

The application can be compiled for any other architecture Golang supports.

## Usage

To generate the collection the following command is used. The output will be written to stdout.
To have the collection written to a file instead, do output redirect:

~~~bash
protocoll generate --folder <path-to-protobuf-files> --name "My collection" > collection.json
~~~

## Contributing

1. [Install Golang](https://go.dev/doc/install).
2. Clone the repository:
   ~~~bash
   git clone git@github.com:gannochenko/protocoll.git
   cd protocoll
   ~~~
3. Install the dependencies:
   ~~~bash
   go mod download
   ~~~
4. Run the application with the demo data:
   ~~~bash
   make run
   ~~~
5. To make a release:
   ~~~bash
   make release version=1.0.0
   ~~~

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT). Explain the license terms and provide a link to the full license file.
