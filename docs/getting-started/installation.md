# Installation

In order to get the Sunbird interpreter up and running you need to compile it from source.
```sh
git clone https://github.com/radeqq007/sunbird.git
cd sunbird
go build ./cmd/sunbird
```

After that you should be able to run the interpreter by typing `./sunbird` in the terminal.

To run a script you just need to provide the path to the script:
```sh
./sunbird main.sb
```