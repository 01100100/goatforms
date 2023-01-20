# GoatForm

[GoatForm]() is a super light weight tool to collect form data POSTED to the `/forms` endpoint.

## Why

I was looking for a way to collect [form](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/form) responses from a [static site](https://en.wikipedia.org/wiki/Static_web_page) in the simplest and cheapest way possible, to use for a hobby project site that would get minimal traffic.

There are some good companies out there offering Faas (FormsAsAService).

- [Google Forms](https://www.google.com/forms/about/)
- [Netlify Forms](https://www.netlify.com/products/forms/)
- [Formspree](https://formspree.io/)

These services all offer a free tiers with a couple hundred form requests a month and features of data validation and anti spam. After the free tier runs out prices jump up pretty high. This may be fine for enterprise customers, but for hobby projects it seems to much.

I couldn't find anything to fit what I was looking for so built this. It can easily be deployed and self hosted cheaply using one of many services.

## Code

The program is written in [go](https://go.dev/). IT implements a basic HTTP handler and uses the a [fork](https://github.com/01100100/forms) of the [albrow/forms](https://github.com/albrow/forms) module. It accepts any incoming post request and validates it for a valid looking email and name.

It also check's that a hidden `spam` field is set to `safe` to avoid bot spam.

It then writes valid form data to a local json file.

## Deployment

It is deployed [via Dockerfile](https://fly.io/docs/languages-and-frameworks/dockerfile) using to the [fly.io](https://fly.io) platform.

The [Dockerfile] defines who to build the app and build a container image.

The file `fly.io` describes how to deploy this image and expose it to the outer world.

### CI

This repo contains a github action to deploy changes on the master branch using the fly.io cli.

### CLI

```bash
# Install flyctl on GNU/Linux
curl -L https://fly.io/install.sh | sh
```

The CLI can be used to create secrets, deploy the app, and ssh into the server to execute commands.

```bash
# This will verify the config and send it to a build and deploy pipeline.
fly deploy
# If the deployment is successful, the defined Entrypoint in the Dockerfile will be executed, creating a starting the server.
```

```bash
➜  goatform git:(main) ✗ flyctl ssh console

Connecting to xxxx:x:xxxx:xxx:xx:xxxx:xxxx:x... complete
/ # cat app/db.json
{"FormData":{"Values":{"email":["example@email.com"],"name":["Foo Bar"],"spam":["safe"],"submit":[""]},"Files":{}}}
...
/ #
```

Build and deploy process logs are available through the [fly.io dashboard](https://fly.io/apps/).

## Future

- [ ] Mount a storage volume. Right now logs are ephemeral.
- [ ] Write a super small front end with password auth.
