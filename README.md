# DanWolf.NET

This repo contains source code for my personal website which I try to add to in tiny regular increments. While I don't intend to focus upon only features that will be used by a wider audience, I will respond to issues, assist in understanding code, and (once I have something functional) consider splitting off smaller packages that are meant to be reused as libraries or separate applications--such as the blog itself.

The contents of this document should be taken with only light consideration as I am in the process of reworking the 

## Install

```
go get github.com/palumacil/dwn
```

## Run

### Development

I use VS Code with the popular [lukehoban.go](https://github.com/Microsoft/vscode-go). Once I hit F5 to run the Go application, I type `npm start` in the terminal to run the Angular application.

### Production

The deployment version of this project will be built in Docker. I need to complete the Dockerfile by adding a volume for the existing database, copying the Angular files to the image, and I will also need to add the docker run command to the build.ps1 file.

## Issues

 - This project is not yet fully functional after login.
 - Replace the Bootstrap login modal with something more Angular-friendly.
 - The first-run code which prompts to create a user if no admin user exists is clunky and needs to be replaced with a setup page that is accessible when the application is run with a `-setup` flag.

## License

This project bears an MIT license, allowing permissive use. If you do something clever, please consider mentioning it to me so that I am able to consider implementing it here or in the future version of this project which will be more meant for public consumption.