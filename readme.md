# slogth

Ingest logs from stdin and re-emit them with specified delay.

## TL;DR

```shell
ping 1.1.1.1 | go run github.com/flashbots/slogth/cmd --delay 5s
```

```shell
# nothing happens for 5 seconds

PING 1.1.1.1 (1.1.1.1): 56 data bytes
64 bytes from 1.1.1.1: icmp_seq=0 ttl=57 time=32.841 ms
64 bytes from 1.1.1.1: icmp_seq=1 ttl=57 time=16.679 ms
64 bytes from 1.1.1.1: icmp_seq=2 ttl=57 time=22.551 ms
64 bytes from 1.1.1.1: icmp_seq=3 ttl=57 time=24.806 ms
64 bytes from 1.1.1.1: icmp_seq=4 ttl=57 time=34.152 ms
64 bytes from 1.1.1.1: icmp_seq=5 ttl=57 time=25.536 ms
64 bytes from 1.1.1.1: icmp_seq=6 ttl=57 time=27.846 ms
64 bytes from 1.1.1.1: icmp_seq=7 ttl=57 time=23.651 ms
^C64 bytes from 1.1.1.1: icmp_seq=8 ttl=57 time=41.213 ms  # <- Ctrl-C
64 bytes from 1.1.1.1: icmp_seq=9 ttl=57 time=20.200 ms
64 bytes from 1.1.1.1: icmp_seq=10 ttl=57 time=46.447 ms
64 bytes from 1.1.1.1: icmp_seq=11 ttl=57 time=27.891 ms
64 bytes from 1.1.1.1: icmp_seq=12 ttl=57 time=26.006 ms   # remaining 5 seconds of logs
```
