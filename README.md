<div align="center">
  <img src="./mocktopus.jpeg" width="180">

  <h1 align="center">
    mocktopus
  </h1>
  <p>
    🐙 <b>GPT powered</b> CLI tool to generate mocks for anything!
  </p>
</div>

### Installation

This project requires:

- Node v18 or above to be installed on the system.
- Your own OpenAI API key that can be found [here](https://platform.openai.com/account/api-keys).

```bash
git clone https://github.com/blenderskool/mocktopus
cd mocktopus
npm i
npm link
```

#### Add OpenAI API key as an env variable

This example is for `zsh` shell, you can add it accordingly for other shell environments.

```bash
nano ~/.zshrc
```

In the file that is opened, add the following line at the end

```bash
export MOCKTOPUS_OPENAI_KEY="<YOUR OPENAI API KEY>"
```

Save the file and exit, then restart the terminal.

### Usage

```
mocktopus [command]

Commands:
  proto [options] <source> <destination>  generate mock data for complex structures by analyzing proto definitions
  placeholder                             generate mock data from natural descriptions
  tests <source> <destination>            generate test cases for code snippets
  persona                                 generate user personas for a product
  help [command]                          display help for command
```

### Uninstall

```bash
npm unlink -g mocktopus
```

And you can optionally remove the env variables starting with `MOCKTOPUS`
