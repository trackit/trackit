var express = require('express');
var app = express();

app.use(express.static(__dirname + '/build')); // set the static files location /public/img will be /img for users

// This route deals enables HTML5Mode by forwarding missing files to the index.html
app.all('/*', function(req, res) {
  res.sendFile(__dirname +'/build/index.html');
});

// listen (start app with node server.js) ======================================
var port = process.env.UI_PORT || 80;
app.listen(port);
console.log("[App listening on port " + port + ']');
