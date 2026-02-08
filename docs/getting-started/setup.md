# Setting Up a Project

You can initialize a project with `sunbird init`, which creates a `sunbird.toml` file and a `src` directory for your code.

```bash
sunbird init
```

`sunbird.toml` contains metadata about your project, like its name and version.

```toml
[package]
name = "my_project"
version = "0.1.0"
description = "A Sunbird project"
authors = ["Your Name <you@example.com>"]
main = "./src/main.sb"

dependencies = []
```

The `main` field specifies the entry point of your application.
The main entry will run when you execute `sunbird run`.
You can also specify a file to run with `sunbird run <file>` (this doesn't require a `sunbird.toml` file, good for quick tests or one-off scripts).

The `src` directory is where you put your source code. You can create multiple `.sb` files and import them as needed.

