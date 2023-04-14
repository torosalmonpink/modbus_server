Modbus Server
=============
This is a simple Modbus server implementation in Go. It listens for Modbus TCP connections and handles incoming requests from clients.

Features
--------
Logging using logrus
Command-line options for listening address and port using spf13/pflag
Built with the torosalmonpink/mbserver library for handling Modbus communication

Usage
-----
To build and run the server, clone the repository and run the following commands:

```sh
git clone https://github.com/torosalmonpink/modbus_server
cd modbus_server
go build ./modbus_server
```
By default, the server listens on address 0.0.0.0 and port 502.
To change these settings, use the --address (or -a) and --port (or -p) command-line options:

```sh
./modbus_server --address 192.168.1.100 --port 1502
```

License
-------
This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.