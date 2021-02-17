Home Dynamic DNS
------------------------
This is a small go program to update A records in aws hosted zones after pinging for public IP address from
[https://www.ipify.org/](https://www.ipify.org/). Meant to be run on a cron or systemd timer from a host on the
wan network you want records updated to. It provides essentially my own program for [DDNS](https://en.wikipedia.org/wiki/Dynamic_DNS)
instead of using builtin providers on my router or signing up for a ddns service.

## But why not use a free DDNS service?
This is a personal project to help explore some go features I don't yet have much experience with such as go routines and
`WaitGroups`. Eventually conditional compiling (for interfacing with systemd journaling, but not on windows),
interfaces (for logging to systemd or stdout) and to stretch some op muscles by packaging it as a `.deb` package  with
systemd unit timer configuration bundled in as well.

## Goals
1. Simple binary that updates route53 A records with new ip address if required
2. systemd unit timer configuration and .deb package
3. log messages to systemd journal if available on the system
4. Log to stdout if not launched through systemd

## TODO
- [x] write initial a record update code
    - [x] Use one go routine per hosted zone
    - [x] log all to std out
- [ ] setup cross compiling (would be nice to be able to target arm64 too to try to run on a pi)
- [ ] add systemd unit timer configuartion  
- [ ] package as simple deb package [read](https://askubuntu.com/questions/1345/what-is-the-simplest-debian-packaging-guide)
- [ ] release script to upload deb to github releases
- [ ] update logging to write to systemd journal
   - [ ] conditionally choose the logger if systemd is detected
   - [ ] don't compile systemd files in a windows build
- [ ] Make less opinionated take a list of hostnames and hosted zones to change instead of changing all
- [ ] More verbose output
