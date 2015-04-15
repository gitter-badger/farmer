'use strict';

var io = require('socket.io-client');


function Client (stationServer, workID) {
    this.serverUrl = stationServer + '/' + workID;
    this.stage = 0;
}

Client.prototype.listen = function () {
    var self = this,
        socket = io.connect(self.serverUrl);

    console.log('>>>>>>>>>> listen');
    socket.on('connect', function() {
        console.log('>>>>>>>>>> listener connected');
        socket.on('event', function (data) {
            console.log('>>>>>>>>>> client receive data', data);
        });
    });


    socket.on('error', function (error) {
        console.log('error >', error);
    });
};

module.exports = Client;