# PHPublisher

A publisher in Golang!

## Download:

Download binary for your operation system and architecture [here](https://github.com/Pelegrinetti/phpublisher/releases/tag/v0.1.0-alpha).

## How to use it:

Make sure you have the PHPublisher available.

```sh
./phpublisher

# Hello, friend!
```

### Publishing a package:

You can publish any package using the command below:

```sh
# On root path of package:

./phpublisher publish --registry REPOSITORY_UPLOAD_URL --version VERSION --vendor VENDOR --project PROJECT_NAME --user USERNAME --password PASSWORD
```

### Ignoring files:

Create a file `.phpublisherignore` on project root and set files to be ignored on publish.

```sh
# Use regex pattern to ignore files.

vendor

\w+\.test\.\w+
```