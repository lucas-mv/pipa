# pipa
Pipa is a simple [go](https://golang.org/) application that can be run by the command line to fetch the most relevant twitter trending topics based on your address.

As of right now you need to fill the settings on `settings.json` providing a [Twitter Developer App ID](https://developer.twitter.com/en/docs/basics/getting-started) and get a [Basic Auth Base64 key](https://developer.twitter.com/en/docs/basics/authentication/basic-auth), a [Bing Maps API key](https://www.bingmapsportal.com/) and your address (really simple stuff, as if you were to look up your house on Google or Bing maps).

## Build

`git clone` the project and run `go build` on the cloned folder.

## Running

After you built the application, you'll get a `pipa` executable file, all you need to do to run it successfully is add a configured `settings.json` file to the folder where you want to run the `pipa` executable.

## Why pipa?

Pipa is the name of a dog that I really like that lives in my mom's house.
