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
- Requesting multiple pieces from the same peer concurrently to allow faster download.
- Downloading multifile torrent
- Keep looking for new peers, after some intervals.
- Reestablishing dropped connection with peers.
- Adding support for WS and HTTP trackers.
- Actually writing the downloaded buffer to a file. (LOL)


Still exploring:
- Controlling max download speed
- Upload
- NAT traversal
- DHT
- uTP
- different piece selection algorithms
- supporting magnet link, using external packages
