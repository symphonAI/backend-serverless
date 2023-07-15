# SymphonAI Serverless Backend

The backend infrastructure-as-code (IAC) and lambda code for the SymphonAI APIs.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Developer Setup

1. Clone this git repo to your machine
2. Install Golang (https://go.dev/doc/install)
3. In your command-line interface (CLI), run the following in the root directory of this repository:
   ```
       go get
   ```
4. Install VS Code (https://code.visualstudio.com/)
5. Install the Go plugin and associated tools for VS Code (https://medium.com/backend-habit/setting-golang-plugin-on-vscode-for-autocomplete-and-auto-import-30bf5c58138a). **Ignore steps 2, 3 and 4 in this guide if you're on a Mac**
6. Install the AWS CDK toolkit: https://docs.aws.amazon.com/cdk/v2/guide/cli.html
7. Setup a AWS profile for deployment from your machine if you haven't already: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html

## Running Locally

You can try running the API locally using the SAM CLI, but no guarantees all the code (like Cognito access etc.) will run properly.

### Setup Steps for Running Locally

1. Carry out all the steps in the _Developer Setup_ section above.
2. Install Docker (https://www.docker.com/)
3. Install the SAM CLI (https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)
4. Create a `env.json` file in this root directory, with this format:

   ```
   {
       "Parameters": {
           "OPENAI_MODEL": "<<model>>",
           "OPENAI_API_KEY": "<<apikey>>",
           "SPOTIFY_CLIENT_ID": "<<spotifyclientid>",
           "SPOTIFY_CLIENT_SECRET": "<<spotifyclientsecret>>"
       }
   }
   ```

   Replace `<<model>>` with one of the following, depending on which model
   you want to use:

   - davinci
   - gpt-3.5-turbo
   - gpt-4

   Replace `<<apikey>>` with the ChatGPT API Key you want to use.

   `env.json` is in the .gitignore, but just in case: **Please do not check in this file to version control!**

   Replace `<<spotifyclientid>>` and `<<spotifyclientsecret>>` with the Spotify Developer App Client ID and App Client Secret.

### Steps to Run Locally

First follow the steps above (**Setup Steps for Running Locally**), then do the following:

1. Make sure Docker is running.

2. In the CLI, navigate to this directory and run: `make start`. This will expose the API on http://localhost:8080, so you will be able to hit API endpoints like http://localhost:8080/prompt.

**Note:** If the `make` command doesn't work on your machine, then you will need to install Make. The exact instructions vary depending on your OS, so you might have to Google....here are instructions for Ubuntu: https://linuxhint.com/install-make-ubuntu/

If you are unable to install make, then run the individual commands under `make start` that are listed in the `Makefile` contained within the root directory.

3. Press Ctrl+C to stop running the API on your machine.

## Deployment

Run `cdk deploy` to run the infrastructure/code changes against AWS.

**TODO: Automate deployment via Github actions**

## Other Useful CDK Commands

- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `cdk synth` emits the synthesized CloudFormation template
- `go test` run unit tests
