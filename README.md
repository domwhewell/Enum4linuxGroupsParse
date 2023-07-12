# Enum4linuxGroupsParse

A golang script to parse Enum4linux groups output and lookup what groups a user is a member of.
This could be done with grep but this script does it recursively looking up each new group name as its found so nested domain admins can be found

## Example
You can recursively lookup a group name
```bash
$ enum4linux_parse --file enum4linux_groups --group_name 'Administrators'
This script will recursively lookup groups or users in an active directory environment from the enum4linux groups output.
- For example 'user1' is a member of 'IT Support' which is itself a member of 'Domain Admins' making 'user1' effectively a domain admin.
This script will colour code the output so 'Groups' are always highlighted as red to make them easier to make out.

[!] Members of group: Administrators
Administrator
[!] Enterprise Admins is a member of Administrators
[!] Members of group: Enterprise Admins
EA-Admin
[!] Domain Admins is a member of Administrators
[!] Members of group: Domain Admins
EA-Admin
X-Admin
X-Dom
[!] IT Support is a member of Domain Admins
[!] Members of group: IT Support
Standard-User
```

Or you can recursively lookup a username
```bash
$ enum4linux_parse --file enum4linux_groups --username 'Standard-User'
This script will recursively lookup groups or users in an active directory environment from the enum4linux groups output.
- For example 'user1' is a member of 'IT Support' which is itself a member of 'Domain Admins' making 'user1' effectively a domain admin.
This script will colour code the output so 'Groups' are always highlighted as red to make them easier to make out.

Standard-User is a member of Domain Users
Domain Users is a member of Users
Standard-User is a member of IT Support
IT Support is a member of Domain Admins
Domain Admins is a member of Administrators
```
