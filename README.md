# Eventsapibeat

Eventsapibeat is the open source libbeat based data shipper for pulling events from the 1Password Events API.  
This beat will fetch successful and failed sign-in attempts and items usage data from public 1Password Events API.

## Installation

Download the latest binaries from [the releases page](https://github.com/1Password/events-api-elastic/releases/latest).  
Or build from sources (_resulting binary will be located at _bin_ folder_):  

```shell
make eventsapibeat
```

## Configuration

Rename the sample configuration file _eventsapibeat-sample.yml_ to _eventsapibeat.yml_.

Create a [1Password Events Reporting](https://support.1password.com/events-reporting-elastic/) integration for your account and configure the `auth_token`.  

```yaml
signin_attempts:
  auth_token: "token"
item_usages:
  auth_token: "token"
```

Configure the remaining options and set your output as usual.

## Run

```
./eventsapibeat -c eventsapibeat.yml -e
```
