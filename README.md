# gash

[![Build](https://img.shields.io/badge/endpoint.svg?url=https://badger.seankhliao.com/r/github_seankhliao_gash)](https://console.cloud.google.com/cloud-build/builds?project=com-seankhliao&query=source.repo_source.repo_name%20%3D%20%22github_seankhliao_gash%22)
[![License](https://img.shields.io/github/license/seankhliao/gash.svg?style=for-the-badge)](LICENSE)

Google Analytics, Self Hosted

## What

Proxy Google Analytics through your own domain

Note: must not allow cloudflare to cache query strings

## Use

same as standard gtag script but replace the url

```js
    <script async src="https://gash.seankhliao.com/js"></script>
    <script>
      window.dataLayer = window.dataLayer || [];
      function gtag() {
        dataLayer.push(arguments);
      }
      gtag("js", new Date());

      gtag("config", "UA-XXXXXXXXX-Y");
```

## What it Does

1. Get the gtag sccript from: `https://https://www.googletagmanager.com/gtag/js`
   a. Replace `https://www.google-analytics.com/` with `https://gash.seankhliao.com`
   b. Cache the script
2. Get the analytics.js script from `https://www.google-analytics.com/analytics.js`
   a. Replace `https://www.google-analytics.com/` with `https://gash.seankhliao.com`
   b. Cache the script
3. Collect reports at `https://gash.seankhliao.com/collect/`
   a. add `uip=<true client ip address`
   b. forward to `https://www.google-analytics.com/collect`
4. Periodically update the scripts

## TODO

- [ ] verify client ip is valid
