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

**Note: This project requires your own API key for one of the supported AI models:**

- OpenAI: Can be found [here](https://platform.openai.com/account/api-keys).
- Google Gemini: Can be found [here](https://aistudio.google.com/).

**Setup CLI**

1. Download the binary for your OS from the [releases page](https://github.com/blenderskool/mocktopus/releases/latest)
2. Rename the binary file to `mocktopus`
3. Copy the binary to `/usr/bin/` paths (for macOS, Linux). If you cannot copy the the binary to the directory, then update your `PATH` env variable to also include the directory where `mocktopus` binary is stored.
4. Proceed to adding the OpenAI/Gemini API key as an environment variable by following the steps in next section

#### Add API key of an AI Model as an env variable

This example is for `zsh` shell, you can add it accordingly for other shell environments.

```bash
nano ~/.zshrc
```

In the file that is opened, add the following line at the end

```bash
export MOCKTOPUS_OPENAI_KEY="<YOUR OPENAI API KEY>"
export MOCKTOPUS_GEMINI_KEY="<YOUR GEMINI API KEY>"
```

Save the file and exit, then restart the terminal.

### Usage

```
mocktopus [global options] command [command options] [arguments...]

Commands:
  proto        proto <source> <destination>
  placeholder
  tests        tests <source> <destination>
  persona
  help, h      Shows a list of commands or help for one command
```

### Uninstall

1. Remove the `mocktopus` binary file
2. Optionally remove the env variables starting with `MOCKTOPUS`
