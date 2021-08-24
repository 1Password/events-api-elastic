Eventsapibeat
=============

Eventsapibeat is the open source libbeat based data shipper for pulling events from the 1Password Events API.  
This beat will fetch successful and failed sign-in attempts and items usage data from public 1Password Events API.

Installation
------------

Download the latest binaries from [the releases page](https://github.com/1Password/events-api-elastic/releases/latest).  
Or build from sources, *resulting binary will be located at 'bin' folder*:

```shell
make eventsapibeat
```

Configuration
-------------

Rename the sample configuration file *eventsapibeat-sample.yml* to *eventsapibeat.yml*.

Create a [1Password Events Reporting](https://support.1password.com/events-reporting-elastic/) integration for your account and configure the `auth_token`.

```yaml
signin_attempts:
  auth_token: "token"
item_usages:
  auth_token: "token"
```

Configure the remaining options and set your output as usual.

Run
---

```
./eventsapibeat -c eventsapibeat.yml -e
```

Elastic Common Schema
---------------------

### Sign-in Attempts fields

| Field                                 | Description                                                                                                                                               |
|---------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `@timestamp`                          | The date and time of the sign-in attempt                                                                                                                  |
| `event.action`                        | The category of the sign-in attempt                                                                                                                       |
| `user.id`                             | The UUID of the user that attempted to sign in to the account                                                                                             |
| `user.full_name`                      | The name of the user, hydrated at the time the event was generated                                                                                        |
| `user.email`                          | The email address of the user, hydrated at the time the event was generated                                                                               |
| `os.name`                             | The name of the operating system of the user that attempted to sign in to the account                                                                     |
| `os.version`                          | The version of the operating system of the user that attempted to sign in to the account                                                                  |
| `source.ip`                             | The IP address that attempted to sign in to the account                                                                                                   |
| `onepassword.uuid`                    | The UUID of the event                                                                                                                                     |
| `onepassword.session_uuid`            | The UUID of the session that created the event                                                                                                            |
| `onepassword.type`                    | Details about the sign-in attempt                                                                                                                         |
| `onepassword.country`                 | The country code of the event. Uses the ISO 3166 standard                                                                                                 |
| `onepassword.details`                 | Additional information about the sign-in attempt, such as any firewall rules that prevent a user from signing in                                          |
| `onepassword.client.app_name`         | The name of the 1Password app that attempted to sign in to the account                                                                                    |
| `onepassword.client.app_version`      | The version number of the 1Password app                                                                                                                   |
| `onepassword.client.platform_name`    | The name of the platform running the 1Password app                                                                                                        |
| `onepassword.client.platform_version` | The version of the browser or computer where the 1Password app is installed, or the CPU of the machine where the 1Password command-line tool is installed |

### Item Usages fields

| Field                                 | Description                                                                                                                                               |
|---------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `@timestamp`                          | The date and time of the sign-in attempt                                                                                                                  |
| `user.id`                             | The UUID of the user that accessed the item                                                                                                               |
| `user.full_name`                      | The name of the user, hydrated at the time the event was generated                                                                                        |
| `user.email`                          | The email address of the user, hydrated at the time the event was generated                                                                               |
| `os.name`                             | The name of the operating system the item was accessed from                                                                                               |
| `os.version`                          | The version of the operating system the item was accessed from                                                                                            |
| `source.ip`                             | The IP address the item was accessed from                                                                                                                 |
| `onepassword.uuid`                    | The UUID of the event                                                                                                                                     |
| `onepassword.used_version`            | The version of the item that was accessed                                                                                                                 |
| `onepassword.vault_uuid`              | The UUID of the vault the item is in                                                                                                                      |
| `onepassword.item_uuid`               | The UUID of the item that was accessed                                                                                                                    |
| `onepassword.client.app_name`         | The name of the 1Password app the item was accessed from                                                                                                  |
| `onepassword.client.app_version`      | The version number of the 1Password app                                                                                                                   |
| `onepassword.client.platform_name`    | The name of the platform the item was accessed from                                                                                                       |
| `onepassword.client.platform_version` | The version of the browser or computer where the 1Password app is installed, or the CPU of the machine where the 1Password command-line tool is installed |
