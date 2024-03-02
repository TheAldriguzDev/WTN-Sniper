# WTN Sniper
With this free WTN (also known as we the new) bot you will be able to sell your shoes easily by sniping the offers that wtn sends to you.

## Why you need the sniper?
If you are already a reseller you know that wtn offers gets accepted from other in few seconds, with this bot you will be able to do the same, but for free!

## What is included?
I made different tools to help your selling experience on WTN:
1. Label download
2. Listing exporter
3. Listing extender
4. Offers monitor and sniper

This is what you need to start selling your shoes on WTN and make profit.

## What do i need to have to run the bot?
You only need to have some proxies, if you don't know where to find them you can send me a dm on discord. 
You don't need any third captcha solver provider since the bot is able to bypass the recaptcha challenge located in the login page.

You can also run without proxies but is not recommended, to do this create in the proxies folder an empty txt file.
To use proxies, just input them in the txt file.

## How to run the bot
First of all you need to download the bot by using this command on your terminal:
```shell
git clone https://github.com/TheAldriguzDev/WTN-Sniper
```

Then download the packes used by using this command:
```shell
go mod init
```

Now you are ready to start for the first time the bot, you have to option here:

Run from the terminal:
```shell
go run .
```

Or build the exe:
```shell
go build
```

In the first run the bot will create some folder, so i suggest to put all the downloaded files inside another main folder.
The created folders will be:
- Proxies
- Data
- Settings

After this first run the bot will close and you need to fill some data:
> Go on /Settings and fill the settings.json file. 
```js
{
    "success_webhook": "DISCORD WEBHOOK",
    "error_webhook": "DISCORD WEBHOOK",
    "wtn_account": {
        "email": "YOUR EMAIL",
        "password": "YOUR PASSWORD"
    },
    "delay": 2
}
```

After this step you will be able to use the WTN bot functions related to:
- Listings
- Label

To use the sniper you need to do an extra step, start the bot, you need to export the listings, so on the first menu select 1, than select your proxy list, than select again 1 to export the listings.
The bot will save a csv file in /Settings, fill with your min price.

After this step you are able to completely use the bot. Remember that it will accept offer that the price is grader than the min you have just set.

## Disclaimer
This project has been created exclusively for educational and demonstrative purposes. Please note the following:

- **Purpose:** This repository is intended to provide code examples for educational purposes related to the study of how requests works.

- **Limitations:** The author disclaims any responsibility for the misuse or harmful use of this code by third parties. The code is usable only for personal use, commercial use is not allowed.

- **No Warranty:** No warranty is provided regarding the accuracy, completeness, or reliability of the content. Users are responsible for their use of the code in this repository.

The use of this project is at the user's own risk and discretion.

## Credits
Special thanks to bogdanfinn for his [TLS client](https://github.com/bogdanfinn/tls-client)

## Contact
If you need any help setting up the bot or any question feel free to dm me on discord: **aldriguz**

## Conclusion
If you liked this repo or you found anything useful a star would be appreciated! ⭐