# Eventsapibeat

Eventsapibeat is the open source libbeat based data shipper for pulling events from the 1Password Events API.
This beat will fetch successful and failed sign-in attempts and items usage data from public 1Password Events API.

## Installation

Download the latest binaries from [the releases page](https://github.com/1Password/events-api-elastic/releases/latest).
Or build from sources, _resulting binary will be located at 'bin' folder_:

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
audit_events:
  auth_token: "token"
```

Configure the remaining options and set your output as usual.

## Run

```
./eventsapibeat -c eventsapibeat.yml -e
```

## Elastic Common Schema

### Sign-in Attempts fields

| Field                                 | Description                                                                                                                                               | Type      |
| ------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- |
| `@timestamp`                          | The date and time of the sign-in attempt                                                                                                                  | date      |
| `event.action`                        | The category of the sign-in attempt                                                                                                                       | keyword   |
| `user.id`                             | The UUID of the user that attempted to sign in to the account                                                                                             | keyword   |
| `user.full_name`                      | The name of the user, hydrated at the time the event was generated                                                                                        | keyword   |
| `user.email`                          | The email address of the user, hydrated at the time the event was generated                                                                               | keyword   |
| `os.name`                             | The name of the operating system of the user that attempted to sign in to the account                                                                     | keyword   |
| `os.version`                          | The version of the operating system of the user that attempted to sign in to the account                                                                  | keyword   |
| `source.ip`                           | The IP address that attempted to sign in to the account                                                                                                   | ip        |
| `geo.country_iso_code`                | The country code of the event. Uses the ISO 3166 standard                                                                                                 | keyword   |
| `geo.region_name`                     | The region name of the event                                                                                                                              | keyword   |
| `geo.city_name`                       | The city name of the event                                                                                                                                | keyword   |
| `geo.location`                        | The longitude and latitude of the event                                                                                                                   | geo_point |
| `onepassword.uuid`                    | The UUID of the event                                                                                                                                     | keyword   |
| `onepassword.session_uuid`            | The UUID of the session that created the event                                                                                                            | keyword   |
| `onepassword.type`                    | Details about the sign-in attempt                                                                                                                         | keyword   |
| `onepassword.country`                 | The country code of the event. Uses the ISO 3166 standard                                                                                                 | keyword   |
| `onepassword.details`                 | Additional information about the sign-in attempt, such as any firewall rules that prevent a user from signing in                                          | keyword   |
| `onepassword.client.app_name`         | The name of the 1Password app that attempted to sign in to the account                                                                                    | keyword   |
| `onepassword.client.app_version`      | The version number of the 1Password app                                                                                                                   | keyword   |
| `onepassword.client.platform_name`    | The name of the platform running the 1Password app                                                                                                        | keyword   |
| `onepassword.client.platform_version` | The version of the browser or computer where the 1Password app is installed, or the CPU of the machine where the 1Password command-line tool is installed | keyword   |

### Item Usages fields

| Field                                 | Description                                                                                                                                               | Type      |
| ------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- |
| `@timestamp`                          | The date and time of the item usage                                                                                                                       | date      |
| `event.action`                        | The action performed on the item                                                                                                                          | keyword   |
| `user.id`                             | The UUID of the user that accessed the item                                                                                                               | keyword   |
| `user.full_name`                      | The name of the user, hydrated at the time the event was generated                                                                                        | keyword   |
| `user.email`                          | The email address of the user, hydrated at the time the event was generated                                                                               | keyword   |
| `os.name`                             | The name of the operating system the item was accessed from                                                                                               | keyword   |
| `os.version`                          | The version of the operating system the item was accessed from                                                                                            | keyword   |
| `source.ip`                           | The IP address the item was accessed from                                                                                                                 | ip        |
| `geo.country_iso_code`                | The country code of the event. Uses the ISO 3166 standard                                                                                                 | keyword   |
| `geo.region_name`                     | The region name of the event                                                                                                                              | keyword   |
| `geo.city_name`                       | The city name of the event                                                                                                                                | keyword   |
| `geo.location`                        | The longitutde and latitude of the event                                                                                                                  | geo_point |
| `onepassword.uuid`                    | The UUID of the event                                                                                                                                     | keyword   |
| `onepassword.used_version`            | The version of the item that was accessed                                                                                                                 | long      |
| `onepassword.vault_uuid`              | The UUID of the vault the item is in                                                                                                                      | keyword   |
| `onepassword.item_uuid`               | The UUID of the item that was accessed                                                                                                                    | keyword   |
| `onepassword.client.app_name`         | The name of the 1Password app the item was accessed from                                                                                                  | keyword   |
| `onepassword.client.app_version`      | The version number of the 1Password app                                                                                                                   | keyword   |
| `onepassword.client.platform_name`    | The name of the platform the item was accessed from                                                                                                       | keyword   |
| `onepassword.client.platform_version` | The version of the browser or computer where the 1Password app is installed, or the CPU of the machine where the 1Password command-line tool is installed | keyword   |

### Audit Events fields

| Field                              | Description                                                        | Type    |
| ---------------------------------- | ------------------------------------------------------------------ | ------- |
| `@timestamp`                       | The date and time of the audit event. Uses the RFC 3339 standard.  | date    |
| `event.action`                     | Details about the action taken for the audit event.                | keyword |
| `user.id`                          | The UUID of the user that performed the audit event.               | keyword |
| `source.ip`                        | The IP address that performed the audit event.                     | ip      |
| `onepassword.uuid`                 | The UUID of the audit event.                                       | keyword |
| `onepassword.object_type`          | The target object type of the audit event.                         | keyword |
| `onepassword.object_uuid`          | The target object UUID of the audit event.                         | keyword |
| `onepassword.aux_id`               | Any auxiliary ID of the audit event.                               | long    |
| `onepassword.aux_uuid`             | Any auxiliary UUID of the audit event.                             | keyword |
| `onepassword.aux_info`             | Any auxiliary info of the audit event.                             | keyword |
| `onepassword.session.session_uuid` | The UUID of the user session that performed the audit event.       | keyword |
| `onepassword.session.device_uuid`  | The UUID of the device that performed the audit event.             | keyword |
| `onepassword.session.login_time`   | The login time of the user session that performed the audit event. | date    |
