# bittorrent-client

The project is very limited in scope for now, although I am still working on it.
A well seeded torrent can be downloaded as fast as any standard bittorrent client.

Current status:
- Only download supported (leeching)
- Supports only UDP tracker
- No NAT traversal
- Can't be used with magnet link, and also no support for DHT.
- Tries to reestablish dropped connection with peers, but currently it does't look for new peers.
- Stores the entire file in memory, until download is finished. Hence no COD.


Currently working on:
- Writing pieces to disc as soon as it is recieved to optimise memory usage.
- Capping download speed.
- Resend Announce after some fixed interval, and spin new goroutines for any new peer.
- Adding support for HTTP trackers.


Still exploring:
- Upload
- NAT traversal
- DHT
- uTP
- different piece selection algorithms
- supporting magnet link, using external packages
