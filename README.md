# SymphonAI Serverless Backend

The backend infrastructure-as-code (IAC) and lambda code for the SymphonAI APIs.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Developer Setup

1. Clone this git repo to your machine
2. Install Golang (https://go.dev/doc/install)
3. In your command-line interface (CLI), run the following:
   ```
       go get
   ```
4. Install VS Code (https://code.visualstudio.com/)
5. Install the Go plugin and associated tools for VS Code (https://medium.com/backend-habit/setting-golang-plugin-on-vscode-for-autocomplete-and-auto-import-30bf5c58138a). **Ignore steps 2, 3 and 4 in this guide if you're on a Mac**

## Useful commands

- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `cdk synth` emits the synthesized CloudFormation template
- `go test` run unit tests
