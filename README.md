# DanWolf.NET

## Install

```
git clone github.com/PaluMacil/dwn
go build
```

## Configuration

The settings table below shows the environmental variables to override 
defaults. The default sometimes varies between prod and dev. In prod mode, 
you must set and encryption key and cannot use the default dev key.

### Core Settings

| setting               | env                | prod                 | dev                                |
|-----------------------|--------------------|----------------------|------------------------------------|
| host                  | DWN_HOST           | danwolf.net          | localhost                          |
| port (api/server)     | DWN_PORT           | 3035                 | 3035                               |
| ui port               | DWN_UI_PROXY_PORT  | 443                  | 4200                               |
| content root (ui)     | DWN_CONTENT_ROOT   | /opt/danwolf.net/ui  | /home/dan/repos/dwn-ui/dist        |
| initial admin         | DWN_INITIAL_ADMIN  | (none configured)    | dcwolf@gmail.com                   |
| data directory        | DWN_DATA_DIR       | data                 | data                               |
| master encryption key | DWN_MASTER_ENC_KEY | (must configure!)    | 3d17618d4297f83665b32e28f9b1c23d   |

Note that if you are developing on **Windows**, encyption is disabled (along with memory-mapped files since Windows uses 
disk IO instead of memory mapped files). You will also need to set your own content root if you are using Windows or 
your username isn't *dan*.

The initial admin is used for recognizing a user that should be marked as verified and added to the admin usergroup upon 
user registration. This user will not have a default password, so setting up an auth provider foreign system such as 
Google will also be required. Once this user logs in, this environmental variable should be removed.

The **ui port** is used for building URLs for things like OAuth2 redirects. This is the port which will be used to 
access the website. The **port** setting, on the other hand, represents the port the server application is listening to. 
On a server without a proxy in front of its web applications, these values might be the same. At this time, however, the 
application depends upon a server like Caddy or Nginx to serve TLS. While running in dev mode, you will probably be 
running `ng serve` to serve the UI from port 4200, and the angular application will proxy backend requests to port 3035.

### Foreign Systems

Foreign systems include authentication providers, SMTP providers, and other APIs. For foreign systems that can be configured 
by environmental variable, doing so will lock those settings in the UI, preventing live modification of settings. Most 
foreign systems require a key and secret or at least a token. These details should be kept in 
a password manager or other type of secure system.

#### Authentication Providers

| setting              | env                     |
|----------------------|-------------------------|
| Google OAuth2 Key    | DWN_OAUTH_GOOGLE_KEY    |
| Google OAuth2 Secret | DWN_OAUTH_GOOGLE_SECRET |

##### External Provider Configuration

- **Google:** Go to [Developer Console > Credentials](https://console.developers.google.com/apis/credentials) and click 
the oauth client.

#### Other Foreign Systems

| setting              | env                      |
|----------------------|--------------------------|
| Sendgrid SMTP Key    | DWN_SMTP_SENDGRID_KEY    |
| Sendgrid SMTP Secret | DWN_SMTP_SENDGRID_SECRET |

##### External System Configuration

- **Sendgrid:** [Settings > API Keys](https://app.sendgrid.com/settings/api_keys)

## Run

For development defaults, leave off the prod subcommand:

```
.\dwn prod
```

To run the UI (dev mode only--prod mode serves its own files), use the command `ng serve` which depends upon 
[Node](https://nodejs.org/) and the [Angular CLI](https://angular.io/).

## Licensing

The license for this project is MIT but some dependencies might have Apache or similar licenses. General reuse isn't necessarily intended and APIs might break, but it is permissible.

## Roadmap

### Section Priorities
- game
- file
- visit
- (then comment docs, tests, and README)
- download
- photo
- blog
- backup
- landlord

### Unsorted Thoughts 
- Friends linked to friend role if appearing in Facebook friends
- Movie list: people in friends role should be able to vote for movies they want to see (to plan movie parties)
- Admins should be able to do some things on behalf of others (mark that someone wants to see a movie even if they don't have an account)
- photo album
- blog
- downloads

Improve:
- get from env config: `TokenName:   "dwn-token"`
- Make sure certain keys (e.g. email) are not case sensitive
- document use of environmental variables
- when routing, first trim off base url (if base url isn't '/', then currently routing will be misaligned when split and break)
