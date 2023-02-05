# bittorrent-client

The project is very limited in scope for now, although I am still working on it.

Current status:
- Only download supported (leeching)
- Can't download multifile torrent
- Supports only UDP tracker
- Doesn't attempt to reestablish dropped TCP connection with peers.
- No NAT traversal
- Can't be used with magnet link, and also no support for DHT.


Currently working on:
- Downloading multifile torrent
- Keep looking for new peers, after some intervals.
- Reestablishing dropped connection with peers.
- Adding support for WS and HTTP trackers.


Still exploring:
- Upload
- NAT traversal
- DHT
- uTP
- different piece selection algorithms
- supporting magnet link, using external packages
