# krapslog

[![Actions Status](https://github.com/acj/krapslog/workflows/CI/badge.svg)](https://github.com/acj/krapslog/actions)

Visualize a log file with [sparklines](https://en.wikipedia.org/wiki/Sparkline)

When troubleshooting a problem with a production service, I often need to get the general shape of a log file. Are there any spikes? Was the load higher during the incident than it was beforehand? Does anything else stand out? Without tooling to help you, a large log file is little more than a blob of data. This tool is designed to quickly surface key features of the log — and then get out of your way.

## Installing

```
$ go get -u github.com/acj/krapslog
```

## Usage

```
$ krapslog -h
Usage of krapslog:
  -concurrency int
        number of log lines to process concurrently (default 16)
  -format string
        date format to look for (see https://golang.org/pkg/time/#Time.Format) (default "02/Jan/2006:15:04:05.000")
  -markers int
        number of time markers to display
  -progress
        display progress while scanning the log file
```

## Examples

Get the basic shape:

```
$ krapslog /var/log/haproxy.log 
▂▂▂▂▂▁▂▁▁▁▁▂▁▁▁▁▂▂▂▁▁▁▁▁▁▁▁▁▂▂▂▂▂▂▂▂▂▃▂▂▂▃▂▂▂▂▃▃▃▃▃▄▅▅▅▄▅▃▄▃▄▄▅▅▆▇▆▆▆▆▆▆▆▆▇▇▇▇██
```

Add points in time:

```
$ krapslog -markers 10 /var/log/haproxy.log
                                                             Sat Nov 23 14:15:56
                                                    Sat Nov 23 13:22:29        |
                                           Sat Nov 23 12:29:01        |        |
                                  Sat Nov 23 11:35:33        |        |        |
                          Sat Nov 23 10:48:02       |        |        |        |
                                            |       |        |        |        |
▂▂▂▂▂▁▂▁▁▁▁▂▁▁▁▁▂▂▂▁▁▁▁▁▁▁▁▁▂▂▂▂▂▂▂▂▂▃▂▂▂▃▂▂▂▂▃▃▃▃▃▄▅▅▅▄▅▃▄▃▄▄▅▅▆▇▆▆▆▆▆▆▆▆▇▇▇▇██
|        |        |       |        |
|        |        |       |        Sat Nov 23 09:54:34
|        |        |       Sat Nov 23 09:01:07
|        |        Sat Nov 23 08:13:36
|        Sat Nov 23 07:20:08
Sat Nov 23 06:26:40
```

## Custom date formats

By default, krapslog assumes that log timestamps are in the format "02/Jan/2006:15:04:05.000". However, you can use the `format` parameter to find timestamps in other formats. The parameter value must use the format given in the [documentation](https://golang.org/pkg/time/#Time.Format) for Go's `Time.Format` type.

For example, if your log contains dates that look like  "Jan 1, 2020 15:04:05", you can run krapslog as follows:

```
$ krapslog -format "Jan 2, 2006 15:04:05"
```

## Contributing

Please be kind. We're all trying to do our best.

If you find a bug, please open an issue. (Or, better, submit a pull request that fixes it!)

If you've added a feature, please open a pull request.
