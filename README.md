# Eventsapibeat

Eventsapibeat is the open source beats shipper for pulling events from the [1Password Events API](#).  
This beat will fetch successful and failed sign-in attempts and items usage data from public 1Password Events API.

## Installation
Build and install the binary for your system.  

```shell
make build_all_apps
```

Resulting packages will be located at _bin_ folder.  

## Configuration

Rename the sample configuration file _eventsapibeat-sample.yml_ to _eventsapibeat.yml_.  

Configure the `api_host` and `auth_token` after you create a 1Password Events API integration for your account in https://1Password.com.
```yaml
  api_host: "https://events.1password.com"
  #api_host: "https://events.ent.1password.ca"
  #api_host: "https://events.1password.eu"
  #api_host: "https://events.1password.eu"
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