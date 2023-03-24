# 🛡️OTGuard

OTGuard is a tool for providing two-factor authentication with zero configuration on services that have single-factor authentication and do not integrate with any 2FA solution. It is not meant to be used in isolation but is intended to be part of a defense in depth strategy where applications listening on the ports have their own authentication. OTGuard aims to provide a small attack surface with only ~ 200 lines of code and uses privilege separation.

One of the key benefits of OTGuard is that it is the only solution that works with applications that expect direct access to services, such as Jellyfin and Nextcloud. While it is not a replacement for Identity Aware Proxies or SSO, it aims to provide 2FA with zero configuration on services.

OTGuard is the second factor in 2FA, and it is expected that services provide the first factor with usernames and passwords. The tool uses OTP codes, which makes it less vulnerable to replay attacks. All firewall rules are reset daily to provide some security against compromised devices. Users will need to authenticate every day, and this reduces the amount of code and the attack surface even further.

The workflow for OTGuard is simple. Users try to access a service, and if it does not work, they log in on OTGuard. Once authenticated, users will have access to the desired service. The tool does not check whether the user is already allowed or use cookies as it seems like unnecessary complexity.

OTGuard is inspired by solutions like port knocking or Fwknop but unlike port knocking, passwords are not sent in the clear, and unlike Fwknop, it provides OTP codes so that passwords cannot be reused.
