# bittorrent-client

The project is very limited in scope for now, although I am still working on it.

Current status:
- Only download supported (leeching)
- Can't download multifile torrent
- Supports only UDP tracker
- No NAT traversal
- Can't be used with magnet link, and also no support for DHT.
- Tries to reestablish dropped connection with peers, but currently it does't look for new peers.


Currently working on:
- Requesting multiple pieces from the same peer concurrently to allow faster download.
- Downloading multifile torrent
- Keep looking for new peers, after some intervals.
- Adding support for WS and HTTP trackers.
- Adding a dynamic timeout for peers, based on the size of a piece and health of the network.


Still exploring:
- Controlling max download speed
- Upload
- NAT traversal
- DHT
- uTP
- different piece selection algorithms
- supporting magnet link, using external packages
