# Microb terminal client

Terminal client for [Microb](https://github.com/synw/microb) services

### Configure

Configure your servers in `config.json`:

   ```javascript
   {
	"servers":[
		{
		"name": "localhost",
		"centrifugo_addr":"localhost:8001",
		"centrifugo_key":"secret_key"
		},
		{
		"name": "another_one",
		"centrifugo_addr":"localhost:8001",
		"centrifugo_key":"secret_key"
		}
	],
	"services":["infos", "http", "mail"]
   }
   ```

### Run

   ```bash
./microb-cli
# or
./microb-cli -u=localhost
# start using localhost server
   ```

## Available commands

- **use**: server to use

- **using**: server beeing used

### Available services:

- [Http](https://github.com/synw/microb-http): an http server
- [Mail](https://github.com/synw/microb-mail): send mails

Credits
-------

- [Go prompt](https://github.com/c-bata/go-prompt): library to handle the terminal prompt

