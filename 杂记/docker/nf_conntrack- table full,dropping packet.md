# nf_conntrack: table full,dropping packet

---

This is normal.
Docker uses iptables connection tracking.
Connection tracking entries will subside for a while, even after the connection is closed.
There are two ways to work around this:

1. increase nf_conntrack_max (it just uses more memory, but AFAIR it's about 100 bytes per connection, so even if you set it to 2000000, you end up using 200 MB of RAM)
2. reduce the nf_conntrack_tcp_timeout values so that connections are removed faster after they are closed.