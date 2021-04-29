DNSCODE
-------

Control DNS records in single file independent of provider.

### Supported providers:

 * Pdd yandex
 * Reg.ru
 * Adman.com

### Usage:

Create file dnscode.json for example
```json
{
 "zones": [
  {
   "provider": "yandex",
   "name": "example.ru",
   "records": null,
   "connection": {
    "PddToken": "token from page https://pddimp.yandex.ru/api2/admin/get_token"
   }
  },
  {
   "provider": "adman",
   "name": "example.org",
   "records": null,
   "connection": {
    "login": "login to adman",
    "mdpass": "api key generated on page https://adman.com/pages/api/"
   }
  },
  {
   "provider": "regru",
   "name": "example.com",
   "records": null,
   "connection": {
    "password": "password",
    "username": "login"
   }
  }
 ]
}
```

Then you can import records:
```shell
dnscode import
```

This command import records. Allowed flags:

 * --filename - filename which will contain imported record
 * --useImport=true|false - if true dnscode will create own file for each zone (default true)

Now you can manage zone by edit file. After you can see changes:

```shell
dnscode plan
```
And apply changes
```shell
dnscode apply
```

Allowed flags for apply command:
 * --force=true|false - delete records from provider (default true)
 * --interactive=true|false - confirm apply? (default true)

Flags for all commands:
 * --proxy=addr.com - proxy server. For example see proxy/server.go

### Add provider
To create a new provider, add provider file to providers folder, 
create provider struct which implement providers.BaseProvider.
Add your provider to function GetProvider() in file providers/base.go