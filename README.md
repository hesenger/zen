# Zen

Zen is a zero setup in-place deployment tool that allows you deploy multiple applications
to your server. It was created from the need to have cheap and easy way to deploy and keep
small apps without huge infrastructure like kubernetes and mainly, without the heavy cost of maintaining multiple servers just to orchestrate a few instances.

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/hesenger/zen/main/scripts/install.sh | sudo bash
```

Once installed access your server under port 8888 from any browser to setup credentials
and configure your applications.

Zen will monitor published releases for the applications, download and execute them
automatically. You can use private repositories on Github through Token authentication.
