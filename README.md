# This is a WIP project

# git-helper-cli

A CLI tool to automate git related tasks by using OpenAI's completion API.

## Roadmap

- [Roadmap Page](docs/roadmap.md)

## Usage

```bash
git-helper-cli
```

## Configuration

Create a TOML configuration file at `~/.config/git-helper-cli/config.toml` with the following structure:

```toml
[api]
api_endpoint = "https://api.openai.com/v1" #default
api_secret = "your_openai_api_secret" #required

[branch]
pattern = "${date}/feature/${description}" #default
description_format = "kebab-case" #default
max_description_length = 50 #default
num_suggestions = 10 #default
```

### Configuration Options

- `api.api_endpoint`: The OpenAI API endpoint (default: "https://api.openai.com/v1")
- `api.api_secret`: Your OpenAI API secret key
- `api.model`: The OpenAI model to use (default: "gpt-4o-mini")
- `branch.pattern`: The pattern for generating branch names. Use `${date}` and `${description}` as placeholders.
- `branch.description_format`: The format for the description part of the branch name (currently supports "kebab-case")
- `branch.max_description_length`: Maximum length for the description part of the branch name
- `branch.num_suggestions`: Number of branch name suggestions to generate
