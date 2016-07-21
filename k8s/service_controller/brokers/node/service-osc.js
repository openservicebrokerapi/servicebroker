var http = require('http');
var url = require('url');
var mysql = require('mysql');

var port
var host
var user
var password

// HTTP options for fetching from service controller
var serviceControllerOptions = {
  host: 'service-controller',
  port: 10000,
  path: '/v2/service_brokers/test/service_instances/guestbook-mysql/service_bindings/guestbook-mysql-binding'
};

// mysql connection, created lazily upon request based on the above params
var connection

http.get(serviceControllerOptions, function(resp){
  var body = '';
  resp.on('data', function(chunk){
    body += chunk;
  });

  resp.on('end', function() {
    var cred = JSON.parse(body);
    
    host = cred['hostname'];
    port = cred['port'];
    user = cred['username'];
    password = cred['password'];
  });
}).on("error", function(e){
  console.log("Got error: " + e.message);
});

http.createServer(function (req, res) {
  if (connection == null) {
    console.log("Creating mysql")
    connection = mysql.createConnection({
      port     : port,
      host     : host,
      user     : user,
      password : password
    });
  }

  reqUrl = url.parse(req.url, true);

  res.useChunkedEncodingByDefault = false;
  res.writeHead(200, {'Content-Type': 'text/html'});

  if (reqUrl.pathname == '/_ah/health') {
    res.end('ok');
  } else if (reqUrl.pathname == '/exit') {
    process.exit(-1)
  } else {
      connection.query('use demo;', function(err) {
          if (err) {
              res.end('error: ' + err + '\n');
              connection.end();
          } else {
              if (reqUrl.query && reqUrl.query['msg']) {
                  var msg = reqUrl.query['msg'];
                  connection.query('insert into log(message) values (' + connection.escape(msg) + ');', function(err) {
                      if (err) {
                          res.end('error: ' + err + '\n');
                      } else {
                          res.end('added');
                      }
                 });
              } else {
                  connection.query('select * from log;', function(err, rows, fields) {
                      if (err) {
                      res.end('error: ' + err + '\n');
                      } else {
                          var result = '';
                          for (var i = 0; i < rows.length; ++i) {
                              result += rows[i].message + '<br>';
                          }
                          res.end(result);
                      }
                  });
              }
          }
      })
  }

}).listen(8080, '0.0.0.0');

console.log('Server running at http://127.0.0.1:8080/');
