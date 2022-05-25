# chrony monitoring with Netdata

[`chrony`](https://chrony.tuxfamily.org/)  is a versatile implementation of the Network Time Protocol (NTP).

The modules will monitor local host `chrony` server.

This module use golang to collect chrony and produces:
* stratum
* frequency
* last offset
* RMS offset
* residual freq
* root delay
* root dispersion
* skew
* leap status
* update interval
* current correction
* current source server address