---
title: Installation
nav_order: 1
layout: default
permalink: installation/
description: "Set-up and Installation"
---


# Installation
{: .no_toc}


Burnout Barometer is easy-to-configure and deployable as a serverless application.
{: .fs-6 .fw-300 }

Assuming that you already know [how to create a Slack
App](https://api.slack.com/start), this page will walk you through on how to set-up
your Barometer, then show various deployment options at your disposal.

---


## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

## Initial Setup 

1. **Download the executable**. Ensure that you have downloaded the `barometer`
   executable. Follow the [download instructions]({{ site.baseurl }}/download)
   for more info.
2. **Initialize configuration**. Run `barometer init`. Answer a series of
   prompts to configure your Barometer. The following config options need to be
   set:

    | Option         | Docker Env Var | Description                                                                                                                                                                                                                                                          |
    |----------------|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
    | Table          | BB_TABLE       | The database connection URL to store Barometer logs. For Bigquery, use the `bq` protocol like so: `bq://my-gcp-project.my-dataset.my-table`                                                                                                                          |
    | Slack Token    | BB_SLACK_TOKEN | The Slack Token generated whenever you create an App. This is used to verify that the incoming request came from the authorized account. See this [page](https://slack.com/intl/en-ph/help/articles/215770388-Create-and-regenerate-API-tokens) for more information |
    | Area           | BB_AREA        | The [IANA compliant area](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) for correcting the timezone. For example, `Asia/Manila`. |
    |Twitter Consumer Key | BB_TWITTER_CONSUMER_KEY | *(Optional)* Your Twitter Consumer Key to fetch Tweets from [tinycarebot](https://twitter.com/tinycarebot). Check [this link](https://dev.twitter.com/apps/new) for details |
    | Twitter Consumer Secret| BB_TWITTER_CONSUMER_SECRET        | *(Optional)* Your Twitter Consumer Secret to fetch Tweets from [tinycarebot](https://twitter.com/tinycarebot). Check [this link](https://dev.twitter.com/apps/new) for details |
    | Twitter Access Key| BB_TWITTER_ACCESS_KEY        | *(Optional)* Your Twitter Access Key to fetch Tweets from [tinycarebot](https://twitter.com/tinycarebot). Check [this link](https://dev.twitter.com/apps/new) for details |
    | Twitter Access Secret| BB_TWITTER_ACCESS_SECRET        | *(Optional)* Your Twitter Access Secret to fetch Tweets from [tinycarebot](https://twitter.com/tinycarebot). Check [this link](https://dev.twitter.com/apps/new) for details |

    You can find more information about the `init` command by running
    `barometer init --help`.

3. **Check if a config file has been generated**. After running the `init`
   command, you should see a `config.json` file with your configuration. We
   will use this later on when deploying or starting the server.


## Deployment Options

Burnout Barometer is a server-side application, and can be deployed by various
means. You can also check out our Docker image located in the Azure
Container Registry.

### Deploy via Google Cloud Functions

To deploy via [Google Cloud Functions](https://cloud.google.com/functions/),
clone the Github repository to access the `function/` directory:

```bash
git clone git@github.com:ljvmiranda921/burnout-barometer.git
```

Copy over the `config.json` that you've generated before into this path:

```bash
cp path/of/config.json function/
```

Lastly, head to the `function/` directory and execute the Cloud Function deploy
command:

```bash
cd function/
gcloud functions deploy BurnoutBarometerFn --runtime go111 --triger-http
```

Once successful, Cloud Functions will provide you a URL that you can now
add in your Slack Application's Slash command.

### Deploy via Google Cloud Run

You can deploy to [Google Cloud Run](https://cloud.google.com/run/) using the
`ljvmiranda.azurecr.io/burnout-barometer` Docker image. You need to set
some [environment variables]({{ site.baseurl  }}/installation.html#initial-setup) to configure the barometer: 

To deploy, run the following command:

```bash
gcloud beta run deploy burnout-barometer \
    --image ljvmiranda.azurecr.io/burnout-barometer \
    --set-env-vars=BB_PROJECT_ID=<PROJECT_ID>,BB_TABLE=<TABLE>,BB_SLACK_TOKEN=<TOKEN>,BB_AREA=<AREA>
```

Or better yet, just click the button below:

[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run?git_repo=https://github.com/ljvmiranda921/burnout-barometer.git)

---

Now that you have configured and deployed your barometer, check-out the [Usage
Section]({{ site.baseurl }}/usage) to learn more about this tool!
{: .info }
