# bittorrent-client

- git clone https://github.com/stillme1/bittorrent-client
- go get .
- go run . "path to torrent" "path to download destination"     
Use different files for stdout and stderr for better logging    
Eg.   
go run . torrent/0.torrent download > output.txt 2> error.txt

The project is very limited in scope for now, although I am still working on it.
A well seeded torrent can be downloaded as fast as any standard bittorrent client.

Current status:
- Only download supported (leeching)
- Supports only UDP tracker
- No NAT traversal
- Can't be used with magnet link, and also no support for DHT.


Currently working on:
- Capping download speed.
- Upload
- NAT traversal
- DHT
- uTP
- different piece selection algorithms
- supporting magnet link, using external packages
