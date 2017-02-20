# Microb terminal client

Terminal client for [Microb](https://github.com/synw/microb) servers

Configure your servers in `config.json`

   ```bash
./microb-cli
# or
./microb-cli -s=localhost
# start using localhost server
   ```

## Available commands

- **use**: server to use

- **using**: server beeing used

### Info

- **ping**

- **time**: measures server response times for an url

- **routes**: view server side routes currently used

### Database

- **get**: returns the content for url

- **db_status**: reports the database status

### Server state

- **reparse_templates**: refresh templates in memory

- **update_routes**: updates client side routes

Credits
-------

- [Ishell](https://github.com/abiosoft/ishell): library for building cli apps

