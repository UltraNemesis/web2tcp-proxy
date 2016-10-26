
(function () {
    
    Base64 = {
        // Copyright (c) 2007-2016 Kevin van Zonneveld (http://kvz.io) and Contributors (http://locutus.io/authors)
        // Based on 
        //          https://github.com/kvz/locutus/blob/master/src/php/url/base64_encode.js
        //          https://github.com/kvz/locutus/blob/master/src/php/url/base64_decode.js
        CODES: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",
        
        encode: function (data) {
            var o1, o2, o3, h1, h2, h3, h4, bits;
            var res = [];

            if (!data) {
                return "";
            }

            var i = 0;
            do { // pack three octets into four hexets
                o1 = data[i++];
                o2 = data[i++];
                o3 = data[i++];

                bits = o1 << 16 | o2 << 8 | o3;

                h1 = bits >> 18 & 0x3f;
                h2 = bits >> 12 & 0x3f;
                h3 = bits >> 6 & 0x3f;
                h4 = bits & 0x3f;

                // use hexets to index into CODES, and append result to encoded string
                res.push(this.CODES.charAt(h1) + this.CODES.charAt(h2) + this.CODES.charAt(h3) + this.CODES.charAt(h4));

            } while (i < data.length);

            var enc = res.join("");

            var r = data.length % 3;

            return (r ? enc.slice(0, r - 3) : enc) + "===".slice(r || 3);
        },
        
        decode: function (data) {

            var o1, o2, o3, h1, h2, h3, h4, bits;

            var res = [];

            if (!data) {
                return [];
            }

            data += "";

            var i = 0;
            do {
                // unpack four hexets into three octets using index points in CODES
                h1 = this.CODES.indexOf(data.charAt(i++));
                h2 = this.CODES.indexOf(data.charAt(i++));
                h3 = this.CODES.indexOf(data.charAt(i++));
                h4 = this.CODES.indexOf(data.charAt(i++));

                bits = h1 << 18 | h2 << 12 | h3 << 6 | h4;

                o1 = bits >> 16 & 0xff;
                o2 = bits >> 8 & 0xff;
                o3 = bits & 0xff;

                if (h3 == 64) {
                    res.push(o1);
                } else if (h4 == 64) {
                    res.push(o1);
                    res.push(o2);
                } else {
                    res.push(o1);
                    res.push(o2);
                    res.push(o3);
                }
            } while (i < data.length);

            return res;
        }
    };

    Web2TcpSocket = function (url) {
        var sock = null;
        var _self = this;

        if (window.SockJS !== undefined) {
            sock = new SockJS(url)
        }
        else if (window.WebSocket !== undefined) {
            sock = new WebSocket(url.replace("http", "ws") + "/websocket");
        }
        else {
            console.error("Neither SockJS or WebSocket support found.")
        }

        sock.onopen = function () {
        }

        sock.onmessage = function (e) {
            if (e.data == "STATUS:CONNECTED") {
                _self.onopen()
            }
            else {
                var bData = {}
                bData.data = Base64.decode(e.data)
                _self.onmessage(bData);
            }

        }

        sock.onerror = function (e) {
            _self.onerror(e);
        }

        sock.onclose = function () {
            _self.onclose();
        }

        this.send = function (data) {
            sock.send(Base64.encode(data));
        }

        this.close = function () {
            sock.close();
        }
    }

    Web2TcpSocket.prototype.onopen = function () { }
    Web2TcpSocket.prototype.onmessage = function (e) { }
    Web2TcpSocket.prototype.onerror = function (e) { }
    Web2TcpSocket.prototype.onclose = function () { }
})();



