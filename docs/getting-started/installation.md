# Installation

In order to get the Sunbird interpreter up and running you need to compile it from source.
```sh
git clone https://github.com/radeqq007/sunbird.git
cd sunbird
go build ./cmd/sunbird -o sunbird
```

After that you should be able to run the interpreter by typing `./sunbird` in the terminal.

To run a script you just need to provide the path to the script:
```sh
./sunbird run main.sb
```

Or you can initialize a package:
```sh
./sunbird init
```

this will generate a `sunbird.toml` file containing project information and a `src/main.sb` file.

Then you can run it with:
```sh
./sunbird run
```

which will take the main file specified in `sunbird.toml`
