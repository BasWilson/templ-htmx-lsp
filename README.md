## What is it
I am busy implementing the Language Server Protocol so [Templ](https://templ.guide) has support for better code completion than current LSP. The goal it to support html, htmx and normal templ syntax and to be able to provide code completion for all of them.

The language server is written in go. Mainly to challenge myself to learn go. Coming from a Typescript background.

It is not published to the extension store, but you can easily install it by following these instructions.

## How to install

You can download the latest release from the [releases page](https://github.com/BasWilson/templ-htmx-lsp/releases)

1. Clone the repository

```bash
git clone https://github.com/baswilson/templ-lsp.git
```

2. Install the go dependencies

```bash
cd client && yarn install
```

```bash
cd server && go mod download
```

3. Install vsce globally

```bash
npm i -g vsce
```

4. Build the extension (server is built by client)
```bash
cd client && vsce package
```

5. Install the extension in VSCode
Right click the created file and select "Install Extension VSIX" or Run this command. Replace the version number with the correct version number.
```bash
code --install-extension templ-lsp-0.0.1.vsix
```

Done. You should now have the extension installed.