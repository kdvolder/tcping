TCPing
======

A `tcping` cli is a tool to measure the latency of a 'req/resp' on a tcp socket.

The idea is that `tcping` provides both a client and server. You run the server
in one place and then run a client in another pointing it at the server.

The client works in a way similar to the standard `ping` utility. But it is focussed
on measuring the latency of request response interactions over an already open socket.

So the client will:

- open a connection to the server
- send small 'request' to the server
- measuer the time it takes to receive a response.
- print out basic statistics about the response times.