import * as path from "path";
import { ExtensionContext, workspace } from "vscode";

import {
  LanguageClient,
  LanguageClientOptions,
  TransportKind
} from "vscode-languageclient/node";

let client: LanguageClient;
let isDev = process.env.NODE_ENV === "development";

export function activate(context: ExtensionContext) {
  // The server is implemented in node
  const serverPath = isDev ? context.asAbsolutePath(
    path.join("out", "server")
  ) : context.asAbsolutePath(
    path.join("..", "server", "bin", "lsp")
  );

  // Options to control the language client
  const clientOptions: LanguageClientOptions = {
    // Register the server for all documents by default
    documentSelector: [{ scheme: "file", language: "templ" }],
    synchronize: {
      // Notify the server about file changes to '.clientrc files contained in the workspace
      fileEvents: workspace.createFileSystemWatcher("**/.clientrc"),
    },

  };

  // Create the language client and start the client.
  client = new LanguageClient(
    "templ-html-htmx",
    "templ-html-htmx-language-server",
    {
      run: {
        command: serverPath,
        transport: TransportKind.stdio,
      },
      debug: {
        command: serverPath,
        transport: TransportKind.stdio,
      }
    },
    clientOptions
  );

  // Start the client. This will also launch the server
  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
