# Sagittarius client :sagittarius:
A [sagittarius](https://github.com/nyugoh/sagittarius) client to monitor, collect and send log files to a central repository.

## Todo

- [ ] Send ping to server for constant monitoring
- [ ] Read config to know which folder to monitor
- [ ] Send logs to ``server``

- have a route for sending changes to server
- have a go routine to watch a certain folder, if file changes, send an update to server
- To start, request jwt from server by sending a post request with app name, hash and port and ip, in return you get
  a jwt which is included in every other request.
  
## Deployment
- Download the file [sagittarius-client](sagittarius-client) , [service](sagittarius-client.service) and [.env](.env)
- Create a directory to host the app, also create a log folder
- Move `sagittarius-client` & `.env` to the app folder
- Edit `.service` to reflect the new app folder
- Edit `.env`, `SAG_SERVER` with the main server and also the `CLIENT_NAME` which will be used to register the app on server.
- Move `sagittarius-client` to `/etc/systemd/system`
- Enable and start app `systemctl start sagittarius-client` and `systemctl enable sagittarius-client`